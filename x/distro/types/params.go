package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyMintingAddress            = []byte("MintingAddress")
	DefaultMintingAddress string = ""
)

var (
	KeyReceivingAddress            = []byte("ReceivingAddress")
	DefaultReceivingAddress string = ""
)

var (
	KeyDenom            = []byte("Denom")
	DefaultDenom string = "uOPT"
)

var (
	KeyMaxSupply            = []byte("MaxSupply")
	DefaultMaxSupply uint64 = 30000000000000000
)

var (
	KeyDistributionStartDate     = []byte("DistributionStartDate")
	DefaultDistributionStartDate = "2024-09-15"
)

var (
	KeyMonthsInHalvingPeriod            = []byte("MonthsInHalvingPeriod")
	DefaultMonthsInHalvingPeriod uint64 = 12
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	mintingAddress string,
	receivingAddress string,
	denom string,
	maxSupply uint64,
	distributionStartDate string,
	monthsInHalvingPeriod uint64,
) Params {
	return Params{
		MintingAddress:        mintingAddress,
		ReceivingAddress:      receivingAddress,
		Denom:                 denom,
		MaxSupply:             maxSupply,
		DistributionStartDate: distributionStartDate,
		MonthsInHalvingPeriod: monthsInHalvingPeriod,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultMintingAddress,
		DefaultReceivingAddress,
		DefaultDenom,
		DefaultMaxSupply,
		DefaultDistributionStartDate,
		DefaultMonthsInHalvingPeriod,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintingAddress, &p.MintingAddress, validateMintingAddress),
		paramtypes.NewParamSetPair(KeyReceivingAddress, &p.ReceivingAddress, validateReceivingAddress),
		paramtypes.NewParamSetPair(KeyDenom, &p.Denom, validateDenom),
		paramtypes.NewParamSetPair(KeyMaxSupply, &p.MaxSupply, validateMaxSupply),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateMintingAddress(p.MintingAddress); err != nil {
		return err
	}

	if err := validateReceivingAddress(p.ReceivingAddress); err != nil {
		return err
	}

	if err := validateDenom(p.Denom); err != nil {
		return err
	}

	if err := validateMaxSupply(p.MaxSupply); err != nil {
		return err
	}

	return nil
}

// validateMintingAddress validates the MintingAddress param
func validateMintingAddress(v interface{}) error {
	mintingAddress, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	_, err := sdk.AccAddressFromBech32(mintingAddress)
	if err != nil {
		return fmt.Errorf("invalid account address: %s", mintingAddress)
	}
	return nil
}

// validateReceivingAddress validates the ReceivingAddress param
func validateReceivingAddress(v interface{}) error {
	receivingAddress, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	_, err := sdk.AccAddressFromBech32(receivingAddress)
	if err != nil {
		return fmt.Errorf("invalid account address: %s", receivingAddress)
	}
	return nil
}

// validateDenom validates the Denom param
func validateDenom(v interface{}) error {
	denom, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if denom == "" {
		return fmt.Errorf("denom cannot be empty")
	}

	return nil
}

// validateMaxSupply validates the MaxSupply param
func validateMaxSupply(v interface{}) error {
	maxSupply, ok := v.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	if maxSupply <= 0 {
		return fmt.Errorf("max supply must be positive")
	}

	return nil
}
