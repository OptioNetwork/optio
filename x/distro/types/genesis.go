package types

import (
	"encoding/json"
	"fmt"
	"os"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	dailyDistributionTotals := make(map[string]uint64)
	data, err := os.ReadFile("app/upgrades/v2/daily_distribution_totals.json")
	if err == nil {
		var jsonDailyDistributionTotals []struct {
			Date   string  `json:"date"`
			Amount float64 `json:"amount"`
		}

		if err := json.Unmarshal(data, &jsonDailyDistributionTotals); err != nil {
			panic(fmt.Errorf("failed to unmarshal daily distribution totals: %w", err))
		}

		for _, dailyDistributionTotal := range jsonDailyDistributionTotals {
			dailyDistributionTotals[dailyDistributionTotal.Date] = uint64(dailyDistributionTotal.Amount)
		}
	}

	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
