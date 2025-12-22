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
	LocksByDateKey    = []byte("locks_by_date")
	LocksByAddressKey = []byte("locks_by_address")
	TotalLockedKey    = []byte("total_locked")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
