package v2_distro

import (
	store "cosmossdk.io/store/types"
	"github.com/OptioNetwork/optio/app/upgrades"
)

const UpgradeName = "v2-distro"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
