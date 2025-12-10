package keeper

import (
	"context"
	"encoding/binary"
	"time"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetTotalLocked sets the total locked amount
func (k Keeper) SetTotalLocked(ctx context.Context, amount math.Int) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := amount.Marshal()
	if err != nil {
		return err
	}
	return store.Set(types.TotalLockedKey, bz)
}

// GetTotalLocked returns the total locked amount
func (k Keeper) GetTotalLocked(ctx context.Context) (math.Int, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.TotalLockedKey)
	if err != nil {
		return math.ZeroInt(), err
	}
	if bz == nil {
		return math.ZeroInt(), nil
	}
	var amount math.Int
	err = amount.Unmarshal(bz)
	if err != nil {
		return math.ZeroInt(), err
	}
	return amount, nil
}

// GetLockExpirationKey creates the key for the lock expiration queue
// Key: Prefix + Timestamp (8 bytes) + Address
func (k Keeper) GetLockExpirationKey(unlockTime time.Time, addr sdk.AccAddress) []byte {
	timeBz := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBz, uint64(unlockTime.Unix()))
	return append(append(types.LockExpirationKey, timeBz...), addr.Bytes()...)
}

// AddToExpirationQueue adds a lock to the expiration queue
// If entry exists, adds the amount
func (k Keeper) AddToExpirationQueue(ctx context.Context, unlockTime time.Time, addr sdk.AccAddress, amount math.Int) error {
	store := k.storeService.OpenKVStore(ctx)
	key := k.GetLockExpirationKey(unlockTime, addr)

	// Check existing
	bz, err := store.Get(key)
	if err != nil {
		return err
	}

	currentAmount := math.ZeroInt()
	if bz != nil {
		if err := currentAmount.Unmarshal(bz); err != nil {
			return err
		}
	}

	newAmount := currentAmount.Add(amount)
	bz, err = newAmount.Marshal()
	if err != nil {
		return err
	}

	return store.Set(key, bz)
}

// IterateAndDeleteExpiredLocks iterates over locks that have expired before or at cutoffTime and deletes them
func (k Keeper) IterateAndDeleteExpiredLocks(ctx context.Context, cutoffTime time.Time, cb func(addr sdk.AccAddress, unlockTime time.Time, amount math.Int) error) error {
	store := k.storeService.OpenKVStore(ctx)

	// End key is Prefix + CutoffTime + 1 second (to include CutoffTime)
	endTimeBz := make([]byte, 8)
	binary.BigEndian.PutUint64(endTimeBz, uint64(cutoffTime.Unix()+1))
	endKey := append(types.LockExpirationKey, endTimeBz...)

	iter, err := store.Iterator(types.LockExpirationKey, endKey)
	if err != nil {
		return err
	}
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		// Parse Key
		// Prefix (len) + Time (8) + Addr (Remainder)
		prefixLen := len(types.LockExpirationKey)
		timeBz := key[prefixLen : prefixLen+8]
		addrBz := key[prefixLen+8:]

		unlockUnix := binary.BigEndian.Uint64(timeBz)
		unlockTime := time.Unix(int64(unlockUnix), 0)
		addr := sdk.AccAddress(addrBz)

		var amount math.Int
		if err := amount.Unmarshal(iter.Value()); err != nil {
			return err
		}

		if err := cb(addr, unlockTime, amount); err != nil {
			return err
		}

		if err := store.Delete(key); err != nil {
			return err
		}
	}

	return nil
}
