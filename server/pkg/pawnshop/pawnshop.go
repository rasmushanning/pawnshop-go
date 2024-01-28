package pawnshop

import (
	"fmt"
	"pawnshop/server/pkg/messages"

	log "github.com/sirupsen/logrus"
)

/*
offerHandler is an interface for an offerHandler that can handle offers as well as
return a string representation of itself.
*/
type offerHandler interface {
	HandleOffer(o messages.Offer) messages.Answer
	fmt.Stringer
}

/*
offerValidator is an interface for a validator that can validate offers.
*/
type offerValidator interface {
	validate(o messages.Offer) error
}

/*
PawnShop is a pawn shop that handles offers from callers and has a backing inventory
and offer validator.
*/
type PawnShop struct {
	inventory offerHandler
	validator offerValidator
}

/*
Creates a new PawnShop with the given inventory and an offer validator.
*/
func NewPawnShop(inv offerHandler) *PawnShop {
	val := newValidator(
		&ensureProfitRule{},
	)

	return &PawnShop{
		inventory: inv,
		validator: val,
	}
}

/*
Handles an offer from a client. It checks if the offer is valid and sane,
and if so, it forwards the offer to the inventory.
*/
func (p *PawnShop) HandleOffer(offer messages.Offer) messages.Answer {
	log.Infof("Inventory before handling offer: %s", p.inventory)

	if err := p.validator.validate(offer); err != nil {
		log.Debugf("Offer %+v is not valid: %s", offer, err)
		log.Infof("Inventory after handling offer: %s", p.inventory)
		return messages.CreateRejectAnswer()
	}

	return p.inventory.HandleOffer(offer)
}
