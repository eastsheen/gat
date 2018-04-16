package util

import (
	"encoding/xml"
	"io/ioutil"

	"github.com/json-iterator/go"
)

// FromJSONFile decode json from file
func FromJSONFile(path string, i interface{}) error {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(bts, i)
}

// FromJSONString decode json string
func FromJSONString(data string, i interface{}) error {
	return jsoniter.UnmarshalFromString(data, i)
}

// FromJSONBytes decode json bytes
func FromJSONBytes(data []byte, i interface{}) error {
	return jsoniter.Unmarshal(data, i)
}

// FromXMLFile decode xml from file
func FromXMLFile(path string, i interface{}) error {
	bts, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return xml.Unmarshal(bts, i)
}

// FromXMLString decode xml from string
func FromXMLString(data string, i interface{}) error {
	return xml.Unmarshal([]byte(data), i)
}

// FromXMLBytes decode xml from bytes
func FromXMLBytes(data []byte, i interface{}) error {
	return xml.Unmarshal(data, i)
}
