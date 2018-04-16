package conf

import (
	"io/ioutil"
	"testing"

	"github.com/json-iterator/go"
)

func TestConfigCenterMarshal(t *testing.T) {
	var cc = configCenter{}

	cc["dev"] = &Config{
		GatewayAddr: "http://47.97.102.217:8080",
	}

	cc["uat"] = &Config{
		GatewayAddr: "http://47.97.102.217:8080",
	}

	bts, _ := jsoniter.MarshalIndent(&cc, "", " ")
	t.Log(string(bts))
}

func TestConfigCenterUnmarshal(t *testing.T) {
	var cc = configCenter{}
	bts, err := ioutil.ReadFile("../testdata/testconfig.json")
	if err != nil {
		t.Error(err)
	}

	if err := jsoniter.Unmarshal(bts, &cc); err != nil {
		t.Error(err)
	}
}
