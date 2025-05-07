package keeper

import (
	"github.com/OptioNetwork/optio/x/distro/types"
)

var _ types.QueryServer = Keeper{}
