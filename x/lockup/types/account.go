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

// Calculate total locked amount at the given time
func (a *Account) GetLockedAmount(currentTime time.Time) math.Int {
	totalLockedAmount := math.ZeroInt()
	for _, lock := range a.Locks {
		if IsLocked(currentTime, lock.UnlockDate) {
			totalLockedAmount = totalLockedAmount.Add(lock.Amount.Amount)
		}
	}

	return totalLockedAmount
}

// Find lockup using binary search
func (a *Account) FindLock(unlockDate string) (*Lock, int, bool) {
	idx := sort.Search(len(a.Locks), func(i int) bool {
		return a.Locks[i].UnlockDate >= unlockDate
	})
	if idx < len(a.Locks) && a.Locks[idx].UnlockDate == unlockDate {
		return a.Locks[idx], idx, true
	}
	return nil, 0, false
}

// Insert lockup maintaining sorted order
func (a *Account) InsertLock(newLockup *Lock) []*Lock {
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
func (a *Account) UpsertLock(unlockDate string, amount math.Int) []*Lock {
	idx := sort.Search(len(a.Locks), func(i int) bool {
		return a.Locks[i].UnlockDate >= unlockDate
	})

	if idx < len(a.Locks) && a.Locks[idx].UnlockDate == unlockDate {
		a.Locks[idx].Amount.Amount = a.Locks[idx].Amount.Amount.Add(amount)
		return a.Locks
	}

	newLockup := &Lock{
		UnlockDate: unlockDate,
		Amount:     &sdk.Coin{Denom: "uOPT", Amount: amount},
	}
	a.Locks = append(a.Locks, nil)
	copy(a.Locks[idx+1:], a.Locks[idx:])
	a.Locks[idx] = newLockup
	return a.Locks
}

// Update lock at a given index
func (a *Account) UpdateLock(idx int, lock *Lock) []*Lock {
	a.Locks[idx] = lock
	return a.Locks
}

// Remove lock by index
func (a *Account) RemoveLock(idx int) []*Lock {
	return append(a.Locks[:idx], a.Locks[idx+1:]...)
}

// Range query: get all locks before a certain date
func (a *Account) GetLocksBeforeDate(date string) []*Lock {
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
