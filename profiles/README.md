# Benchmark for profile

### Bench Split param CPU + MEM

**command**

```shell
go test -bench=BenchmarkBuildMetricByParam -cpuprofile=profiles/cpu_split.out -memprofile=profiles/mem_split.out ./internal/handlers/
```

**output**

```shell
goos: linux
goarch: amd64
pkg: github.com/andreevym/metric-collector/internal/handlers
cpu: Intel(R) Core(TM) i7-10870H CPU @ 2.20GHz
BenchmarkBuildMetricByParam-16           1000000              1439 ns/op             136 B/op          3 allocs/op
PASS
ok      github.com/andreevym/metric-collector/internal/handlers 1.616s
```

### Bench CHI param CPU + MEM

**command**

```shell
go test -bench=BenchmarkBuildMetricByChiParam -cpuprofile=profiles/cpu_chi.out -memprofile=profiles/mem_chi.out ./internal/handlers/
```

**output**

```shell
goos: linux
goarch: amd64
pkg: github.com/andreevym/metric-collector/internal/handlers
cpu: Intel(R) Core(TM) i7-10870H CPU @ 2.20GHz
BenchmarkBuildMetricByChiParam-16        1000000              1088 ns/op              56 B/op          2 allocs/op
PASS
ok      github.com/andreevym/metric-collector/internal/handlers 1.214s
```

### diff cpu

**command**

```shell
go tool pprof -top -diff_base=profiles/cpu_split.out profiles/cpu_chi.out
```

**output**

```shell
File: handlers.test
Type: cpu
Time: Mar 31, 2024 at 7:24am (MSK)
Duration: 2.81s, Total samples = 1.49s (52.97%)
Showing nodes accounting for -0.35s, 23.49% of 1.49s total
Dropped 1 node (cum <= 0.01s)
      flat  flat%   sum%        cum   cum%
    -0.08s  5.37%  5.37%     -0.07s  4.70%  runtime.pcvalue
     0.04s  2.68%  2.68%      0.27s 18.12%  github.com/andreevym/metric-collector/internal/handlers.BuildMetricByChiParam
    -0.04s  2.68%  5.37%     -0.50s 33.56%  github.com/andreevym/metric-collector/internal/handlers.BuildMetricBySplitParam
     0.04s  2.68%  2.68%      0.04s  2.68%  github.com/go-chi/chi/v5.(*Context).URLParam (inline)
    -0.04s  2.68%  5.37%     -0.11s  7.38%  github.com/stretchr/testify/require.NoError
    -0.04s  2.68%  8.05%     -0.10s  6.71%  runtime.heapBitsSetType
     0.04s  2.68%  5.37%     -0.14s  9.40%  runtime.mallocgc
     0.03s  2.01%  3.36%      0.04s  2.68%  context.(*valueCtx).Value
    -0.03s  2.01%  5.37%     -0.03s  2.01%  indexbytebody
    -0.03s  2.01%  7.38%     -0.05s  3.36%  runtime.gentraceback
    -0.03s  2.01%  9.40%     -0.03s  2.01%  runtime.mapaccess2_fast64
    -0.03s  2.01% 11.41%     -0.04s  2.68%  runtime.newobject (partial-inline)
     0.03s  2.01%  9.40%      0.03s  2.01%  runtime.readvarint (inline)
    -0.03s  2.01% 11.41%     -0.04s  2.68%  runtime.writeHeapBits.flush
    -0.03s  2.01% 13.42%     -0.03s  2.01%  strconv.ParseUint
    -0.03s  2.01% 15.44%     -0.25s 16.78%  strings.genSplit
     0.02s  1.34% 14.09%      0.01s  0.67%  runtime.assertI2I2
    -0.02s  1.34% 15.44%     -0.08s  5.37%  runtime.callers
    -0.02s  1.34% 16.78%     -0.02s  1.34%  runtime.deductAssistCredit
    -0.02s  1.34% 18.12%     -0.02s  1.34%  runtime.duffzero
    -0.02s  1.34% 19.46%     -0.02s  1.34%  runtime.fastrand (inline)
     0.02s  1.34% 18.12%      0.04s  2.68%  runtime.findfunc
     0.02s  1.34% 16.78%     -0.04s  2.68%  runtime.funcspdelta
    -0.02s  1.34% 18.12%     -0.02s  1.34%  runtime.memclrNoHeapPointers
    -0.02s  1.34% 19.46%     -0.02s  1.34%  runtime.readUintptr (inline)
    -0.02s  1.34% 20.81%     -0.07s  4.70%  strings.Index
     0.02s  1.34% 19.46%      0.02s  1.34%  sync.(*Mutex).Unlock (inline)
    -0.01s  0.67% 20.13%     -0.01s  0.67%  countbody
     0.01s  0.67% 19.46%      0.01s  0.67%  gosave_systemstack_switch
    -0.01s  0.67% 20.13%     -0.01s  0.67%  internal/bytealg.IndexByteString
    -0.01s  0.67% 20.81%     -0.01s  0.67%  runtime.(*gcBitsArena).tryAlloc (inline)
    -0.01s  0.67% 21.48%     -0.01s  0.67%  runtime.(*gcControllerState).heapGoalInternal
    -0.01s  0.67% 22.15%     -0.01s  0.67%  runtime.(*m).becomeSpinning (inline)
    -0.01s  0.67% 22.82%     -0.01s  0.67%  runtime.(*mheap).alloc.func1
     0.01s  0.67% 22.15%      0.01s  0.67%  runtime.(*moduledata).textAddr
     0.01s  0.67% 21.48%      0.01s  0.67%  runtime.(*moduledata).textOff (inline)
    -0.01s  0.67% 22.15%     -0.01s  0.67%  runtime.acquirem (inline)
     0.01s  0.67% 21.48%      0.01s  0.67%  runtime.add (inline)
    -0.01s  0.67% 22.15%     -0.01s  0.67%  runtime.arenaIndex (inline)
    -0.01s  0.67% 22.82%     -0.06s  4.03%  runtime.callers.func1
     0.01s  0.67% 22.15%      0.01s  0.67%  runtime.convI2I
     0.01s  0.67% 21.48%      0.01s  0.67%  runtime.efaceeq
     0.01s  0.67% 20.81%      0.01s  0.67%  runtime.elideWrapperCalling (inline)
     0.01s  0.67% 20.13%      0.01s  0.67%  runtime.findmoduledatap (inline)
     0.01s  0.67% 19.46%      0.02s  1.34%  runtime.funcInfo.entry (inline)
    -0.01s  0.67% 20.13%     -0.01s  0.67%  runtime.funcInfo.valid (inline)
    -0.01s  0.67% 20.81%     -0.01s  0.67%  runtime.futex
     0.01s  0.67% 20.13%      0.01s  0.67%  runtime.getMCache (inline)
    -0.01s  0.67% 20.81%     -0.01s  0.67%  runtime.getitab
    -0.01s  0.67% 21.48%     -0.01s  0.67%  runtime.heapBits.next
    -0.01s  0.67% 22.15%     -0.01s  0.67%  runtime.heapBits.nextFast (inline)
     0.01s  0.67% 21.48%      0.01s  0.67%  runtime.lock2
     0.01s  0.67% 20.81%      0.01s  0.67%  runtime.madvise
    -0.01s  0.67% 21.48%     -0.01s  0.67%  runtime.nextFreeFast (inline)
     0.01s  0.67% 20.81%      0.01s  0.67%  runtime.pcvalueCacheKey (inline)
    -0.01s  0.67% 21.48%     -0.01s  0.67%  runtime.releasem (inline)
    -0.01s  0.67% 22.15%     -0.03s  2.01%  runtime.scanobject
    -0.01s  0.67% 22.82%     -0.02s  1.34%  strings.Count
    -0.01s  0.67% 23.49%     -0.05s  3.36%  strings.IndexByte (inline)
    -0.01s  0.67% 24.16%     -0.02s  1.34%  sync.(*RWMutex).Lock
     0.01s  0.67% 23.49%      0.01s  0.67%  sync/atomic.(*Int32).Add (inline)
         0     0% 23.49%     -0.01s  0.67%  github.com/andreevym/metric-collector/internal/handlers.NewRouter
         0     0% 23.49%      1.08s 72.48%  github.com/andreevym/metric-collector/internal/handlers_test.BenchmarkBuildMetricByChiParam
         0     0% 23.49%     -1.42s 95.30%  github.com/andreevym/metric-collector/internal/handlers_test.BenchmarkBuildMetricByParam
         0     0% 23.49%     -0.01s  0.67%  github.com/andreevym/metric-collector/internal/handlers_test.TestHandler_GaugeEndToEnd
         0     0% 23.49%     -0.01s  0.67%  github.com/go-chi/chi/v5.(*Mux).Post (inline)
         0     0% 23.49%     -0.01s  0.67%  github.com/go-chi/chi/v5.(*Mux).handle
         0     0% 23.49%     -0.01s  0.67%  github.com/go-chi/chi/v5.(*node).InsertRoute
         0     0% 23.49%     -0.01s  0.67%  github.com/go-chi/chi/v5.(*node).setEndpoint
         0     0% 23.49%      0.04s  2.68%  github.com/go-chi/chi/v5.RouteContext (inline)
         0     0% 23.49%      0.08s  5.37%  github.com/go-chi/chi/v5.URLParam
         0     0% 23.49%     -0.01s  0.67%  github.com/go-chi/chi/v5.endpoints.Value (inline)
         0     0% 23.49%     -0.01s  0.67%  runtime.(*gcControllerState).trigger
         0     0% 23.49%     -0.01s  0.67%  runtime.(*mcache).nextFree
         0     0% 23.49%     -0.01s  0.67%  runtime.(*mcache).refill
         0     0% 23.49%     -0.01s  0.67%  runtime.(*mcentral).cacheSpan
         0     0% 23.49%     -0.01s  0.67%  runtime.(*mcentral).grow
         0     0% 23.49%     -0.01s  0.67%  runtime.(*mheap).alloc
         0     0% 23.49%      0.01s  0.67%  runtime.(*mheap).freeSpan
         0     0% 23.49%      0.01s  0.67%  runtime.(*mheap).freeSpan.func1
         0     0% 23.49%      0.01s  0.67%  runtime.(*pageAlloc).scavenge
         0     0% 23.49%      0.01s  0.67%  runtime.(*pageAlloc).scavenge.func1
         0     0% 23.49%      0.01s  0.67%  runtime.(*pageAlloc).scavengeOne
         0     0% 23.49%      0.01s  0.67%  runtime.(*scavengerState).init.func2
         0     0% 23.49%      0.01s  0.67%  runtime.(*scavengerState).run
         0     0% 23.49%     -0.08s  5.37%  runtime.Callers (inline)
         0     0% 23.49%      0.01s  0.67%  runtime.bgscavenge
         0     0% 23.49%     -0.02s  1.34%  runtime.fastrandn (inline)
         0     0% 23.49%     -0.02s  1.34%  runtime.findRunnable
         0     0% 23.49%     -0.01s  0.67%  runtime.futexwakeup
         0     0% 23.49%     -0.03s  2.01%  runtime.gcBgMarkWorker
         0     0% 23.49%     -0.03s  2.01%  runtime.gcBgMarkWorker.func2
         0     0% 23.49%     -0.03s  2.01%  runtime.gcDrain
         0     0% 23.49%     -0.01s  0.67%  runtime.gcTrigger.test
         0     0% 23.49%     -0.01s  0.67%  runtime.gopreempt_m
         0     0% 23.49%     -0.01s  0.67%  runtime.goschedImpl
         0     0% 23.49%      0.01s  0.67%  runtime.lock (inline)
         0     0% 23.49%      0.01s  0.67%  runtime.lockWithRank (inline)
         0     0% 23.49%     -0.13s  8.72%  runtime.makeslice
         0     0% 23.49%     -0.01s  0.67%  runtime.mapassign_fast64
         0     0% 23.49%     -0.01s  0.67%  runtime.mcall
         0     0% 23.49%     -0.01s  0.67%  runtime.morestack
         0     0% 23.49%     -0.01s  0.67%  runtime.newMarkBits
         0     0% 23.49%     -0.01s  0.67%  runtime.newstack
         0     0% 23.49%     -0.01s  0.67%  runtime.notewakeup
         0     0% 23.49%     -0.01s  0.67%  runtime.park_m
         0     0% 23.49%      0.01s  0.67%  runtime.pcdatastart (inline)
         0     0% 23.49%     -0.01s  0.67%  runtime.runSafePointFn
         0     0% 23.49%     -0.02s  1.34%  runtime.schedule
         0     0% 23.49%      0.01s  0.67%  runtime.sysUnused
         0     0% 23.49%      0.01s  0.67%  runtime.sysUnusedOS
         0     0% 23.49%     -0.07s  4.70%  runtime.systemstack
         0     0% 23.49%     -0.03s  2.01%  strconv.ParseInt
         0     0% 23.49%     -0.25s 16.78%  strings.Split (inline)
         0     0% 23.49%      0.04s  2.68%  sync.(*RWMutex).Unlock
         0     0% 23.49%     -0.34s 22.82%  testing.(*B).launch
         0     0% 23.49%     -0.34s 22.82%  testing.(*B).runN
         0     0% 23.49%     -0.09s  6.04%  testing.(*common).Helper
         0     0% 23.49%     -0.01s  0.67%  testing.tRunner
```

### diff mem

**command**

```shell
go tool pprof -top -diff_base=profiles/mem_split.out profiles/mem_chi.out
```

**output**

```shell
File: handlers.test
Type: alloc_space
Time: Mar 31, 2024 at 7:24am (MSK)
Showing nodes accounting for -82.86MB, 55.41% of 149.54MB total
Dropped 9 nodes (cum <= 0.75MB)
      flat  flat%   sum%        cum   cum%
  -74.51MB 49.82% 49.82%   -74.51MB 49.82%  strings.genSplit
     -60MB 40.12% 89.95%  -134.51MB 89.95%  github.com/andreevym/metric-collector/internal/handlers.BuildMetricBySplitParam
      55MB 36.78% 53.17%       55MB 36.78%  github.com/andreevym/metric-collector/internal/handlers.BuildMetricByChiParam
   -4.41MB  2.95% 56.11%    -2.67MB  1.79%  compress/flate.NewWriter
   -1.16MB  0.77% 56.89%    -1.16MB  0.77%  runtime/pprof.StartCPUProfile
    1.10MB  0.74% 56.15%     1.10MB  0.74%  compress/flate.newDeflateFast (inline)
    0.64MB  0.42% 55.73%     1.74MB  1.16%  compress/flate.(*compressor).init
   -0.52MB  0.34% 56.07%    -0.52MB  0.34%  compress/flate.(*dictDecoder).init (inline)
   -0.52MB  0.34% 56.41%    -0.52MB  0.34%  regexp.(*bitState).reset
    0.50MB  0.34% 56.08%     0.50MB  0.34%  bufio.NewWriterSize (inline)
    0.50MB  0.33% 55.75%     0.50MB  0.33%  net/http.(*Transport).queueForIdleConn
   -0.50MB  0.33% 56.08%    -1.02MB  0.68%  io.ReadAll
    0.50MB  0.33% 55.75%     0.50MB  0.33%  net/http.NewRequestWithContext
    0.50MB  0.33% 55.41%     0.50MB  0.33%  github.com/go-chi/chi/v5.(*node).InsertRoute
   -0.50MB  0.33% 55.75%    -0.50MB  0.33%  net/http.ReadResponse
    0.50MB  0.33% 55.41%     0.50MB  0.33%  regexp/syntax.parse
         0     0% 55.41%    -0.52MB  0.34%  compress/flate.NewReader
         0     0% 55.41%    -0.52MB  0.34%  compress/gzip.(*Reader).Reset
         0     0% 55.41%    -0.52MB  0.34%  compress/gzip.(*Reader).readHeader
         0     0% 55.41%    -2.67MB  1.79%  compress/gzip.(*Writer).Write
         0     0% 55.41%    -0.52MB  0.34%  compress/gzip.NewReader
         0     0% 55.41%    -0.88MB  0.59%  github.com/andreevym/metric-collector/internal/compressor.Compress
         0     0% 55.41%     0.50MB  0.33%  github.com/andreevym/metric-collector/internal/handlers.NewRouter
         0     0% 55.41%    -0.81MB  0.54%  github.com/andreevym/metric-collector/internal/handlers.ServiceHandlers.PostUpdateHandler
         0     0% 55.41%    -0.88MB  0.59%  github.com/andreevym/metric-collector/internal/handlers.TestGzipCompressionUpdate.func3
         0     0% 55.41%     0.50MB  0.33%  github.com/andreevym/metric-collector/internal/handlers.TestGzipCompressionValue
         0     0% 55.41%    -0.54MB  0.36%  github.com/andreevym/metric-collector/internal/handlers.TestGzipCompressionValue.func3
         0     0% 55.41%    -0.50MB  0.33%  github.com/andreevym/metric-collector/internal/handlers.buildMetricByBody
         0     0% 55.41%    -0.50MB  0.33%  github.com/andreevym/metric-collector/internal/handlers.buildMetricByRequest
         0     0% 55.41%    -0.52MB  0.34%  github.com/andreevym/metric-collector/internal/handlers.testCompressRequest
         0     0% 55.41%       55MB 36.78%  github.com/andreevym/metric-collector/internal/handlers_test.BenchmarkBuildMetricByChiParam
         0     0% 55.41%  -134.51MB 89.95%  github.com/andreevym/metric-collector/internal/handlers_test.BenchmarkBuildMetricByParam
         0     0% 55.41%     0.50MB  0.33%  github.com/andreevym/metric-collector/internal/handlers_test.TestHandler_GaugeEndToEnd
         0     0% 55.41%     0.50MB  0.33%  github.com/andreevym/metric-collector/internal/handlers_test.TestUpdateHandler.func1
         0     0% 55.41%        1MB  0.67%  github.com/andreevym/metric-collector/internal/handlers_test.signedTestRequest
         0     0% 55.41%        1MB  0.67%  github.com/andreevym/metric-collector/internal/handlers_test.testRequest (inline)
         0     0% 55.41%    -0.74MB  0.49%  github.com/andreevym/metric-collector/internal/middleware.(*Middleware).ResponseGzipMiddleware.func1
         0     0% 55.41%     0.50MB  0.33%  github.com/go-chi/chi/v5.(*Mux).Post (inline)
         0     0% 55.41%    -0.74MB  0.49%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 55.41%     0.50MB  0.33%  github.com/go-chi/chi/v5.(*Mux).handle
         0     0% 55.41%    -0.74MB  0.49%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
         0     0% 55.41%     0.50MB  0.33%  gopkg.in/yaml%2ev3.init
         0     0% 55.41%    -1.67MB  1.12%  main.main
         0     0% 55.41%     0.50MB  0.33%  net/http.(*Client).Do (inline)
         0     0% 55.41%     0.50MB  0.33%  net/http.(*Client).do
         0     0% 55.41%     0.50MB  0.33%  net/http.(*Client).send
         0     0% 55.41%     0.50MB  0.33%  net/http.(*Transport).RoundTrip
         0     0% 55.41%     0.50MB  0.33%  net/http.(*Transport).getConn
         0     0% 55.41%     0.50MB  0.33%  net/http.(*Transport).roundTrip
         0     0% 55.41%     0.50MB  0.34%  net/http.(*conn).readRequest
         0     0% 55.41%    -0.52MB  0.34%  net/http.(*gzipReader).Read
         0     0% 55.41%    -0.50MB  0.33%  net/http.(*persistConn).readLoop
         0     0% 55.41%    -0.50MB  0.33%  net/http.(*persistConn).readResponse
         0     0% 55.41%    -0.74MB  0.49%  net/http.HandlerFunc.ServeHTTP
         0     0% 55.41%     0.50MB  0.33%  net/http.NewRequest (inline)
         0     0% 55.41%     0.50MB  0.34%  net/http.newBufioWriterSize
         0     0% 55.41%     0.50MB  0.33%  net/http.send
         0     0% 55.41%    -0.74MB  0.49%  net/http.serverHandler.ServeHTTP
         0     0% 55.41%    -0.52MB  0.34%  regexp.(*Regexp).MatchString (inline)
         0     0% 55.41%    -0.52MB  0.34%  regexp.(*Regexp).backtrack
         0     0% 55.41%    -0.52MB  0.34%  regexp.(*Regexp).doExecute
         0     0% 55.41%    -0.52MB  0.34%  regexp.(*Regexp).doMatch (inline)
         0     0% 55.41%     0.50MB  0.33%  regexp.Compile (inline)
         0     0% 55.41%     0.50MB  0.33%  regexp.MustCompile
         0     0% 55.41%     0.50MB  0.33%  regexp.compile
         0     0% 55.41%     0.50MB  0.33%  regexp/syntax.Parse (inline)
         0     0% 55.41%     0.50MB  0.33%  runtime.doInit
         0     0% 55.41%    -1.17MB  0.78%  runtime.main
         0     0% 55.41%    -2.05MB  1.37%  runtime/pprof.(*profileBuilder).build
         0     0% 55.41%    -1.55MB  1.04%  runtime/pprof.(*profileBuilder).flush
         0     0% 55.41%    -2.09MB  1.40%  runtime/pprof.(*profileBuilder).pbSample
         0     0% 55.41%    -2.05MB  1.37%  runtime/pprof.profileWriter
         0     0% 55.41%   -74.51MB 49.82%  strings.Split (inline)
         0     0% 55.41%   -79.51MB 53.17%  testing.(*B).launch
         0     0% 55.41%   -79.51MB 53.17%  testing.(*B).runN
         0     0% 55.41%    -1.67MB  1.12%  testing.(*M).Run
         0     0% 55.41%    -1.16MB  0.77%  testing.(*M).before
         0     0% 55.41%    -0.52MB  0.34%  testing.runExamples
         0     0% 55.41%    -0.52MB  0.34%  testing/internal/testdeps.TestDeps.MatchString
         0     0% 55.41%    -1.16MB  0.77%  testing/internal/testdeps.TestDeps.StartCPUProfile
```

## Profile tool

Example how to work from code and for recive cpu mem profile from golang code

Diff between source version 'base' and version after changes for split param uri without chi 'result'

In result parse arg with chi work faster in this block of code, because chi parsed params earlier and put them through
context.

**CPU**

**command**

```shell
go tool pprof -top -diff_base=profiles/base_cpu.pprof profiles/result_cpu.pprof
```

**result**

```
File: ___go_build_github_com_andreevym_metric_collector_profiles
Type: cpu
Time: Mar 31, 2024 at 12:09am (MSK)
Duration: 401.27ms, Total samples = 0 
Showing nodes accounting for 0, 0% of 0 total
      flat  flat%   sum%        cum   cum%
```

**MEM**

**command**

```shell
go tool pprof -top -diff_base=profiles/base_mem.pprof profiles/result_mem.pprof
```

**result**

```
File: ___go_build_github_com_andreevym_metric_collector_profiles
Type: inuse_space
Time: Mar 31, 2024 at 12:09am (MSK)
Showing nodes accounting for 0, 0% of 1.72MB total
      flat  flat%   sum%        cum   cum%

```