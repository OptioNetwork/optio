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
	ParamsKey = []byte("p_lockup")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
