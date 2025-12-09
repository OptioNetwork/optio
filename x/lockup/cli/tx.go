package cli

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
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
	cmd.AddCommand(CmdExtend())
	cmd.AddCommand(CmdSendDelegateAndLock())
	cmd.AddCommand(CmdMultiSendDelegateAndLock())

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
				Address: clientCtx.GetFromAddress().String(),
				Locks:   locks,
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

func CmdExtend() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "extend [from-date:to-date:amount] [from-date:to-date:amount]...",
		Short: "Extend lock unlock dates",
		Long: `Extend the unlock date of existing locks. You can specify multiple extensions.
		Example: '2026-12-01:2027-12-01:1000000000uOPT' extends a lock that unlocks on 2026-12-01 by locking an additional 1000000000uOPT until 2027-12-01.`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			extensions := make([]*types.Extension, 0, len(args))
			for i, arg := range args {
				parts := strings.Split(arg, ":")
				if len(parts) != 3 {
					return fmt.Errorf("invalid extension format at position %d: expected 'from-date:to-date:amount', got '%s'", i, arg)
				}

				fromDate := parts[0]
				toDate := parts[1]
				coin, err := sdk.ParseCoinNormalized(parts[2])
				if err != nil {
					return fmt.Errorf("invalid coin at position %d: %w", i, err)
				}

				extensions = append(extensions, &types.Extension{
					From: fromDate,
					Lock: &types.Lock{
						UnlockDate: toDate,
						Coin:       coin,
					},
				})
			}

			msg := &types.MsgExtend{
				Address:    clientCtx.GetFromAddress().String(),
				Extensions: extensions,
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

func CmdSendDelegateAndLock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-delegate-and-lock [to-address] [validator-address] [unlock-date] [amount]",
		Short: "Send tokens to an address, delegate them to a validator, and lock them",
		Long: `Send tokens from your address to another address, delegate them to a validator, and lock them until a specific unlock date.
Example: 
  send-delegate-and-lock optio1abc... optiovaloper1xyz... 2026-12-01 1000uOPT`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			toAddress := args[0]
			validatorAddress := args[1]
			unlockDate := args[2]
			coin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return fmt.Errorf("invalid coin amount: %w", err)
			}

			msg := &types.MsgSendDelegateAndLock{
				FromAddress:      clientCtx.GetFromAddress().String(),
				ToAddress:        toAddress,
				ValidatorAddress: validatorAddress,
				Lock: &types.Lock{
					UnlockDate: unlockDate,
					Coin:       coin,
				},
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

func CmdMultiSendDelegateAndLock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi-send-delegate-and-lock [to-address:validator-address:unlock-date:amount] [to-address:validator-address:unlock-date:amount]...",
		Short: "Send tokens to multiple addresses, delegate them, and lock them",
		Long: `Send tokens to multiple addresses, delegate them to validators, and lock them until specific unlock dates.
Example: 
  multi-send-delegate-and-lock optio1abc...:optiovaloper1xyz...:2026-12-01:1000uOPT optio1def...:optiovaloper1uvw...:2027-01-01:2000uOPT`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			outputs := make([]*types.Output, 0, len(args))
			totalAmount := math.ZeroInt()

			for i, arg := range args {
				parts := strings.Split(arg, ":")
				if len(parts) != 4 {
					return fmt.Errorf("invalid output format at position %d: expected 'to-address:validator-address:unlock-date:amount', got '%s'", i, arg)
				}

				toAddress := parts[0]
				validatorAddress := parts[1]
				unlockDate := parts[2]
				coin, err := sdk.ParseCoinNormalized(parts[3])
				if err != nil {
					return fmt.Errorf("invalid coin at position %d: %w", i, err)
				}

				amount, ok := math.NewIntFromString(coin.Amount.String())
				if !ok {
					return fmt.Errorf("invalid coin amount at position %d: %w", i, err)
				}
				totalAmount = totalAmount.Add(amount)

				outputs = append(outputs, &types.Output{
					ToAddress:  toAddress,
					ValAddress: validatorAddress,
					Lock: &types.Lock{
						UnlockDate: unlockDate,
						Coin:       coin,
					},
				})
			}

			msg := &types.MsgMultiSendDelegateAndLock{
				FromAddress: clientCtx.GetFromAddress().String(),
				TotalAmount: totalAmount.String(),
				Outputs:     outputs,
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
