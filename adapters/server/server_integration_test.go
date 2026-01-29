package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bryack/lgwt_app/domain"
	"github.com/bryack/lgwt_app/filesystem"
	"github.com/bryack/lgwt_app/testhelpers"
	"github.com/stretchr/testify/assert"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	db, cleanDatabase := testhelpers.CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := filesystem.NewFileSystemPlayerStore(db)
	if err != nil {
		t.Fatalf("failed to create file system player store: %v", err)
	}
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetScoreRequest(player))

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "3")
	})
	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, err := newLeagueRequest(t)
		assert.NoError(t, err)

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)

		got := getLeagueFromResponse(t, response.Body)
		want := []domain.Player{
			{Name: "Pepper", Wins: 3},
		}

		assert.Equal(t, want, got)
	})
}
