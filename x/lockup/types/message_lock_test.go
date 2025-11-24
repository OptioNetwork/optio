package types

import (
	"testing"

	"github.com/OptioNetwork/optio/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgLock_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgLock
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgLock{
				LockupAddress: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgLock{
				LockupAddress: sample.AccAddress(),
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
