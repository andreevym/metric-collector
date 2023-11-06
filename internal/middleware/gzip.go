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
		acceptEncoding := r.Header.Get("Accept-Encoding")
		isRespGzip := strings.Contains(acceptEncoding, compressor.AcceptEncoding)
		if isRespGzip {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := compressor.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer func() {
				err := cw.Close()
				if err != nil {
					logger.Log.Fatal(err.Error())
				}
			}()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		isReqGzip := strings.Contains(contentEncoding, compressor.ContentEncoding)
		if isReqGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := compressor.NewCompressReader(r.Body)
			if err != nil {
				logger.Log.Fatal(err.Error())
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer func() {
				err = cr.Close()
				if err != nil {
					logger.Log.Fatal(err.Error())
				}
			}()
		}

		h.ServeHTTP(ow, r)
	})
}
