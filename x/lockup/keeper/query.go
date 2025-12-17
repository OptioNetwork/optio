package keeper

import (
	"context"
	"encoding/binary"
	"time"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ActiveLocks(goCtx context.Context, req *types.QueryActiveLocksRequest) (*types.QueryActiveLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := k.storeService.OpenKVStore(ctx)

	var locks []types.ActiveLock

	var startKey []byte
	if req.Pagination != nil && len(req.Pagination.Key) != 0 {
		startKey = req.Pagination.Key
	} else {
		startKey = types.LockExpirationKey
	}

	iterator, err := store.Iterator(startKey, prefixEndBytes(types.LockExpirationKey))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer iterator.Close()

	limit := uint64(100)
	if req.Pagination != nil && req.Pagination.Limit != 0 {
		limit = req.Pagination.Limit
	}

	count := uint64(0)
	var nextKey []byte

	for ; iterator.Valid(); iterator.Next() {
		if count >= limit {
			nextKey = iterator.Key()
			break
		}

		key := iterator.Key()
		value := iterator.Value()

		// Decode
		// Key: Prefix + Timestamp (8) + Address
		prefixLen := len(types.LockExpirationKey)
		if len(key) < prefixLen+8 {
			continue
		}

		timeBz := key[prefixLen : prefixLen+8]
		addrBz := key[prefixLen+8:]

		unlockUnix := binary.BigEndian.Uint64(timeBz)
		unlockTime := time.Unix(int64(unlockUnix), 0)
		addr := sdk.AccAddress(addrBz)

		var amount math.Int
		if err := amount.Unmarshal(value); err != nil {
			return nil, status.Error(codes.Internal, "failed to unmarshal amount")
		}

		locks = append(locks, types.ActiveLock{
			Address:    addr.String(),
			UnlockDate: unlockTime.Format(time.DateOnly),
			Amount:     amount,
		})

		count++
	}

	return &types.QueryActiveLocksResponse{
		Locks: locks,
		Pagination: &query.PageResponse{
			NextKey: nextKey,
		},
	}, nil
}

func (k Keeper) TotalLockedAmount(goCtx context.Context, req *types.QueryTotalLockedAmountRequest) (*types.QueryTotalLockedAmountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Calculate total locked by iterating all active locks in the expiration queue
	totalLocked := math.ZeroInt()

	blockTime := ctx.BlockTime()
	blockTimeDateOnly, err := time.Parse(time.DateOnly, blockTime.Format(time.DateOnly))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = k.IterateActiveLocks(ctx, blockTimeDateOnly, func(addr sdk.AccAddress, unlockTime time.Time, amount math.Int) error {
		totalLocked = totalLocked.Add(amount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryTotalLockedAmountResponse{
		TotalLocked: totalLocked,
	}, nil
}

func (k Keeper) AccountLocks(goCtx context.Context, req *types.QueryAccountLocksRequest) (*types.QueryAccountLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	const maxAddresses = 100
	if len(req.Addresses) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one address is required")
	}
	if len(req.Addresses) > maxAddresses {
		return nil, status.Error(codes.InvalidArgument, "too many addresses: maximum is 100")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	accountLocks := make([]types.AccountLocks, 0, len(req.Addresses))

	for _, addrStr := range req.Addresses {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid address: "+addrStr)
		}

		acc := k.accountKeeper.GetAccount(ctx, addr)
		if acc == nil {
			accountLocks = append(accountLocks, types.AccountLocks{
				Address: addrStr,
				Locks:   []types.Lock{},
			})
			continue
		}

		lockupAcc, ok := acc.(*types.Account)
		if !ok {
			accountLocks = append(accountLocks, types.AccountLocks{
				Address: addrStr,
				Locks:   []types.Lock{},
			})
			continue
		}

		now := ctx.BlockTime()
		activeLocks := make([]types.Lock, 0)
		for _, lock := range lockupAcc.Locks {
			unlockTime, err := time.Parse(time.DateOnly, lock.UnlockDate)
			if err != nil {
				continue
			}
			if unlockTime.After(now) {
				activeLocks = append(activeLocks, *lock)
			}
		}

		accountLocks = append(accountLocks, types.AccountLocks{
			Address: addrStr,
			Locks:   activeLocks,
		})
	}

	return &types.QueryAccountLocksResponse{
		Locks: accountLocks,
	}, nil
}

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
