package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/forderation/go-app/http_app"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}

	server := &http_app.PlayerServer{Store: &store}

	t.Run("return the Pepper score", func(t *testing.T) {
		request := NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "20"
		AssertResponseBody(t, got, want)
		AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("return the Floyd score", func(t *testing.T) {
		request := NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := response.Body.String()
		want := "10"
		AssertResponseBody(t, got, want)
		AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		scores:   map[string]int{},
		winCalls: nil,
	}
	server := &http_app.PlayerServer{Store: &store}

	t.Run("it returns accepted on post", func(t *testing.T) {
		player := "Pepper"
		request := NewPostWinRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to recorded want %d", len(store.winCalls), 1)
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
	})
}

func NewGetScoreRequest(Name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", Name), nil)
	return req
}

func NewPostWinRequest(Name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", Name), nil)
	return req
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong got %q, want %q", got, want)
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not correct status got %d, want %d", got, want)
	}
}
