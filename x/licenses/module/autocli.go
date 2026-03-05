package licenses

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/OptioNetwork/optio/api/optio/licenses"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query module parameters",
				},
				{
					RpcMethod: "LicenseType",
					Use:       "license-type [id]",
					Short:     "Query a license type by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "LicenseTypes",
					Use:       "license-types",
					Short:     "Query all license types",
				},
				{
					RpcMethod: "License",
					Use:       "license [type-id] [id]",
					Short:     "Query a license by type id and license id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "type_id"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "LicensesByType",
					Use:       "licenses-by-type [type-id]",
					Short:     "Query all licenses for a given type",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "type_id"},
					},
				},
				{
					RpcMethod: "LicensesByHolder",
					Use:       "licenses-by-holder [holder]",
					Short:     "Query all licenses held by an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "holder"},
					},
				},
				{
					RpcMethod: "LicensesByHolderAndType",
					Use:       "licenses-by-holder-and-type [holder] [type-id]",
					Short:     "Query all licenses held by an address for a specific type",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "holder"},
						{ProtoField: "type_id"},
					},
				},
				{
					RpcMethod: "AdminKey",
					Use:       "admin-key [address]",
					Short:     "Query admin key grants for an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "AdminKeys",
					Use:       "admin-keys",
					Short:     "Query all admin keys",
				},
				{
					RpcMethod: "AdminKeysByLicenseType",
					Use:       "admin-keys-by-license-type [license-type-id]",
					Short:     "Query admin keys that have grants for a given license type",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "license_type_id"},
					},
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // allows custom CLI commands (e.g. set-admin-key) to be merged
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params [authority] [owner]",
					Short:     "Update module parameters (governance)",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "authority"},
					},
				},
				{
					RpcMethod: "CreateLicenseType",
					Use:       "create-license-type [owner] [id]",
					Short:     "Create a new license type",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "SetAdminKey",
					Skip:      true, // provided by custom CLI (see client/cli/tx.go)
				},
				{
					RpcMethod: "RemoveAdminKey",
					Use:       "remove-admin-key [owner] [address]",
					Short:     "Remove admin key for an address",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "address"},
					},
				},
				{
					RpcMethod: "IssueLicense",
					Use:       "issue-license [issuer] [license-type-id] [holder] [start-date]",
					Short:     "Issue one or more licenses",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "issuer"},
						{ProtoField: "license_type_id"},
						{ProtoField: "holder"},
						{ProtoField: "start_date"},
					},
				},
				{
					RpcMethod: "RevokeLicense",
					Use:       "revoke-license [revoker] [license-type-id] [id]",
					Short:     "Revoke a license",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "revoker"},
						{ProtoField: "license_type_id"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "UpdateLicense",
					Use:       "update-license [updater] [license-type-id] [id] [status]",
					Short:     "Update the status of a license",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "updater"},
						{ProtoField: "license_type_id"},
						{ProtoField: "id"},
						{ProtoField: "status"},
					},
				},
				{
					RpcMethod: "TransferLicense",
					Use:       "transfer-license [holder] [license-type-id] [id] [recipient]",
					Short:     "Transfer a license to a new holder",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "holder"},
						{ProtoField: "license_type_id"},
						{ProtoField: "id"},
						{ProtoField: "recipient"},
					},
				},
				{
					RpcMethod: "UpdateLicenseType",
					Use:       "update-license-type [owner] [id]",
					Short:     "Update a license type's properties",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "owner"},
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "BatchIssueLicense",
					Skip:      true, // provided by custom CLI (repeated message field)
				},
			},
		},
	}
}
