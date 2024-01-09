package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOffer(t *testing.T) {
	cases := []struct {
		name     string
		offer    int
		demand   int
		expOffer Offer
	}{
		{
			name:   "Create offer",
			offer:  5,
			demand: 6,
			expOffer: Offer{
				Code:   "PAWN",
				Offer:  5,
				Demand: 6,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expOffer, CreateOffer(c.offer, c.demand))
		})
	}
}

func TestCreateAnswer(t *testing.T) {
	cases := []struct {
		name      string
		value     int
		expAnswer Answer
	}{
		{
			name:  "Create answer",
			value: 5,
			expAnswer: Answer{
				Code:  "ACCEPT",
				Value: 5,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expAnswer, CreateAcceptedAnswer(c.value))
		})
	}
}

func TestCreateRejectAnswer(t *testing.T) {
	cases := []struct {
		name      string
		value     int
		expAnswer Answer
	}{
		{
			name: "Create reject answer",
			expAnswer: Answer{
				Code: "REJECT",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expAnswer, CreateRejectAnswer())
		})
	}
}
