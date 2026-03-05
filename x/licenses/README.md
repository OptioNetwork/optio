# x/licenses

The `licenses` module implements a governance-controlled licensing system for the Optio network. A single **owner** address (set as a module parameter, updatable via governance) controls license type creation and admin delegation. Admins with scoped permissions can issue, revoke, and update licenses. License holders can transfer their own licenses if the type allows it.

## Architecture

```
Owner (module param, gov-updatable)
  ├── creates/updates LicenseTypes
  ├── sets/removes AdminKeys (scoped permission grants)
  │     └── AdminKey grants: {permission, license_types[]}
  │           ├── "issue"  → can mint new licenses (single or batch)
  │           ├── "revoke" → can delete licenses
  │           └── "update" → can change license status
  └── AdminKeys operate on Licenses
        └── License holders can transfer (if type is transferrable)
```

## Data Model

### Params

| Field | Type | Description |
|-------|------|-------------|
| `owner` | `string` | Bech32 address that controls the module. Updated via governance `MsgUpdateParams`. |

### LicenseType

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | Owner-supplied identifier (e.g. `"optio.node"`). Immutable once created. |
| `transferrable` | `bool` | Whether holders can transfer licenses of this type. |
| `max_supply` | `cosmos.Int` | Maximum number of active licenses. `0` = unlimited. |
| `issued_count` | `cosmos.Int` | Current number of active licenses. Incremented on issue, decremented on revoke. |

License types can be updated by the owner via `MsgUpdateLicenseType` (transferrable, max_supply).

### License

| Field | Type | Description |
|-------|------|-------------|
| `id` | `uint64` | Auto-incremented per type (IDs restart per license type). |
| `type` | `string` | The `LicenseType.id` this license belongs to. |
| `holder` | `string` | Bech32 address of the current holder. |
| `start_date` | `string` | Validity start date (`YYYY-MM-DD` format, required). |
| `end_date` | `string` | Expiry date (`YYYY-MM-DD` format). Empty = no expiry. |
| `status` | `string` | `"active"` or `"revoked"`. |

A license is uniquely identified by the composite key `(type, id)`.

### AdminKey

| Field | Type | Description |
|-------|------|-------------|
| `address` | `string` | Bech32 address holding the grants. |
| `grants` | `[]AdminKeyGrant` | List of permission grants. |

Each `AdminKeyGrant` contains:
- `permission` — one of `"issue"`, `"revoke"`, `"update"`
- `license_types` — list of license type IDs this permission applies to

## Transactions

| Message | Signer | Description |
|---------|--------|-------------|
| `MsgUpdateParams` | governance authority | Update the module `owner` parameter. |
| `MsgCreateLicenseType` | owner | Create a new license type. |
| `MsgSetAdminKey` | owner | Set/replace all grants for an address. |
| `MsgRemoveAdminKey` | owner | Remove all grants for an address. |
| `MsgUpdateLicenseType` | owner | Update a license type's `transferrable` and `max_supply` fields. |
| `MsgIssueLicense` | admin (needs `issue` grant) | Mint one or more licenses to a holder (supports `count` field). |
| `MsgBatchIssueLicense` | admin (needs `issue` grant) | Issue licenses to multiple holders in a single transaction. |
| `MsgRevokeLicense` | admin (needs `revoke` grant) | Delete a license (physically removed from store). |
| `MsgUpdateLicense` | admin (needs `update` grant) | Change a license's status (`"active"` or `"revoked"`). |
| `MsgTransferLicense` | license holder | Transfer a license to another address (type must be transferrable). |

### Validation Rules

**SetAdminKey:**
- Address must be valid bech32
- Each grant must have a valid permission (`issue`/`revoke`/`update`)
- Each grant must include at least one license type

**UpdateLicenseType:**
- Only the module owner can call
- `max_supply` cannot be reduced below `issued_count`

**IssueLicense:**
- Holder must be valid bech32
- `start_date` is required, must be `YYYY-MM-DD`
- `end_date` (if provided) must be `YYYY-MM-DD` and >= `start_date`
- Max supply is enforced (if non-zero)
- `count` field (optional, default 1) creates multiple identical licenses for the same holder

**BatchIssueLicense:**
- Same validation rules as IssueLicense per entry
- All entries validated before any licenses are created
- Max supply checked against total entry count upfront

**TransferLicense:**
- Recipient must be valid bech32
- Cannot transfer to self
- License type must be transferrable

**UpdateLicense:**
- Status must be `"active"` or `"revoked"`

## Queries

| RPC | REST Endpoint | Description |
|-----|---------------|-------------|
| `Params` | `GET /optio/licenses/params` | Module parameters |
| `LicenseType` | `GET /optio/licenses/license_type/{id}` | Single license type |
| `LicenseTypes` | `GET /optio/licenses/license_types` | All license types (paginated) |
| `License` | `GET /optio/licenses/license/{type_id}/{id}` | Single license |
| `LicensesByType` | `GET /optio/licenses/licenses_by_type/{type_id}` | All licenses for a type (paginated) |
| `LicensesByHolder` | `GET /optio/licenses/licenses_by_holder/{holder}` | All licenses for a holder (paginated) |
| `LicensesByHolderAndType` | `GET /optio/licenses/licenses_by_holder/{holder}/{type_id}` | Holder's licenses for a specific type (paginated) |
| `AdminKey` | `GET /optio/licenses/admin_key/{address}` | Grants for an address |
| `AdminKeys` | `GET /optio/licenses/admin_keys` | All admin keys (paginated) |
| `AdminKeysByLicenseType` | `GET /optio/licenses/admin_keys_by_license_type/{license_type_id}` | Admin keys with grants for a license type (paginated, optional `permission` filter) |

## CLI Usage

Most commands are auto-generated via AutoCLI with positional arguments. The `set-admin-key` and `batch-issue-license` commands use custom CLIs.

### Transactions

```bash
# Create a license type (positional: owner, id; flags: transferrable, max-supply)
optiod tx licenses create-license-type <owner-addr> optio.node \
  --transferrable true --max-supply 1000 --from owner

# Update a license type (positional: owner, id; flags: transferrable, max-supply)
optiod tx licenses update-license-type <owner-addr> optio.node \
  --transferrable true --max-supply 2000 --from owner

# Grant admin permissions (custom CLI — positional, comma-delimited)
optiod tx licenses set-admin-key optio1abc... issue,revoke optio.node,optio.validator \
  --from owner

# Remove admin key (positional: owner, address)
optiod tx licenses remove-admin-key <owner-addr> optio1abc... --from owner

# Issue a license (positional: issuer, license-type-id, holder, start-date; flags: end-date, count)
optiod tx licenses issue-license <issuer-addr> optio.node optio1xyz... 2026-01-01 \
  --end-date 2027-01-01 --from admin

# Issue multiple licenses to one holder
optiod tx licenses issue-license <issuer-addr> optio.node optio1xyz... 2026-01-01 \
  --count 5 --from admin

# Batch issue to multiple holders (custom CLI — colon-delimited entries)
optiod tx licenses batch-issue-license optio.node \
  optio1abc...:2026-01-01:2027-01-01 \
  optio1def...:2026-01-01 \
  --from admin

# Revoke a license (positional: revoker, license-type-id, id)
optiod tx licenses revoke-license <revoker-addr> optio.node 1 --from admin

# Update license status (positional: updater, license-type-id, id, status)
optiod tx licenses update-license <updater-addr> optio.node 1 revoked --from admin

# Transfer a license (positional: holder, license-type-id, id, recipient)
optiod tx licenses transfer-license <holder-addr> optio.node 1 optio1new... --from holder
```

### Queries

```bash
optiod query licenses params
optiod query licenses license-type optio.node
optiod query licenses license-types
optiod query licenses license optio.node 1
optiod query licenses licenses-by-type optio.node
optiod query licenses licenses-by-holder optio1abc...
optiod query licenses licenses-by-holder-and-type optio1abc... optio.node
optiod query licenses admin-key optio1abc...
optiod query licenses admin-keys
optiod query licenses admin-keys-by-license-type optio.node
optiod query licenses admin-keys-by-license-type optio.node --permission issue
```

## Store Layout

| Key Pattern | Value |
|-------------|-------|
| `params` | `Params` |
| `license_type/{type_id}` | `LicenseType` |
| `license_counter/{type_id}` | `uint64` (8-byte big-endian per-type ID counter) |
| `license/{type_id}/{8-byte-id}` | `License` (primary index) |
| `license_by_holder/{holder}/{type_id}/{8-byte-id}` | `License` (secondary holder index) |
| `admin_key/{address}` | `AdminKey` |

The holder secondary index supports efficient prefix scans for both `LicensesByHolder` (prefix on `{holder}/`) and `LicensesByHolderAndType` (prefix on `{holder}/{type_id}/`).

## Events

| Event Type | Attributes |
|------------|------------|
| `create_license_type` | `license_type_id` |
| `set_admin_key` | `address`, `permissions`, `grant_license_types` |
| `remove_admin_key` | `address` |
| `update_license_type` | `license_type_id` |
| `issue_license` | `license_type_id`, `holder`, `count` |
| `batch_issue_license` | `license_type_id`, `count` |
| `revoke_license` | `license_type_id`, `license_id` |
| `update_license` | `license_type_id`, `license_id`, `status` |
| `transfer_license` | `license_type_id`, `license_id`, `holder`, `recipient` |
| `update_params` | `owner` |

## Error Codes

| Code | Name | Description |
|------|------|-------------|
| 1100 | `ErrInvalidSigner` | Expected governance account as signer |
| 1101 | `ErrLicenseTypeNotFound` | License type not found |
| 1102 | `ErrLicenseTypeExists` | License type already exists |
| 1103 | `ErrMaxSupplyReached` | License type max supply reached |
| 1104 | `ErrLicenseNotFound` | License not found |
| 1105 | `ErrNotLicenseHolder` | Signer is not the license holder |
| 1106 | `ErrLicenseNotTransferable` | License type is not transferrable |
| 1107 | `ErrUnauthorized` | Missing required admin key permission |

## Genesis

The genesis state includes `params`, `license_types`, `licenses`, and `admin_keys`. On import (`InitGenesis`), per-type license ID counters are automatically rebuilt from the maximum license ID seen per type, preventing ID collisions after chain export/import.

Validation checks:
- No duplicate license type IDs
- No duplicate `(type, id)` license pairs
- All licenses reference a declared license type
- All license statuses are `"active"` or `"revoked"`
- All holder and admin key addresses are valid bech32
- No duplicate admin key addresses
