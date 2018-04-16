package invoker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Invoker interface {
	Invoke(ctx context.Context, data string) (string, error)
}

func NewGatewayInvoker(gatewayAddr string) Invoker {
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	return &httpGatewayInvoker{
		apiurl: fmt.Sprintf("%s/api", gatewayAddr),
		client: &http.Client{},
	}
}

type httpGatewayInvoker struct {
	apiurl string
	client *http.Client
}

func (i *httpGatewayInvoker) Invoke(ctx context.Context, data string) (string, error) {

	resp, err := i.client.Post(i.apiurl, "application/json", strings.NewReader(data))
	if err != nil {
		log.Println("err post:", err)
		return "", err
	}
	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("err read body:", err)
		return "", err
	}
	return string(bts), nil
}
