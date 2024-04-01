package handlers_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/andreevym/metric-collector/internal/handlers"
	"github.com/andreevym/metric-collector/internal/storage/mocks"
	"github.com/golang/mock/gomock"
)

// ExampleServiceHandlers_GetPingHandler demonstrates how to use GetPingHandler.
func ExampleServiceHandlers_GetPingHandler() {
	mockCtrl := gomock.NewController(nil)
	defer mockCtrl.Finish()
	mockStorage := mocks.NewMockStorage(mockCtrl)
	mockClient := mocks.NewMockClient(mockCtrl)
	mockCtrl.RecordCall(mockClient, "Ping").Return(nil)

	// Create an instance of ServiceHandlers
	serviceHandlers := handlers.NewServiceHandlers(mockStorage, mockClient)

	// Create a new HTTP request (GET /ping)
	req, _ := http.NewRequest(http.MethodGet, handlers.PathGetPing, nil)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Create an HTTP handler function from GetPingHandler method
	handler := http.HandlerFunc(serviceHandlers.GetPingHandler)

	// Serve the HTTP request to the ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Output: 200
	fmt.Println(rr.Code)
}
