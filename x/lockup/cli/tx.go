package cli

import (
	"fmt"
	"strings"

	"github.com/OptioNetwork/optio/x/lockup/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdLock())

	return cmd
}

func CmdLock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [unlock-date:amount] [unlock-date:amount]...",
		Short: "Lock tokens until specific dates",
		Long:  "Lock tokens until specific unlock dates. You can specify multiple locks.",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			locks := make([]*types.Lock, 0, len(args))
			for i, arg := range args {
				parts := strings.Split(arg, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid lock format at position %d: expected 'date:amount', got '%s'", i, arg)
				}

				unlockDate := parts[0]
				coin, err := sdk.ParseCoinNormalized(parts[1])
				if err != nil {
					return fmt.Errorf("invalid coin at position %d: %w", i, err)
				}

				locks = append(locks, &types.Lock{
					UnlockDate: unlockDate,
					Coin:       coin,
				})
			}

			msg := &types.MsgLock{
				LockupAddress: clientCtx.GetFromAddress().String(),
				Locks:         locks,
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
