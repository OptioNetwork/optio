package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/OptioNetwork/optio/testutil/sample"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgLock_ValidateBasic(t *testing.T) {
	validAddr := sample.AccAddress()

	tests := []struct {
		name string
		msg  MsgLock
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgLock{
				Address:    "invalid_address",
				UnlockDate: "2026-12-01",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(1000)),
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty unlock date",
			msg: MsgLock{
				Address:    validAddr,
				UnlockDate: "",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(1000)),
			},
			err: ErrInvalidDate,
		},
		{
			name: "invalid unlock date format",
			msg: MsgLock{
				Address:    validAddr,
				UnlockDate: "12/01/2026",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(1000)),
			},
			err: ErrInvalidDate,
		},
		{
			name: "invalid unlock date format - not a date",
			msg: MsgLock{
				Address:    validAddr,
				UnlockDate: "not-a-date",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(1000)),
			},
			err: ErrInvalidDate,
		},
		{
			name: "invalid amount - zero",
			msg: MsgLock{
				Address:    validAddr,
				UnlockDate: "2026-12-01",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(0)),
			},
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "invalid amount - negative",
			msg: MsgLock{
				Address:    validAddr,
				UnlockDate: "2026-12-01",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(-1000)),
			},
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "valid single lock",
			msg: MsgLock{
				Address:    validAddr,
				UnlockDate: "2026-12-01",
				Amount:     sdk.NewCoin(bondDenom, math.NewInt(1000)),
			},
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
