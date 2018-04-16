package testengine

import (
	"context"
	"gat/conf"
	"gat/invoker"
	"gat/metric"
	"log"
	"reflect"
	"strings"
	"time"
)

type TestEngine interface {
	Test(v interface{})
	Bench(v interface{}, conf *conf.Benchmark)
	Invoke(data string) *metric.Result
}

func NewTestEngine(conf *conf.Config, benchmark bool) TestEngine {
	return &testEngine{
		invoker:   invoker.NewGatewayInvoker(conf.GatewayAddr),
		benchmark: benchmark,
	}
}

type testEngine struct {
	invoker   invoker.Invoker
	benchmark bool
}

func (t *testEngine) Test(v interface{}) {
	tfs, _ := allTestFuncs(v)
	for name, tf := range tfs {
		log.Printf("TEST:\t[%s] <<<<<<\n", name)
		var succeed = "FAILED\n"
		if tf() {
			succeed = "SUCCEED >>>>>>\n"
		}
		log.Println(succeed)
	}

}

func (t *testEngine) Bench(v interface{}, conf *conf.Benchmark) {
	if !t.benchmark {
		return
	}
	_, bfs := allTestFuncs(v)
	log.Println("all bench funcs", len(bfs))
	for name, bf := range bfs {
		t.runBenchmark(name, bf, conf)
	}
}

func (t *testEngine) runBenchmark(name string, bf BenchFunc, conf *conf.Benchmark) {
	log.Printf("BENCHMARK:\t[%s] <<<<<<\n", name)
	// TODO: run benchmark
	tokenChan := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Duration(conf.Duration)*time.Second, cancel)

	metric := metric.NewMetric(name, time.Second)
	go tokenGen(cancel, conf, tokenChan)

	for range tokenChan {
		go attack(ctx, metric, bf)
	}

}

func tokenGen(cancel context.CancelFunc, conf *conf.Benchmark, tokenChan chan<- struct{}) {

	defer close(tokenChan)

	if len(conf.Goroutine) == 0 {
		return
	}
	var prevtokens = conf.Goroutine[0]
	var nextIndex = 1
	for i := 0; i < prevtokens; i++ {
		tokenChan <- struct{}{}
	}
	ticker := time.NewTicker(time.Second * time.Duration(conf.Duration/len(conf.Goroutine)))
	for range ticker.C {
		if len(conf.Goroutine) == nextIndex {
			return
		}
		tokens := conf.Goroutine[nextIndex] - prevtokens
		prevtokens = conf.Goroutine[nextIndex]
		nextIndex++
		for i := 0; i < tokens; i++ {
			tokenChan <- struct{}{}
		}
	}
}

func attack(ctx context.Context, metric *metric.Metric, bf BenchFunc) {
	log.Println("attack")
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		metric.OnSuccess(bf())
	}
}

func (t *testEngine) Invoke(data string) *metric.Result {
	began := time.Now()
	res, err := t.invoker.Invoke(context.Background(), data)
	return &metric.Result{
		Err:      err,
		Request:  data,
		Response: res,
		Latency:  time.Now().Sub(began),
	}
}

// TestFunc to run test
type TestFunc func() bool
type BenchFunc func() (bool, time.Duration)

func allTestFuncs(v interface{}) (map[string]TestFunc, map[string]BenchFunc) {
	var tfs = map[string]TestFunc{}
	var bfs = map[string]BenchFunc{}

	val := reflect.ValueOf(v)
	typ := val.Type()

	for i := 0; i < val.NumMethod(); i++ {
		mth := val.Method(i)
		tym := typ.Method(i)
		mty := mth.Type()
		mname := tym.Name
		if tym.PkgPath != "" {
			continue
		}

		if mty.NumIn() != 0 {
			continue
		}

		if strings.HasPrefix(mname, "Test") {
			if mty.NumOut() != 1 {
				continue
			}

			if firstRet := mty.Out(0); firstRet.Kind() != reflect.Bool {
				continue
			}
			tfs[mname] = mth.Interface().(func() bool)
		} else if strings.HasPrefix(mname, "Benchmark") {
			if mty.NumOut() != 2 {
				continue
			}

			if firstRet := mty.Out(0); firstRet.Kind() != reflect.Bool {
				continue
			}
			if secondRet := mty.Out(1); secondRet.Kind() != reflect.TypeOf(time.Second).Kind() {
				continue
			}

			bfs[mname] = mth.Interface().(func() (bool, time.Duration))
		}
	}
	return tfs, bfs
}
