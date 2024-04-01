// Package compressor provides utilities for compressing and decompressing data using gzip encoding.
package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
)

// ContentEncoding represents the content encoding type for gzip compression.
const ContentEncoding = "gzip"

// AcceptEncoding represents the accepted encoding type for gzip compression.
const AcceptEncoding = "gzip"

// GzipWriter implements the http.ResponseWriter interface and allows transparent
// compression of data being transmitted by the server, while setting the correct HTTP headers.
type GzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write writes compressed data to the underlying writer.
func (w GzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Compress compresses the given data using gzip encoding.
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decompress decompresses the given gzip compressed data.
func Decompress(compressedData []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
