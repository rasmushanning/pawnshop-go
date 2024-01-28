package inventory

import (
	"fmt"
	"math"
	"pawnshop/server/pkg/messages"

	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const (
	defaultItemValue = 1
)

/*
Inventory is a data structure that manages a list of items in a thread safe manner.
It also provides a function for handling offers from clients, accepting them if they are profitable
and rejecting them if they are not.
*/
type Inventory struct {
	items              []int
	smallestValue      int
	smallestValueIndex int
	lock               sync.Mutex
}

/*
Creates a new inventory with the given size.
*/
func NewInventory(size int) *Inventory {
	items := make([]int, size)

	for i := range items {
		items[i] = defaultItemValue
	}

	return &Inventory{
		items:              items,
		smallestValue:      defaultItemValue,
		smallestValueIndex: 0,
		lock:               sync.Mutex{},
	}
}

/*
Handles an offer from the caller. It checks if the offer would be profitable
for the inventory, and if so, it will allow the offer and return the exchanged value.
If the offer does not align with the inventory's requirements, it will reject the offer.
*/
func (i *Inventory) HandleOffer(o messages.Offer) messages.Answer {
	defer log.Infof("Inventory after handling offer: %s", i)
	i.lock.Lock()
	defer i.lock.Unlock()

	isP, idx := i.isProfitable(o)
	if !isP {
		log.Debugf("Offer %+v is not profitable for the inventory, or not possible for the inventory to accept", o)
		return messages.CreateRejectAnswer()
	}

	// Hold value to return to client in answer
	valToRet := i.items[idx]

	// Replace the item in the inventory that was decided to be the most
	// profitable to give up, with the received offer
	i.items[idx] = o.Offer

	// If the replaced object was the smallest value in the inventory,
	// we find the new smallest value and its index to avoid the work of iteratively
	// finding the smallest value for every new offer
	if idx == i.smallestValueIndex {
		newSmValue := i.items[0]
		newSmValueIdx := 0
		for j := 1; j < len(i.items); j++ {
			if i.items[j] < newSmValue {
				newSmValue = i.items[j]
				newSmValueIdx = j
			}
		}
		// Cache the new smallest value and its index
		i.smallestValue = newSmValue
		i.smallestValueIndex = newSmValueIdx
	}

	return messages.CreateAcceptedAnswer(valToRet)
}

/*
Returns a string representation of the inventory.
*/
func (i *Inventory) String() string {
	i.lock.Lock()
	defer i.lock.Unlock()

	s := make([]string, len(i.items))
	for j, item := range i.items {
		s[j] = strconv.Itoa(item)
	}

	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}

/*
Checks if the offer is possible and can be profitable for the inventory.
If it is profitable, it will also return the index of the most profitable item in the inventory
that satisfies the demand. It is NOT thread-safe and should be called from another thread-safe
function in the inventory.
*/
func (i *Inventory) isProfitable(o messages.Offer) (bool, int) {
	// If the offer is less than or equal to the smallest value in the inventory, it can not be profitable
	if o.Offer <= i.smallestValue {
		return false, 0
	}

	// Find the item in the inventory that gives the highest profit.
	// It must be greater than or equal to the demand, less than the offer,
	// and provide a greater profit than any other item in the inventory.
	maxPrItemVal := math.MaxInt
	maxPrItemIdx := -1
	for idx := 0; idx < len(i.items); idx++ {
		item := i.items[idx]
		if item >= o.Demand && // Satisfies demand
			o.Offer > item && // Will ensure profit
			item < maxPrItemVal { // Is it the best item so far?
			maxPrItemVal = item
			maxPrItemIdx = idx
		}
	}

	// If no item satisfies the demand and gives profit, the offer is not profitable
	if maxPrItemIdx == -1 {
		return false, 0
	}

	return true, maxPrItemIdx
}
