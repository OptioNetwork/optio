package keeper

import (
	"github.com/OptioNetwork/optio/x/lockup/types"
)

var _ types.QueryServer = Keeper{}
