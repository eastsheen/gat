package testengine

import (
	"gat/conf"
	"testing"
)

func TestAllTestFuncs(t *testing.T) {
	ts := new(TestStruct)

	tfs, bfs := allTestFuncs(ts)

	if len(tfs) != 1 {
		t.Fail()
	}

	if _, ok := tfs["TestSth"]; !ok {
		t.Fail()
	}
}

var testst bool
var benchst bool

func TestRunTests(t *testing.T) {
	ts := new(TestStruct)
	engine := NewTestEngine(&conf.Config{})
	engine.Test(ts)
	if !testst {
		t.Fail()
	}
}

type TestStruct struct {
}

func (t *TestStruct) TestSth() bool {
	testst = true
	return true
}

func (t *TestStruct) BenchmarkSth() bool {
	benchst = true
	return true
}
