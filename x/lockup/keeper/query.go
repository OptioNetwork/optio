package keeper

import (
	"context"
	"encoding/binary"
	"strings"
	"time"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		startKey = types.LocksByDateKey
	}

	iterator, err := store.Iterator(startKey, prefixEndBytes(types.LocksByDateKey))
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
		prefixLen := len(types.LocksByDateKey)
		if len(key) < prefixLen+8 {
			continue
		}

		timeBz := key[prefixLen : prefixLen+8]
		addrBz := key[prefixLen+8:]

		unlockUnix := binary.BigEndian.Uint64(timeBz)
		unlockTime := time.Unix(int64(unlockUnix), 0)
		unlockDate := unlockTime.UTC().Format(time.DateOnly)

		blockTime := ctx.BlockTime()
		blockDay := time.Date(blockTime.Year(), blockTime.Month(), blockTime.Day(), 0, 0, 0, 0, time.UTC)

		if !types.IsLocked(blockDay, unlockDate) {
			continue
		}

		addr := sdk.AccAddress(addrBz)

		var amount math.Int
		if err := amount.Unmarshal(value); err != nil {
			return nil, status.Error(codes.Internal, "failed to unmarshal amount")
		}

		bondDenom, err := k.stakingKeeper.BondDenom(ctx)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		locks = append(locks, types.ActiveLock{
			Address:    addr.String(),
			UnlockDate: unlockDate,
			Amount:     sdk.NewCoin(bondDenom, amount),
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

	totalLocked := math.ZeroInt()

	blockTime := ctx.BlockTime()
	blockDay := time.Date(blockTime.Year(), blockTime.Month(), blockTime.Day(), 0, 0, 0, 0, time.UTC)

	err := k.IterateActiveLocks(ctx, blockDay, func(addr sdk.AccAddress, unlockTime time.Time, amount math.Int) error {
		totalLocked = totalLocked.Add(amount)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, sdkerrors.ErrInvalidType.Wrapf("invalid bond denom: %s", err)
	}

	return &types.QueryTotalLockedAmountResponse{
		TotalLocked: sdk.NewCoin(bondDenom, totalLocked),
	}, nil
}

func (k Keeper) LocksForAddresses(goCtx context.Context, req *types.QueryLocksForAddressesRequest) (*types.QueryLocksForAddressesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Addresses == "" {
		return nil, status.Error(codes.InvalidArgument, "addresses parameter is required")
	}

	addressList := strings.Split(req.Addresses, ",")
	for i, addr := range addressList {
		addressList[i] = strings.TrimSpace(addr)
	}

	if len(addressList) == 0 {
		return nil, status.Error(codes.InvalidArgument, "at least one address is required")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	accountLocks := make([]types.AccountLocks, 0, len(addressList))

	for _, addrStr := range addressList {
		addr, err := sdk.AccAddressFromBech32(addrStr)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid address: "+addrStr)
		}

		locks, err := k.GetLocksByAddress(ctx, addr)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		blockTime := ctx.BlockTime()
		blockDay := time.Date(blockTime.Year(), blockTime.Month(), blockTime.Day(), 0, 0, 0, 0, time.UTC)

		activeLocks := make([]types.Lock, 0)
		for _, lock := range locks {

			if types.IsLocked(blockDay, lock.UnlockDate) {
				activeLocks = append(activeLocks, *lock)
			}

		}

		accountLocks = append(accountLocks, types.AccountLocks{
			Address: addrStr,
			Locks:   activeLocks,
		})
	}

	return &types.QueryLocksForAddressesResponse{
		Locks: accountLocks,
	}, nil
}

// Locks implements types.QueryServer.
func (k Keeper) Locks(goCtx context.Context, req *types.QueryLocksRequest) (*types.QueryLocksResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address: "+req.Address)
	}

	locks, err := k.GetLocksByAddress(ctx, addr)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	blockTime := ctx.BlockTime()
	blockDay := time.Date(blockTime.Year(), blockTime.Month(), blockTime.Day(), 0, 0, 0, 0, time.UTC)

	activeLocks := make([]types.Lock, 0)
	for _, lock := range locks {

		if types.IsLocked(blockDay, lock.UnlockDate) {
			activeLocks = append(activeLocks, *lock)
		}

	}

	return &types.QueryLocksResponse{
		Locks: activeLocks,
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
