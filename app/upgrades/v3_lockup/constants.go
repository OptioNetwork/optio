package v3_lockup

import (
	store "cosmossdk.io/store/types"
	"github.com/OptioNetwork/optio/app/upgrades"
	lockuptypes "github.com/OptioNetwork/optio/x/lockup/types"
)

const UpgradeName = "v3-lockup"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{lockuptypes.ModuleName},
	},
}
