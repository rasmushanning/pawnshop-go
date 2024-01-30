package pawnshop

import (
	"errors"
	"pawnshop/server/pkg/messages"
)

/*
offerValidationRule is an interface for a rule that can validate an offer.
*/
type offerValidationRule interface {
	validate(o messages.Offer) error
}

/*
validator is a struct that contains a list of offerValidationRules.
*/
type validator struct {
	rules []offerValidationRule
}

/*
Creates a new validator with the given rules.
*/
func newValidator(rules ...offerValidationRule) (*validator, error) {
	for _, rule := range rules {
		if rule == nil {
			return nil, errors.New("validator rules can not be nil")
		}
	}

	return &validator{
		rules: rules,
	}, nil
}

/*
Validates an offer with the validator's rules.
*/
func (v validator) validate(o messages.Offer) error {
	for _, rule := range v.rules {
		if err := rule.validate(o); err != nil {
			return err
		}
	}
	return nil
}
