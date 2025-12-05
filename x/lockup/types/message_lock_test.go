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
				LockupAddress: "invalid_address",
				Locks: []*Lock{
					{
						UnlockDate: "2026-12-01",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
				},
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty locks array",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks:         []*Lock{},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "nil locks array",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks:         nil,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "nil lock in array",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					nil,
				},
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty unlock date",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
				},
			},
			err: ErrInvalidDate,
		},
		{
			name: "invalid unlock date format",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "12/01/2026",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
				},
			},
			err: ErrInvalidDate,
		},
		{
			name: "invalid unlock date format - not a date",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "not-a-date",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
				},
			},
			err: ErrInvalidDate,
		},
		{
			name: "invalid coin - zero amount",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "2026-12-01",
						Coin:       sdk.NewInt64Coin("uOPT", 0),
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "invalid coin - negative amount",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "2026-12-01",
						Coin:       sdk.Coin{Denom: "uOPT", Amount: math.NewInt(-1000)},
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "invalid coin - invalid denom",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "2026-12-01",
						Coin:       sdk.Coin{Denom: "INVALID!", Amount: math.NewInt(1000)},
					},
				},
			},
			err: sdkerrors.ErrInvalidCoins,
		},
		{
			name: "valid single lock",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "2026-12-01",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
				},
			},
		},
		{
			name: "valid multiple locks",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "2026-06-01",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
					{
						UnlockDate: "2026-12-01",
						Coin:       sdk.NewInt64Coin("uOPT", 2000),
					},
					{
						UnlockDate: "2027-06-01",
						Coin:       sdk.NewInt64Coin("uOPT", 3000),
					},
				},
			},
		},
		{
			name: "second lock invalid - should catch it",
			msg: MsgLock{
				LockupAddress: validAddr,
				Locks: []*Lock{
					{
						UnlockDate: "2026-06-01",
						Coin:       sdk.NewInt64Coin("uOPT", 1000),
					},
					{
						UnlockDate: "invalid-date",
						Coin:       sdk.NewInt64Coin("uOPT", 2000),
					},
				},
			},
			err: ErrInvalidDate,
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
