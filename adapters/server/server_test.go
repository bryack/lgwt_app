package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bryack/lgwt_app/domain"
	"github.com/bryack/lgwt_app/store"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestGETPlayers(t *testing.T) {
	store := &testhelpers.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := mustMakePlayerServer(t, store)

	tests := []struct {
		name               string
		player             string
		expectedHTTPStatus int
		expectedScore      string
	}{
		{
			name:               "Returns Pepper's score",
			player:             "Pepper",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "20",
		},
		{
			name:               "Returns Floyd's score",
			player:             "Floyd",
			expectedHTTPStatus: http.StatusOK,
			expectedScore:      "10",
		},
		{
			name:               "Returns 404 on missing players",
			player:             "Apollo",
			expectedHTTPStatus: http.StatusNotFound,
			expectedScore:      "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := newGetScoreRequest(tt.player)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, tt.expectedHTTPStatus)
			assertResponseBody(t, response.Body.String(), tt.expectedScore)
		})
	}
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, want %q, but got %q", want, got)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, want %d, but got %d", want, got)
	}
}

func TestStoreWins(t *testing.T) {
	store := &testhelpers.StubPlayerStore{
		Scores:   map[string]int{},
		WinCalls: nil,
	}

	server := mustMakePlayerServer(t, store)
	t.Run("it returns accepted on POST", func(t *testing.T) {
		player := "Pepper"
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin, want %d", len(store.WinCalls), 1)
		}

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.WinCalls[0], player)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func TestLeague(t *testing.T) {
	wantedLeague := []domain.Player{
		{Name: "Cleo", Wins: 32},
		{Name: "Chris", Wins: 20},
		{Name: "Tiest", Wins: 14},
	}
	t.Run("it returns 200 on /league", func(t *testing.T) {
		store := &testhelpers.StubPlayerStore{Scores: nil, WinCalls: nil, League: wantedLeague, Err: nil}
		server := mustMakePlayerServer(t, store)

		request, err := newLeagueRequest(t)
		assert.NoError(t, err)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assert.Equal(t, got, wantedLeague)
		assert.Equal(t, jsonContentType, response.Result().Header.Get("content-type"), "response did not have content-type of application/json")
	})

	t.Run("handle 500", func(t *testing.T) {
		store := &testhelpers.StubPlayerStore{Scores: nil, WinCalls: nil, League: wantedLeague, Err: errors.New("database connection failed")}
		server := mustMakePlayerServer(t, store)

		request, err := newLeagueRequest(t)
		assert.NoError(t, err)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
	})
}

func newLeagueRequest(t *testing.T) (*http.Request, error) {
	t.Helper()
	return http.NewRequest(http.MethodGet, "/league", nil)
}

func getLeagueFromResponse(t *testing.T, body io.Reader) (league []domain.Player) {
	t.Helper()
	league, err := domain.NewLeague(body)
	if err != nil {
		t.Fatalf("failed to get league from response body: %v", err)
	}
	return
}

func TestGame(t *testing.T) {

	t.Run("GET /game returns 200", func(t *testing.T) {
		store := &testhelpers.StubPlayerStore{}
		server := mustMakePlayerServer(t, store)
		request, err := newGameRequest(t)
		assert.NoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
	t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &testhelpers.StubPlayerStore{}
		winner := "Ruth"
		server := httptest.NewServer(mustMakePlayerServer(t, store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
		ws := mustDialWS(t, wsURL)
		defer ws.Close()

		writeWSMessage(t, ws, winner)

		time.Sleep(10 * time.Millisecond)
		testhelpers.AssertPlayerWin(t, store, winner)
	})
}

func newGameRequest(t *testing.T) (*http.Request, error) {
	t.Helper()
	return http.NewRequest(http.MethodGet, "/game", nil)
}

func mustMakePlayerServer(t *testing.T, store store.PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("failed to create player server", err)
	}
	return server
}

func mustDialWS(t *testing.T, wsURL string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to open a ws connection on %q: %v", wsURL, err)
	}
	return ws
}

func writeWSMessage(t testing.TB, ws *websocket.Conn, winner string) {
	t.Helper()
	if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
		t.Fatalf("failed to send message over ws connection: %v", err)
	}
}
