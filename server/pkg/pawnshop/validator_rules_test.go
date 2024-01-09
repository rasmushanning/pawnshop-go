package pawnshop

import (
	"pawnshop/server/pkg/messages"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnsureProfitRuleValidate(t *testing.T) {
	cases := []struct {
		name     string
		offer    messages.Offer
		expError bool
	}{
		{
			name: "demand < offer, should not return error",
			offer: messages.Offer{
				Offer:  2,
				Demand: 1,
			},
			expError: false,
		},
		{
			name: "demand > offer, should return error",
			offer: messages.Offer{
				Offer:  2,
				Demand: 3,
			},
			expError: true,
		},
		{
			name: "demand == offer, should return error",
			offer: messages.Offer{
				Offer:  2,
				Demand: 2,
			},
			expError: true,
		},
		{
			name: "Offer is missing (defaults to 0), demand is > 0, should return error",
			offer: messages.Offer{
				Demand: 1,
			},
			expError: true,
		},
		{
			name: "Demand is missing (defaults to 0), offer is > 0, should not return error",
			offer: messages.Offer{
				Offer: 1,
			},
			expError: false,
		},
		{
			name:     "Demand AND offer are missing (defaults to 0), offer == demand, should return error",
			offer:    messages.Offer{},
			expError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			epr := ensureProfitRule{}

			err := epr.validate(c.offer)
			if c.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
