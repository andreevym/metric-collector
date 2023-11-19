package middleware

import (
	"net/http"
	"strings"

	"github.com/andreevym/metric-collector/internal/compressor"
	"github.com/andreevym/metric-collector/internal/logger"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		contentType := r.Header.Get("Content-Type")
		isCompressedContentType := strings.Contains(contentType, "application/json") ||
			strings.Contains(contentType, "text/html")
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isRespGzip := strings.Contains(acceptEncoding, compressor.AcceptEncoding)
		contentEncoding := r.Header.Get("Content-Encoding")
		isReqGzip := strings.Contains(contentEncoding, compressor.ContentEncoding)
		if isCompressedContentType && isRespGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := compressor.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer func() {
				err := cw.Close()
				if err != nil {
					logger.Log.Error(err.Error())
				}
			}()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		if isCompressedContentType && isReqGzip {
			compressReader, err := compressor.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}

			r.Body = compressReader

			defer func() {
				err := r.Body.Close()
				if err != nil {
					logger.Log.Error(err.Error())
				}
			}()
		}

		h.ServeHTTP(ow, r)
	})
}
