package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

const ContentEncoding = "gzip"
const AcceptEncoding = "gzip"

// GzipWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close() //NOT SUFFICIENT, DON'T DEFER WRITER OBJECTS
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	// NEED TO CLOSE EXPLICITLY
	if err := gz.Close(); err != nil {
		return nil, err
	}
	compressedData := b.Bytes()
	return compressedData, nil
}

func Decompress(compressedData []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	res, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}
