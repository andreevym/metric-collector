package internation_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	t.Skip("todo: need to implement integration test")
	t.Run("[good] counter - valid url", func(t *testing.T) {
		response, err := http.Post(
			"http://localhost:8080/update/counter/someMetric/527",
			"text/plain",
			nil,
		)

		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				require.NoError(t, err)
			}
		}(response.Body)

		require.NotNil(t, response)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("[good] gauge - valid url", func(t *testing.T) {
		response, err := http.Post(
			"http://localhost:8080/update/gauge/someMetric/527",
			"text/plain",
			nil,
		)

		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				require.NoError(t, err)
			}
		}(response.Body)

		require.NotNil(t, response)
		require.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("[bad] wrong metric type", func(t *testing.T) {
		response, err := http.Post(
			"http://localhost:8080/update/UNKNOWN/someMetric/527",
			"text/plain",
			nil,
		)

		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				require.NoError(t, err)
			}
		}(response.Body)

		require.NotNil(t, response)
		require.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
	t.Run("[bad] metric value didn't set", func(t *testing.T) {
		response, err := http.Post(
			"http://localhost:8080/update/counter/someMetric",
			"text/plain",
			nil,
		)

		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				require.NoError(t, err)
			}
		}(response.Body)

		require.NotNil(t, response)
		require.Equal(t, http.StatusNotFound, response.StatusCode)
	})
	t.Run("[bad] metric value didn't set", func(t *testing.T) {
		response, err := http.Post(
			"http://localhost:8080/update/a",
			"text/plain",
			nil,
		)

		require.NoError(t, err)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				require.NoError(t, err)
			}
		}(response.Body)

		require.NotNil(t, response)
		require.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}
