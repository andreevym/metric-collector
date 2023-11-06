package middleware

import (
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/logger"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isRespGzip := strings.Contains(acceptEncoding, "gzip")
		if isRespGzip {
			cw := compressor.NewCompressWriter(w)
			w = cw
			err := cw.Close()
			if err != nil {
				logger.Log.Fatal(err.Error())
			}
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		isReqGzip := strings.Contains(contentEncoding, "gzip")
		if isReqGzip {
			decompressedBody, err := compressor.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Fatal(err.Error())
			}
			r.Body = decompressedBody
			err = decompressedBody.Close()
			if err != nil {
				logger.Log.Fatal(err.Error())
			}
		}

		h.ServeHTTP(w, r)
	})
}
