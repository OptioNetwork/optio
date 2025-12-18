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
				Address: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty extensions",
			msg: MsgExtend{
				Address:    sample.AccAddress(),
				Extensions: []*Extension{},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "nil extension",
			msg: MsgExtend{
				Address:    sample.AccAddress(),
				Extensions: []*Extension{nil},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "empty extension from date",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: "",
						Amount:   sdk.NewCoin(bondDenom, math.NewInt(100)),
						ToDate:   time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "malformed from date",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: "12-04-2025",
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "nil lock in extension",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid lock amount",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
						Amount:   sdk.NewCoin(bondDenom, math.NewInt(-100)),
						ToDate:   time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "invalid lock amount",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
						ToDate:   time.Now().AddDate(2, 0, 0).Format(time.DateOnly),
						Amount:   sdk.NewCoin(bondDenom, math.ZeroInt()),
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		}, {
			name: "empty unlock date",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
						ToDate:   "",
						Amount:   sdk.NewCoin(bondDenom, math.NewInt(1000000000)),
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "malformed unlock date",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
						ToDate:   "12-04-2025",
						Amount:   sdk.NewCoin(bondDenom, math.NewInt(1000000000)),
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "unlock date before from date",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
						ToDate:   time.Now().AddDate(0, 0, -1).Format(time.DateOnly),
						Amount:   sdk.NewCoin(bondDenom, math.NewInt(1000000000)),
					},
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "successful validation",
			msg: MsgExtend{
				Address: sample.AccAddress(),
				Extensions: []*Extension{
					{
						FromDate: time.Now().Format(time.DateOnly),
						ToDate:   time.Now().AddDate(0, 0, 1).Format(time.DateOnly),
						Amount:   sdk.NewCoin(bondDenom, math.NewInt(1000000000)),
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
