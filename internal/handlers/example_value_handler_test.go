package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/storage/mocks"
	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
)

// ExampleGetValueHandler demonstrates how to use GetPingHandler.
func ExampleGetValueHandler() {
	mockCtrl := gomock.NewController(nil)
	defer mockCtrl.Finish()
	mockStorage := mocks.NewMockStorage(mockCtrl)
	mockClient := mocks.NewMockClient(mockCtrl)

	value := float64(1)
	metricName := "metricName"
	mType := store.MTypeGauge
	metric := &store.Metric{
		ID:    metricName,
		MType: mType,
		Value: &value,
	}
	mockCtrl.RecordCall(mockStorage, "Read", gomock.Any(), gomock.Any(), gomock.Any()).Return(metric, nil)

	// Create an instance of ServiceHandlers
	serviceHandlers := handlers.NewServiceHandlers(mockStorage, mockClient)

	// Create a new HTTP request (GET /value/{metricType}/{metricName})
	url := fmt.Sprintf("%s/%s/%s", handlers.PathValue, mType, metricName)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create an HTTP handler function from GetPingHandler method
	r := chi.NewRouter()
	r.Get(handlers.PathValue+"/{metricType}/{metricName}", serviceHandlers.GetValueHandler)

	// Serve the HTTP request to the ResponseRecorder
	r.ServeHTTP(rr, req)

	// Output: 200
	fmt.Println(rr.Code)
}
