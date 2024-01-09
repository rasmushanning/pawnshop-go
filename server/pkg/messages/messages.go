package messages

const (
	PawnCode        = "PAWN"
	RejectCode      = "REJECT"
	AcceptCode      = "ACCEPT"
	UnsupportedCode = "UNSUPPORTED"
)

/*
Offer is a struct that represents an offer.
*/
type Offer struct {
	Code   string `json:"code"`
	Offer  int    `json:"offer"`
	Demand int    `json:"demand"`
}

/*
Creates a new Offer with the given offer and demand.
*/
func CreateOffer(offer int, demand int) Offer {
	return Offer{
		Code:   PawnCode,
		Offer:  offer,
		Demand: demand,
	}
}

/*
Answer is a struct that represents an answer.
*/
type Answer struct {
	Code  string `json:"code"`
	Value int    `json:"value,omitempty"`
}

/*
Creates a new Answer with the given value.
*/
func CreateAcceptedAnswer(value int) Answer {
	return Answer{
		Code:  AcceptCode,
		Value: value,
	}
}

/*
Creates a new Answer with the RejectCode.
*/
func CreateRejectAnswer() Answer {
	return Answer{
		Code: RejectCode,
	}
}
