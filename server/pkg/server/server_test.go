package server

import (
	"encoding/json"
	"fmt"
	"net"
	"pawnshop/server/pkg/messages"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	invSz := 2
	s := startServerAndWait(t, invSz)
	defer func() {
		err := s.Stop()
		require.NoError(t, err)
	}()

	cases := []struct {
		name        string
		addr        string
		offerString string
		expAnswer   messages.Answer
	}{
		{
			name: "Accepted offer",
			offerString: `{
				"code": "PAWN",
				"offer": 5,
				"demand": -2
				}`,
			expAnswer: messages.Answer{
				Code:  messages.AcceptCode,
				Value: 1,
			},
		},
		{
			name: "Rejected offer",
			offerString: `
				"code": "PAWN",
				"offer": 5,
				"demand": 6
				}`,
			expAnswer: messages.Answer{
				Code: messages.RejectCode,
			},
		},
		{
			name: "Unsupported offer",
			offerString: `{
			"code": "unsupported code",
			"offer": 5,
			"demand": -2
			}`,
			expAnswer: messages.Answer{
				Code: messages.RejectCode,
			},
		},
		{
			name:        "Non-JSON body",
			offerString: `not a JSON body`,
			expAnswer: messages.Answer{
				Code: messages.RejectCode,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			conn, err := net.Dial("tcp", s.addr)
			require.NoError(t, err)

			_, err = conn.Write([]byte(c.offerString))
			require.NoError(t, err)

			buf := make([]byte, 128)

			n, err := conn.Read(buf)
			require.NoError(t, err)

			var answer messages.Answer
			err = json.Unmarshal(buf[:n], &answer)
			require.NoError(t, err)

			require.Equal(t, c.expAnswer, answer)
		})
	}
}

func TestServerErrors(t *testing.T) {
	cases := []struct {
		name          string
		addr          string
		size          int
		expNewError   bool
		expStartError bool
	}{
		{
			name:          "Invalid size",
			size:          0,
			expNewError:   true,
			expStartError: false,
		},
		{
			name:          "Invalid address",
			size:          1,
			addr:          "invalid",
			expNewError:   false,
			expStartError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := NewPawnShopServer(c.size)
			if c.expNewError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			s.addr = c.addr

			err = s.Start()
			if c.expStartError {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			err = s.Stop()
			require.NoError(t, err)
		})
	}
}

func startServerAndWait(t *testing.T, size int) *PawnShopServer {
	s, err := NewPawnShopServer(size)
	require.NoError(t, err)

	availablePort, err := getAvailablePort()
	require.NoError(t, err)
	s.addr = fmt.Sprintf("127.0.0.1:%d", availablePort)

	go func() {
		err = s.Start()
		require.NoError(t, err)
	}()

	// Wait for server to start
	for i := 0; i < 40; i++ {
		if s.IsRunning() {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}
	require.True(t, s.IsRunning())

	return s
}

func getAvailablePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
