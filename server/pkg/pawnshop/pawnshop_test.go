package pawnshop

import (
	"pawnshop/server/pkg/messages"
	"pawnshop/server/pkg/mocks"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleOffer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockOfferHandler := mocks.NewMockOfferHandler(ctrl)

	cases := []struct {
		name         string
		offer        messages.Offer
		expected     messages.Answer
		expectations func()
	}{
		{
			name: "Offer is accepted, should return ACCEPT",
			offer: messages.Offer{
				Code:   "PAWN",
				Offer:  2,
				Demand: 1,
			},
			expected: messages.Answer{
				Code:  messages.AcceptCode,
				Value: 1,
			},
			expectations: func() {
				mockOfferHandler.EXPECT().HandleOffer(messages.Offer{
					Code:   "PAWN",
					Offer:  2,
					Demand: 1,
				}).Return(messages.Answer{
					Code:  messages.AcceptCode,
					Value: 1,
				}).Times(1)
				mockOfferHandler.EXPECT().String().Return("[2]").Times(2) // Used for logging inventory
			},
		},
		{
			name: "Offer is rejected, should return REJECT",
			offer: messages.Offer{
				Code:   "PAWN",
				Offer:  2,
				Demand: 5,
			},
			expected: messages.Answer{
				Code: messages.RejectCode,
			},
			expectations: func() {
				mockOfferHandler.EXPECT().String().Return("[]").Times(1) // Used for logging inventory
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.expectations()

			shop, err := NewPawnShop(mockOfferHandler)
			require.NoError(t, err)
			require.Equal(t, c.expected, shop.HandleOffer(c.offer))
		})
	}
}
