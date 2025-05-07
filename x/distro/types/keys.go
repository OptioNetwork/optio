package types

const (
	// ModuleName defines the module name
	ModuleName = "distro"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_distro"
)

var (
	ParamsKey = []byte("p_distro")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
