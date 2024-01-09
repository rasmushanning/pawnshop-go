package pawnshop

import (
	"pawnshop/server/pkg/messages"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		offer    messages.Offer
		rules    []offerValidationRule
		expError bool
	}{
		{
			name: "ensureProfitRule - demand < offer, should not return error",
			offer: messages.Offer{
				Offer:  2,
				Demand: 1,
			},
			rules: []offerValidationRule{
				&ensureProfitRule{},
			},
			expError: false,
		},
		{

			name: "ensureProfitRule - demand > offer, should return error",
			offer: messages.Offer{
				Offer:  2,
				Demand: 3,
			},
			rules: []offerValidationRule{
				&ensureProfitRule{},
			},
			expError: true,
		},
		{

			name: "nil rule - will not count as a rule, should not return error",
			offer: messages.Offer{
				Offer:  2,
				Demand: 3,
			},
			rules: []offerValidationRule{
				nil,
			},
			expError: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			validator := newValidator(
				c.rules...,
			)

			err := validator.validate(c.offer)
			if c.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
