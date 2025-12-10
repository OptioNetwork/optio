package types

const (
	// ModuleName defines the module name
	ModuleName = "lockup"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_lockup"
)

var (
	LockExpirationKey = []byte("lock_expiration")
	TotalLockedKey    = []byte("total_locked")
	AccountLocksKey   = []byte("account_locks")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
