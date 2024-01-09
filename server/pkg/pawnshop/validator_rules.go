package pawnshop

import (
	"errors"
	"pawnshop/server/pkg/messages"
)

/*
ensureProfitRule is a rule that ensures that the offer is greater than the demand.
*/
type ensureProfitRule struct{}

/*
validate validates an offer with the ensureProfitRule.
*/
func (e *ensureProfitRule) validate(o messages.Offer) error {
	if o.Offer <= o.Demand {
		return errors.New("offer must be greater than demand")
	}

	return nil
}
