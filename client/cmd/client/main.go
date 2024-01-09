package main

import (
	"flag"
	"fmt"

	"pawnshop/client/pkg/client"
	"pawnshop/server/pkg/messages"
)

/*
Runs a lightweight client used to test the pawn shop server.
It accepts two flags: offer and demand, which are the offer and demand values
which will be used in the offer sent to the server.
*/
func main() {
	offer := flag.Int("offer", 0, "offer")
	demand := flag.Int("demand", 0, "demand")
	flag.Parse()

	client := &client.Client{}
	err := client.Run(
		messages.CreateOffer(
			*offer,
			*demand,
		),
	)
	if err != nil {
		fmt.Println("Client failed to run: ", err)
	}
}
