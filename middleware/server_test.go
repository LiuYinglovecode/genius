package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerStart(t *testing.T) {
	config := &Config{}
	server, err := NewServer(config)
	assert.Nil(t, err)
	assert.NotNil(t, server)

	t.Run("test ping", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/1/ping", nil)
		server.Route.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}
