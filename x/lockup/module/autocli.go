package lockup

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "github.com/OptioNetwork/optio/api/optio/lockup"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service:           modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				// Lock command is provided by custom CLI (see cli/tx.go)
				{
					RpcMethod:      "Extend",
					Use:            "extend [extensions]",
					Short:          "Send a Extend tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "extensions"}},
				},
				{
					RpcMethod: "SendDelegateAndLock",
					Use:       "send-delegate-and-lock [to-address] [val-address] [amount] [unlockDate]",
					Short:     "Send a SendDelegateAndLock tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "toAddress"},
						{ProtoField: "valAddress"},
						{ProtoField: "amount"},
						{ProtoField: "unlockDate"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
