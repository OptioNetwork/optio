package types

import "encoding/binary"

const (
	ModuleName = "licenses"
	StoreKey   = ModuleName
	MemStoreKey = "mem_licenses"
)

var (
	ParamsKey          = []byte("params")
	LicenseTypeKey     = []byte("license_type/")
	LicenseCounterKey  = []byte("license_counter/")
	LicenseKey         = []byte("license/")
	LicenseByHolderKey = []byte("license_by_holder/")
	AdminKeyStoreKey   = []byte("admin_key/")
)

func GetLicenseTypeKey(typeID string) []byte {
	return append(LicenseTypeKey, []byte(typeID)...)
}

func GetLicenseCounterKey(typeID string) []byte {
	return append(LicenseCounterKey, []byte(typeID)...)
}

func GetLicenseKey(typeID string, id uint64) []byte {
	prefix := append(LicenseKey, []byte(typeID+"/")...)
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return append(prefix, idBz...)
}

func GetLicensesByTypePrefix(typeID string) []byte {
	return append(LicenseKey, []byte(typeID+"/")...)
}

func GetLicenseByHolderKey(holder, typeID string, id uint64) []byte {
	prefix := append(LicenseByHolderKey, []byte(holder+"/"+typeID+"/")...)
	idBz := make([]byte, 8)
	binary.BigEndian.PutUint64(idBz, id)
	return append(prefix, idBz...)
}

func GetLicenseByHolderPrefix(holder string) []byte {
	return append(LicenseByHolderKey, []byte(holder+"/")...)
}

func GetLicenseByHolderAndTypePrefix(holder, typeID string) []byte {
	return append(LicenseByHolderKey, []byte(holder+"/"+typeID+"/")...)
}

func GetAdminKeyKey(address string) []byte {
	return append(AdminKeyStoreKey, []byte(address)...)
}
