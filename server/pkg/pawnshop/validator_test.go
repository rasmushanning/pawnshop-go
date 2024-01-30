package pawnshop

import (
	"pawnshop/server/pkg/messages"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	cases := []struct {
		name     string
		offer    messages.Offer
		rules    []offerValidationRule
		expError bool
	}{
		{
			name: "ensureProfitRule - should not return error",
			rules: []offerValidationRule{
				&ensureProfitRule{},
			},
			expError: false,
		},
		{

			name: "nil rule - should fail creating validator - should return error",
			rules: []offerValidationRule{
				nil,
			},
			expError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			validator, err := newValidator(
				c.rules...,
			)
			if c.expError {
				require.Nil(t, validator)
				require.Error(t, err)
			} else {
				require.NotNil(t, validator)
				require.NoError(t, err)
			}
		})
	}
}

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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			validator, err := newValidator(
				c.rules...,
			)
			require.NoError(t, err)

			err = validator.validate(c.offer)
			if c.expError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
