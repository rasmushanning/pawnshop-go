package inventory

import (
	"pawnshop/server/pkg/messages"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInventory(t *testing.T) {
	cases := []struct {
		name     string
		size     int
		expected *Inventory
	}{
		{
			name: "size 1",
			size: 1,
			expected: &Inventory{
				items:              []int{1},
				smallestValue:      1,
				smallestValueIndex: 0,
			},
		},
		{
			name: "size 5",
			size: 5,
			expected: &Inventory{
				items:              []int{1, 1, 1, 1, 1},
				smallestValue:      1,
				smallestValueIndex: 0,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.expected, NewInventory(c.size))
		})
	}
}

func TestHandleOffer(t *testing.T) {
	cases := []struct {
		name                     string
		offer                    messages.Offer
		expected                 messages.Answer
		oldSmallestValue         int
		oldSmallestValueIndex    int
		oldItems                 []int
		expNewItems              []int
		expNewSmallestValue      int
		expNewSmallestValueIndex int
	}{
		{
			name: "offer > smallestValue > demand, fresh inventory, should be accepted",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 0,
			},
			expected: messages.Answer{
				Code:  messages.AcceptCode,
				Value: 1,
			},
			oldSmallestValue:         1,
			oldSmallestValueIndex:    0,
			oldItems:                 []int{1, 1, 1, 1, 1},
			expNewSmallestValue:      1,
			expNewSmallestValueIndex: 1,
			expNewItems:              []int{2, 1, 1, 1, 1},
		},
		{
			name: "offer > smallestValue == demand, fresh inventory, should be accepted",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 1,
			},
			expected: messages.Answer{
				Code:  messages.AcceptCode,
				Value: 1,
			},
			oldSmallestValue:         1,
			oldSmallestValueIndex:    0,
			oldItems:                 []int{1, 1, 1, 1, 1},
			expNewSmallestValue:      1,
			expNewSmallestValueIndex: 1,
			expNewItems:              []int{2, 1, 1, 1, 1},
		},
		{
			name: "offer > smallestValue < demand, fresh inventory, should be rejected",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 3,
			},
			expected: messages.Answer{
				Code: messages.RejectCode,
			},
			oldSmallestValue:         1,
			oldSmallestValueIndex:    0,
			oldItems:                 []int{1, 1, 1, 1, 1},
			expNewSmallestValue:      1,
			expNewSmallestValueIndex: 0,
			expNewItems:              []int{1, 1, 1, 1, 1},
		},
		{
			name: "offer == smallestValue > demand, fresh inventory, should be rejected",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  1,
				Demand: 0,
			},
			expected: messages.Answer{
				Code: messages.RejectCode,
			},
			oldSmallestValue:         1,
			oldSmallestValueIndex:    0,
			oldItems:                 []int{1, 1, 1, 1, 1},
			expNewSmallestValue:      1,
			expNewSmallestValueIndex: 0,
			expNewItems:              []int{1, 1, 1, 1, 1},
		},
		{
			name: "offer < smallestValue > demand, fresh inventory, should be rejected",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  1,
				Demand: 0,
			},
			expected: messages.Answer{
				Code: messages.RejectCode,
			},
			oldSmallestValue:         2,
			oldSmallestValueIndex:    0,
			oldItems:                 []int{2, 2, 2, 2, 2},
			expNewSmallestValue:      2,
			expNewSmallestValueIndex: 0,
			expNewItems:              []int{2, 2, 2, 2, 2},
		},
		{
			name: "offer > smallestValue > demand, mixed inventory, should be accepted",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  5,
				Demand: 4,
			},
			expected: messages.Answer{
				Code:  messages.AcceptCode,
				Value: 4,
			},
			oldSmallestValue:         2,
			oldSmallestValueIndex:    3,
			oldItems:                 []int{7, 4, 5, 2, 7},
			expNewSmallestValue:      2,
			expNewSmallestValueIndex: 3,
			expNewItems:              []int{7, 5, 5, 2, 7},
		},
		{
			name: "offer > smallestValue > demand, mixed inventory, should be accepted",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  5,
				Demand: 2,
			},
			expected: messages.Answer{
				Code:  messages.AcceptCode,
				Value: 2,
			},
			oldSmallestValue:         2,
			oldSmallestValueIndex:    3,
			oldItems:                 []int{7, 4, 5, 2, 7},
			expNewSmallestValue:      4,
			expNewSmallestValueIndex: 1,
			expNewItems:              []int{7, 4, 5, 5, 7},
		},
		{
			name: "offer > smallestValue, smallestValue < demand, offer > demand, mixed inventory, should be accepted",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  150,
				Demand: 100,
			},
			expected: messages.Answer{
				Code: messages.RejectCode,
			},
			oldSmallestValue:         2,
			oldSmallestValueIndex:    3,
			oldItems:                 []int{7, 4, 5, 2, 7},
			expNewSmallestValue:      2,
			expNewSmallestValueIndex: 3,
			expNewItems:              []int{7, 4, 5, 2, 7},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Inventory{
				items:              c.oldItems,
				smallestValue:      c.oldSmallestValue,
				smallestValueIndex: c.oldSmallestValueIndex,
			}

			assert.Equal(t, c.expected, i.HandleOffer(c.offer))
			assert.Equal(t, c.expNewItems, i.items)
			assert.Equal(t, c.expNewSmallestValue, i.smallestValue)
			assert.Equal(t, c.expNewSmallestValueIndex, i.smallestValueIndex)
		})
	}
}

func TestEnsureProfitability(t *testing.T) {
	cases := []struct {
		name          string
		offer         messages.Offer
		exp           bool
		items         []int
		smallestValue int
		expIdx        int
	}{
		{
			name: "offer < smallestValue, should always return false",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 1,
			},
			smallestValue: 3,
			items:         []int{3, 3, 3, 3, 3},
			exp:           false,
		},
		{
			name: "offer > demand, offer > smallestValue, 1st item allows for maximum profit, should return true",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 1,
			},
			smallestValue: 1,
			items:         []int{1, 1, 1, 1, 1},
			expIdx:        0,
			exp:           true,
		},
		{
			name: "offer > demand, offer > smallestValue, 3rd item allows for maximum profit, should return true",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 1,
			},
			smallestValue: 1,
			items:         []int{2, 2, 1, 2, 2},
			expIdx:        2,
			exp:           true,
		},
		{
			name: "offer > negative demand, offer > smallestValue, 3rd item allows for maximum profit, should return true",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: -13,
			},
			smallestValue: 1,
			items:         []int{2, 2, 1, 2, 2},
			expIdx:        2,
			exp:           true,
		},
		{
			name: "offer == demand, offer > smallestValue, no items satisfy demand & give profit, should return false",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  2,
				Demand: 2,
			},
			smallestValue: 1,
			items:         []int{2, 2, 1, 2, 2},
			exp:           false,
		},
		{
			name: "offer > demand, offer > smallestValue, 5th item allows for maximum profit, should return true",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  5,
				Demand: 4,
			},
			smallestValue: 3,
			items:         []int{5, 3, 7, 10, 4},
			expIdx:        4,
			exp:           true,
		},
		{
			name: "offer > demand, offer > smallestValue, 5th item allows maximum profit, should return true",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  11,
				Demand: 3,
			},
			smallestValue: 3,
			items:         []int{5, 3, 3, 10, 4},
			expIdx:        1,
			exp:           true,
		},
		{
			name: "offer > demand, offer > smallestValue, no items satisfy demand & give profit, should return true",
			offer: messages.Offer{
				Code:   messages.PawnCode,
				Offer:  4,
				Demand: 2,
			},
			smallestValue: 1,
			items:         []int{1, 1, 1, 1, 1},
			exp:           false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Inventory{
				items:         c.items,
				smallestValue: c.smallestValue,
			}

			ok, idx := i.isProfitable(c.offer)
			assert.Equal(t, c.exp, ok)
			assert.Equal(t, c.expIdx, idx)
		})
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		name     string
		items    []int
		expected string
	}{
		{
			name:     "inventory with one item",
			items:    []int{1},
			expected: "[1]",
		},
		{
			name:     "inventory with multiple items",
			items:    []int{1, 2, 3},
			expected: "[1, 2, 3]",
		},
		{
			name:     "empty inventory",
			items:    []int{},
			expected: "[]",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			i := Inventory{
				items: c.items,
			}

			assert.Equal(t, c.expected, i.String())
		})
	}
}
