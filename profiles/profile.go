package main

import (
	context2 "context"
	"github.com/andreevym/metric-collector/internal/transport/http/handlers"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/andreevym/metric-collector/internal/storage/store"
	"github.com/go-chi/chi/v5"
)

const (
	requestURL        = "http://localhost:8080/update/counter/test/1"
	resultFilenameCPU = "profiles/result_cpu.pprof"
	resultFilenameMEM = "profiles/result_mem.pprof"
)

// profile runs the specified handler function and profiles it
func profile(handler func(*http.Request) (*store.Metric, error), cpuFilename, memFilename string) {
	// Prepare test data
	r, err := prepareTestData()
	if err != nil {
		panic(err)
	}

	// Open CPU profile file
	fcpu, err := os.Create(cpuFilename)
	if err != nil {
		panic(err)
	}
	defer fcpu.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(fcpu); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	// Perform garbage collection
	runtime.GC()

	// Execute the handler function
	if _, err := handler(r); err != nil {
		panic(err)
	}

	// Open memory profile file
	fmem, err := os.Create(memFilename)
	if err != nil {
		panic(err)
	}
	defer fmem.Close()

	// Write memory profile
	if err := pprof.WriteHeapProfile(fmem); err != nil {
		panic(err)
	}
}

// prepareTestData prepares a mock HTTP request
func prepareTestData() (*http.Request, error) {
	router := chi.NewRouter()
	router.Post(handlers.PathPostUpdate+"/{metricType}/{metricName}/{metricValue}", nil)

	// Simulate URL params
	params := chi.RouteParams{
		Keys:   []string{"metricType", "metricName", "metricValue"},
		Values: []string{"counter", "test", "1"},
	}

	context := chi.NewRouteContext()
	context.Routes = router
	context.URLParams = params
	context.RouteMethod = "POST"
	context.RoutePath = "/update/{metricType}/{metricName}/{metricValue}"

	request, err := http.NewRequest("POST", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Add context to request
	r := request.WithContext(context2.WithValue(context2.Background(), chi.RouteCtxKey, context))
	return r, nil
}

func main() {
	// Profile BuildMetricByChiParam
	profile(handlers.BuildMetricByChiParam, "profiles/base_cpu.pprof", "profiles/base_mem.pprof")

	// Profile BuildMetricBySplitParam
	profile(handlers.BuildMetricBySplitParam, "profiles/result_cpu.pprof", "profiles/result_mem.pprof")
}
