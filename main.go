package main

import (
	"flag"
	"fmt"
	"gat/conf"
	"gat/testengine"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/json-iterator/go"
)

var config = flag.String("config", "./config.json", "Path to json config file")
var bench = flag.Bool("bench", false, "Run benchmarks")
var env = flag.String("env", "dev", "Env")
var rate = flag.Int("rate", 10, "Specifies the requests per second rate to issue against the targets.")
var output = flag.String("out", "", "Output file (default stdout)")
var cpus = flag.Int("cpus", runtime.NumCPU(), "Number of CPUs to use (default 4)")

var configCenter conf.ConfigCenter

var writer = os.Stdout
var engine testengine.TestEngine

func main() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	if *output != "" {
		f, err := os.OpenFile(*output, os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			log.Println("err open file:", *output, ", err:", err)
			return
		}
		defer f.Close()
		writer = f
		log.SetOutput(writer)
	}

	cc, err := conf.InitWithJSONFile(*config)
	if err != nil {
		log.Println(err)
		return
	}
	configCenter = cc

	config, err := configCenter.GetConfig(*env)
	if err != nil {
		log.Println("err get config env: ", *env, ", err: ", err)
		return
	}

	runtime.GOMAXPROCS(*cpus)

	log.Printf("%#v", config)

	engine = testengine.NewTestEngine(config, *bench)

	ustc := &UserServiceTestCase{
		Version: "4.1.1",
	}

	engine.Test(ustc)

	engine.Bench(ustc, &conf.Benchmark{
		Goroutine:     []int{20, 40, 60, 80},
		Duration:      60,
		PrintInterval: 60,
	})
}

type UserServiceTestCase struct {
	Version string
}

func (u *UserServiceTestCase) TestVersionCheck() bool {
	// 测试逻辑
	result := engine.Invoke(fmt.Sprintf(`{
		"action": "client.version.check.test",
		"systemCode":"62",
		"sourceId":"a0009999",
		"appVersion":"%s"
      }`, u.Version))

	// 打印日志
	result.Report(writer)
	if result.Err != nil {
		return false
	}

	if jsoniter.Get([]byte(result.Response), "code").ToInt() != 0 {
		return false
	}

	return true
}

func (u *UserServiceTestCase) BenchmarkVersionCheck() (bool, time.Duration) {
	result := engine.Invoke(fmt.Sprintf(`{
		"action": "client.version.check.test",
		"systemCode":"62",
		"sourceId":"a0009999",
		"appVersion":"%s"
      }`, u.Version))
	// 打印日志
	// result.Report(writer)
	if result.Err != nil {
		return false, result.Latency
	}

	if jsoniter.Get([]byte(result.Response), "code").ToInt() != 0 {
		return false, result.Latency
	}
	return true, result.Latency
}
