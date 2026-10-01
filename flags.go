package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// Flags are compiled into the binary, so turning one on ships through a PR,
// CI, and a deploy like any other change.
//
//go:embed flags.json
var flagsJSON []byte

// Flags gates customer-facing features.
type Flags struct {
	// CamouflageMode hides balances behind a shimmer until the customer taps
	// to reveal them, for checking your account on a crowded train.
	CamouflageMode bool `json:"camouflage_mode"`

	// InkCloudCardFreeze lets a customer freeze a card from the dashboard.
	InkCloudCardFreeze bool `json:"ink_cloud_card_freeze"`
}

func loadFlags() (Flags, error) {
	var f Flags
	if err := json.Unmarshal(flagsJSON, &f); err != nil {
		return Flags{}, fmt.Errorf("parsing flags.json: %w", err)
	}
	return f, nil
}
