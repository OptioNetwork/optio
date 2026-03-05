package keeper

import (
	"context"
	"encoding/binary"

	"github.com/OptioNetwork/optio/x/licenses/types"
)

// SetLicenseCounter writes the per-type license ID counter directly.
// Used by InitGenesis to rebuild counters from imported licenses.
func (k Keeper) SetLicenseCounter(ctx context.Context, typeID string, id uint64) error {
	store := k.storeService.OpenKVStore(ctx)
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return store.Set(types.GetLicenseCounterKey(typeID), idBz)
}

// GetNextLicenseID reads the per-type counter, increments it, and returns the new ID.
func (k Keeper) GetNextLicenseID(ctx context.Context, typeID string) (uint64, error) {
	store := k.storeService.OpenKVStore(ctx)
	counterKey := types.GetLicenseCounterKey(typeID)

	bz, err := store.Get(counterKey)
	if err != nil {
		return 0, err
	}

	var id uint64
	if bz != nil {
		id = binary.BigEndian.Uint64(bz)
	}
	id++

	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	if err := store.Set(counterKey, idBz); err != nil {
		return 0, err
	}
	return id, nil
}

// SetLicense writes the license to the primary key and updates the holder secondary index.
func (k Keeper) SetLicense(ctx context.Context, l types.License) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&l)
	if err != nil {
		return err
	}
	if err := store.Set(types.GetLicenseKey(l.Type, l.Id), bz); err != nil {
		return err
	}
	return store.Set(types.GetLicenseByHolderKey(l.Holder, l.Type, l.Id), bz)
}

// GetLicense retrieves a license by type and id.
func (k Keeper) GetLicense(ctx context.Context, typeID string, id uint64) (types.License, bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetLicenseKey(typeID, id))
	if err != nil {
		return types.License{}, false, err
	}
	if bz == nil {
		return types.License{}, false, nil
	}
	var l types.License
	if err := k.cdc.Unmarshal(bz, &l); err != nil {
		return types.License{}, false, err
	}
	return l, true, nil
}

// DeleteLicense removes the license from both the primary key and the holder index.
func (k Keeper) DeleteLicense(ctx context.Context, l types.License) error {
	store := k.storeService.OpenKVStore(ctx)
	if err := store.Delete(types.GetLicenseKey(l.Type, l.Id)); err != nil {
		return err
	}
	return store.Delete(types.GetLicenseByHolderKey(l.Holder, l.Type, l.Id))
}

// GetLicensesByType returns all licenses for a given license type by iterating the primary key prefix.
func (k Keeper) GetLicensesByType(ctx context.Context, typeID string) ([]types.License, error) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.GetLicensesByTypePrefix(typeID)
	iterator, err := store.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var licenses []types.License
	for ; iterator.Valid(); iterator.Next() {
		var l types.License
		if err := k.cdc.Unmarshal(iterator.Value(), &l); err != nil {
			return nil, err
		}
		licenses = append(licenses, l)
	}
	return licenses, nil
}

// GetLicensesByHolder returns all licenses held by a given address.
func (k Keeper) GetLicensesByHolder(ctx context.Context, holder string) ([]types.License, error) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.GetLicenseByHolderPrefix(holder)
	iterator, err := store.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var licenses []types.License
	for ; iterator.Valid(); iterator.Next() {
		var l types.License
		if err := k.cdc.Unmarshal(iterator.Value(), &l); err != nil {
			return nil, err
		}
		licenses = append(licenses, l)
	}
	return licenses, nil
}

// GetLicensesByHolderAndType returns all licenses held by a given address for a specific type.
func (k Keeper) GetLicensesByHolderAndType(ctx context.Context, holder, typeID string) ([]types.License, error) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.GetLicenseByHolderAndTypePrefix(holder, typeID)
	iterator, err := store.Iterator(prefix, prefixEndBytes(prefix))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var licenses []types.License
	for ; iterator.Valid(); iterator.Next() {
		var l types.License
		if err := k.cdc.Unmarshal(iterator.Value(), &l); err != nil {
			return nil, err
		}
		licenses = append(licenses, l)
	}
	return licenses, nil
}

// GetAllLicenses returns every license in the store (used for genesis export).
func (k Keeper) GetAllLicenses(ctx context.Context) ([]types.License, error) {
	store := k.storeService.OpenKVStore(ctx)
	iterator, err := store.Iterator(types.LicenseKey, prefixEndBytes(types.LicenseKey))
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var licenses []types.License
	for ; iterator.Valid(); iterator.Next() {
		var l types.License
		if err := k.cdc.Unmarshal(iterator.Value(), &l); err != nil {
			return nil, err
		}
		licenses = append(licenses, l)
	}
	return licenses, nil
}

// prefixEndBytes returns the end key for a prefix iterator.
func prefixEndBytes(prefix []byte) []byte {
	if len(prefix) == 0 {
		return nil
	}
	end := make([]byte, len(prefix))
	copy(end, prefix)
	for i := len(end) - 1; i >= 0; i-- {
		end[i]++
		if end[i] != 0 {
			return end
		}
	}
	return nil
}
