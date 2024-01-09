package acceptancetests

import (
	"pawnshop/client/pkg/client"
	"pawnshop/server/pkg/messages"
	"pawnshop/server/pkg/server"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

/*
The end to end tests test a sequential flow of offers sent from a client to a server,
and prints the results to stdout. It begins with an inventory size of 5.
*/
func TestSequentialEndToEndTests(t *testing.T) {
	invSz := 5
	srv, err := server.New(invSz)
	require.NoError(t, err)

	go func() {
		err := srv.Start()
		require.NoError(t, err)
	}()

	// Wait for server to start
	for i := 0; i < 40; i++ {
		time.Sleep(5 * time.Millisecond)
		if srv.IsRunning() {
			break
		}
	}
	require.True(t, srv.IsRunning())

	defer srv.Stop()

	// A client sends an offer {"offer": 7, "demand": 1}
	// The server accepts it and should now have an inventory of [7, 1, 1, 1, 1]
	t.Run("First offer - ACCEPT", func(t *testing.T) {
		client := &client.Client{}
		err := client.Run(messages.CreateOffer(7, 1))
		require.NoError(t, err)
	})

	// A client start sending an offer {"offer": 5, "demand": 3}
	// The pawnshop server will not accept it because it would infer a loss of 2.
	// The inventory of the server should still be [7, 1, 1, 1, 1]
	t.Run("Second offer - REJECT", func(t *testing.T) {
		client := &client.Client{}
		err := client.Run(messages.CreateOffer(5, 3))
		require.NoError(t, err)
	})

	// A client sends an offer {"offer": 4, "demand": 1}
	// The server accepts it and should now have an inventory of [7, 4, 1, 1, 1]
	t.Run("Third offer - ACCEPT", func(t *testing.T) {
		client := &client.Client{}
		err := client.Run(messages.CreateOffer(4, 1))
		require.NoError(t, err)
	})

	// A client sends an offer {"offer": 25, "demand": 8}
	// The pawnshop server can not accept it.
	// The inventory of the server should still be [7, 4, 1, 1, 1]
	t.Run("Fourth offer - REJECT", func(t *testing.T) {
		client := &client.Client{}
		err := client.Run(messages.CreateOffer(25, 8))
		require.NoError(t, err)
	})
}
