package lib

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type Fetcher interface {
	Get(url string) (resp *http.Response, err error)
}

type defaultFetcher struct {
}

func NewDefaultFetcher() Fetcher {
	return defaultFetcher{}
}

func (f defaultFetcher) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

type dummyFetcherCount struct {
	index     int
	responses []string
}

// NewDummyFetcherCount creates a fetcher which returns response based on number of calls.
func NewDummyFetcherCount(responses []string) Fetcher {
	return &dummyFetcherCount{
		responses: responses,
	}
}

func (f *dummyFetcherCount) Get(url string) (*http.Response, error) {
	if len(f.responses) <= f.index {
		return nil, errors.New("no more fetch")
	}
	reader := ioutil.NopCloser(strings.NewReader(f.responses[f.index]))
	f.index++
	return &http.Response{
		StatusCode: 200,
		Body:       reader,
	}, nil
}
