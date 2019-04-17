package script

import (
	"errors"
	"flag"
	"net/http"
	"testing"
)

var (
	totalRequest = flag.Int("total", 100000, "total request")
	clientNum    = flag.Int("n", 200, "the number of concurrent request")
)

var client *http.Client

func init() {
	client = new(http.Client)
	flag.Parse()
}

func testFunc() error {
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("err")
	}

	return nil
}

func TestBenchScript(t *testing.T) {
	BenchScript(*totalRequest, *clientNum, testFunc)
}
