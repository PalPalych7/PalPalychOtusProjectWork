package internalhttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
	"github.com/PalPalych7/OtusProjectWork/mocks"
	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		m := new(mocks.ServerInterface)
		m.On("Serve").Return(nil)
		err := m.Serve()
		require.NoError(t, err)

		m.On("Stop").Return(nil)
		err = m.Stop()
		require.NoError(t, err)
	})
}

func TestHandler(t *testing.T) {
	cases := []struct {
		name         string
		target       string
		body         io.Reader
		responseCode int
	}{
		{"AddBannerSlot", "/AddBannerSlot", nil, http.StatusInternalServerError},
		{"AddBannerSlot2", "/AddBannerSlot", bytes.NewBufferString(""), http.StatusInternalServerError},
		{"BannerClick", "/BannerClick", nil, http.StatusInternalServerError},
		{"DelBannerSlot", "/DelBannerSlot", nil, http.StatusInternalServerError},
		{"GetBannerForSlot", "/GetBannerForSlot", nil, http.StatusInternalServerError},
	}

	myLogger := logger.New("", "")
	service := NewServer(context.Background(), nil, ms.HTTPConf{}, myLogger)
	w := httptest.NewRecorder()
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, c.target, c.body)
			switch c.target {
			case "/AddBannerSlot":
				service.AddBannerSlot(w, r)
			case "/BannerClick":
				service.BannerClick(w, r)
			case "/DelBannerSlot":
				service.DelBannerSlot(w, r)
			case "/GetBannerForSlot":
				service.GetBannerForSlot(w, r)
			}
			resp := w.Result()
			require.Equal(t, c.responseCode, resp.StatusCode)
		})
	}
}
