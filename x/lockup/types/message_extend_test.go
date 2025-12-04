package types

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgExtend_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgExtend
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgExtend{
				ExtendingAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty extensions",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions:       []*Extension{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "nil extension",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions:       []*Extension{nil},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "empty extension from date",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: "",
						Lock: &Lock{
							Coin:       sdk.NewInt64Coin("uOPT", 100),
							UnlockDate: time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "malformed from date",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: "12-04-2025",
						Lock: &Lock{},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "nil lock in extension",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: nil,
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid lock amount",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin: sdk.Coin{
								Denom:  "uOPT",
								Amount: math.NewInt(-100),
							},
							UnlockDate: time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "invalid lock amount",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin: sdk.Coin{
								Denom:  "uOPT",
								Amount: math.NewInt(0),
							},
							UnlockDate: time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "empty lock denom",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin: sdk.Coin{
								Denom:  "",
								Amount: math.NewInt(0),
							},
							UnlockDate: time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "invalid lock denom",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin: sdk.Coin{
								Denom:  "invalid_denom",
								Amount: math.NewInt(0),
							},
							UnlockDate: time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "empty unlock date",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin:       sdk.NewInt64Coin("uOPT", 1000000000),
							UnlockDate: "",
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "malformed unlock date",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin:       sdk.NewInt64Coin("uOPT", 1000000000),
							UnlockDate: "12-04-2025",
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "unlock date before from date",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin:       sdk.NewInt64Coin("uOPT", 1000000000),
							UnlockDate: time.Now().AddDate(0, 0, -1).Format(time.DateOnly),
						},
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "successful validation",
			msg: MsgExtend{
				ExtendingAddress: sample.AccAddress(),
				Extensions: []*Extension{
					{
						From: time.Now().Format(time.DateOnly),
						Lock: &Lock{
							Coin:       sdk.NewInt64Coin("uOPT", 1000000000),
							UnlockDate: time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
						},
					},
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
