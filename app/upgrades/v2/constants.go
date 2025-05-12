package v2

import (
	store "cosmossdk.io/store/types"
	"github.com/OptioNetwork/optio/app/upgrades"
	distrotypes "github.com/OptioNetwork/optio/x/distro/types"
)

const UpgradeName = "v2"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{distrotypes.ModuleName},
	},
}
