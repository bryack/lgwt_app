package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bryack/lgwt_app/domain"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestGETPlayers(t *testing.T) {
	store := testhelpers.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := NewPlayerServer(&store)

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
	store := testhelpers.StubPlayerStore{
		Scores:   map[string]int{},
		WinCalls: nil,
	}

	server := NewPlayerServer(&store)
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
		store := testhelpers.StubPlayerStore{Scores: nil, WinCalls: nil, League: wantedLeague, Err: nil}
		server := NewPlayerServer(&store)

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
		store := testhelpers.StubPlayerStore{Scores: nil, WinCalls: nil, League: wantedLeague, Err: errors.New("database connection failed")}
		server := NewPlayerServer(&store)

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
		store := testhelpers.StubPlayerStore{}
		server := NewPlayerServer(&store)
		request, err := newGameRequest(t)
		assert.NoError(t, err)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func newGameRequest(t *testing.T) (*http.Request, error) {
	t.Helper()
	return http.NewRequest(http.MethodGet, "/game", nil)
}
