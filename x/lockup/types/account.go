package types

import (
	"sort"
	"time"

	"cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func NewLockupAccount(baseAccount *authtypes.BaseAccount) *Account {
	return &Account{
		BaseAccount: baseAccount,
		Locks:       []*Lock{},
	}
}

func (a *Account) GetLockedAmount(currentTime time.Time) math.Int {
	totalLockedAmount := math.ZeroInt()
	for _, lock := range a.Locks {
		if IsLocked(currentTime, lock.UnlockDate) {
			totalLockedAmount = totalLockedAmount.Add(lock.Coin.Amount)
		}
	}

	return totalLockedAmount
}

// Find lockup using binary search
func (a *Account) FindLockup(unlockDate string) (*Lock, int, bool) {
	idx := sort.Search(len(a.Locks), func(i int) bool {
		return a.Locks[i].UnlockDate >= unlockDate
	})
	if idx < len(a.Locks) && a.Locks[idx].UnlockDate == unlockDate {
		return a.Locks[idx], idx, true
	}
	return nil, 0, false
}

// Insert lockup maintaining sorted order
func (a *Account) InsertLockup(newLockup *Lock) []*Lock {
	idx := sort.Search(len(a.Locks), func(i int) bool {
		return a.Locks[i].UnlockDate >= newLockup.UnlockDate
	})

	// Insert at the correct position
	a.Locks = append(a.Locks, nil)
	copy(a.Locks[idx+1:], a.Locks[idx:])
	a.Locks[idx] = newLockup
	return a.Locks
}

// Upsert (update if exists, insert if not) maintaining sorted order
func (a *Account) UpsertLockup(unlockDate string, coin sdk.Coin) []*Lock {
	idx := sort.Search(len(a.Locks), func(i int) bool {
		return a.Locks[i].UnlockDate >= unlockDate
	})

	if idx < len(a.Locks) && a.Locks[idx].UnlockDate == unlockDate {
		// Update existing
		a.Locks[idx].Coin.Amount = a.Locks[idx].Coin.Amount.Add(coin.Amount)
		return a.Locks
	}

	// Insert new at correct position
	newLockup := &Lock{
		Coin: sdk.NewCoin(coin.Denom, coin.Amount),
	}
	a.Locks = append(a.Locks, nil)
	copy(a.Locks[idx+1:], a.Locks[idx:])
	a.Locks[idx] = newLockup
	return a.Locks
}

// Remove lockup by index
func (a *Account) RemoveLockup(idx int) []*Lock {
	return append(a.Locks[:idx], a.Locks[idx+1:]...)
}

// Range query: get all lockups before a certain date
func (a *Account) GetLockupsBeforeDate(date string) []*Lock {
	idx := sort.Search(len(a.Locks), func(i int) bool {
		return a.Locks[i].UnlockDate >= date
	})
	return a.Locks[:idx]
}

// Implement sdk.AccountI interface
func (a *Account) GetAddress() sdk.AccAddress {
	if a.BaseAccount != nil {
		return a.BaseAccount.GetAddress()
	}
	return nil
}

func (a *Account) SetAddress(addr sdk.AccAddress) error {
	if a.BaseAccount != nil {
		return a.BaseAccount.SetAddress(addr)
	}
	return nil
}

func (a *Account) GetPubKey() cryptotypes.PubKey {
	if a.BaseAccount != nil {
		return a.BaseAccount.GetPubKey()
	}
	return nil
}

func (a *Account) SetPubKey(pubKey cryptotypes.PubKey) error {
	if a.BaseAccount != nil {
		return a.BaseAccount.SetPubKey(pubKey)
	}
	return nil
}

func (a *Account) GetAccountNumber() uint64 {
	if a.BaseAccount != nil {
		return a.BaseAccount.GetAccountNumber()
	}
	return 0
}

func (a *Account) SetAccountNumber(accNumber uint64) error {
	if a.BaseAccount != nil {
		return a.BaseAccount.SetAccountNumber(accNumber)
	}
	return nil
}

func (a *Account) GetSequence() uint64 {
	if a.BaseAccount != nil {
		return a.BaseAccount.GetSequence()
	}
	return 0
}

func (a *Account) SetSequence(seq uint64) error {
	if a.BaseAccount != nil {
		return a.BaseAccount.SetSequence(seq)
	}
	return nil
}
