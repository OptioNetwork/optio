package lockup

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/OptioNetwork/optio/testutil/sample"
	lockupsimulation "github.com/OptioNetwork/optio/x/lockup/simulation"
	"github.com/OptioNetwork/optio/x/lockup/types"
)

// avoid unused import issue
var (
	_ = lockupsimulation.FindAccount
	_ = rand.Rand{}
	_ = sample.AccAddress
	_ = sdk.AccAddress{}
	_ = simulation.MsgEntryKind
)

const (
	opWeightMsgLock = "op_weight_msg_lock"
	// TODO: Determine the simulation weight value
	defaultWeightMsgLock int = 100

	opWeightMsgExtend = "op_weight_msg_extend"
	// TODO: Determine the simulation weight value
	defaultWeightMsgExtend int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	lockupGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&lockupGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgLock int
	simState.AppParams.GetOrGenerate(opWeightMsgLock, &weightMsgLock, nil,
		func(_ *rand.Rand) {
			weightMsgLock = defaultWeightMsgLock
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgLock,
		lockupsimulation.SimulateMsgLock(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgExtend int
	simState.AppParams.GetOrGenerate(opWeightMsgExtend, &weightMsgExtend, nil,
		func(_ *rand.Rand) {
			weightMsgExtend = defaultWeightMsgExtend
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgExtend,
		lockupsimulation.SimulateMsgExtend(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgLock,
			defaultWeightMsgLock,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				lockupsimulation.SimulateMsgLock(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgExtend,
			defaultWeightMsgExtend,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				lockupsimulation.SimulateMsgExtend(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
