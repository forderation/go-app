package main

import (
	"net/http/httptest"
	"testing"

	"github.com/forderation/go-app/http_app"
)

func BenchmarkRecordingWinsAndRetrievingThem(b *testing.B) {
	store := NewInMemoryPayerStore()
	server := &http_app.PlayerServer{Store: store}
	player := "Pepper"

	for i := 0; i < b.N; i++ {
		go func() {
			server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
		}()
	}

	response := httptest.NewRecorder()
	server.ServeHTTP(response, NewGetScoreRequest(player))

	// AssertStatus(t, response.Code, http.StatusOK)
	// AssertResponseBody(t, response.Body.String(), fmt.Sprintf("%d", b.N))
}
