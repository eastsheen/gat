package conf

import (
	"errors"
	"gat/util"
)

var (
	ErrConfNotExists = errors.New("config not exists")
)

type ConfigCenter interface {
	GetConfig(env string) (*Config, error)
}

type configCenter map[string]*Config

func (c configCenter) GetConfig(env string) (*Config, error) {
	if ret, ok := c[env]; ok {
		return ret, nil
	}
	return nil, ErrConfNotExists
}

type Config struct {
	// http ip 地址, 从网关调用
	GatewayAddr string `json:"gateway,omitempty"`
	// zookeeper 注册中心地址, soa调用
	ZKAddr string `json:"zkaddr,omitempty"`
}

type Benchmark struct {
	Goroutine []int `json:"goroutine"`
	// Second
	Duration int `json:"duration"`
	// PrintInterval of print bench info
	PrintInterval int `json:"interval"`
}

func InitWithJSONFile(path string) (ConfigCenter, error) {
	var cc = configCenter{}
	if err := util.FromJSONFile(path, &cc); err != nil {
		return nil, err
	}
	return cc, nil
}

func InitWithJSONString(data string) (ConfigCenter, error) {
	var cc = configCenter{}
	if err := util.FromJSONString(data, &cc); err != nil {
		return nil, err
	}
	return cc, nil
}

func InitWithJSONBytes(data []byte) (ConfigCenter, error) {
	var cc = configCenter{}
	if err := util.FromJSONBytes(data, &cc); err != nil {
		return nil, err
	}
	return cc, nil
}
