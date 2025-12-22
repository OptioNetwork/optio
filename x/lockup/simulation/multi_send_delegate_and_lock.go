package simulation

import (
	"math/rand"

	"github.com/OptioNetwork/optio/x/lockup/keeper"
	"github.com/OptioNetwork/optio/x/lockup/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

func SimulateMsgMultiSendDelegateAndLock(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgMultiSendDelegateAndLock{
			FromAddress: simAccount.Address.String(),
		}

		// TODO: Handling the MultiSendDelegateAndLock simulation

		return simtypes.NoOpMsg(types.ModuleName, sdk.MsgTypeURL(msg), "MultiSendDelegateAndLock simulation not implemented"), nil, nil
	}
}
