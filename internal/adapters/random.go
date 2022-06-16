package adapters

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

type randomOrg struct {
	*sync.Mutex
	client *http.Client
}

const (
	baseURI = "https://www.random.org/integers/?"
)

func NewRandomOrg() (NumbersGenerator, error) {
	return &randomOrg{
		Mutex:  &sync.Mutex{},
		client: &http.Client{},
	}, nil
}

func (r randomOrg) Get(start, stop, count int64) ([]int64, error) {
	r.Lock()
	defer r.Unlock()

	var params = map[string][]string{
		"num":    {strconv.Itoa(int(count))},
		"min":    {strconv.Itoa(int(start))},
		"max":    {strconv.Itoa(int(stop))},
		"col":    {"1"},
		"base":   {"10"},
		"format": {"plain"},
		"rnd":    {"new"},
	}

	p := url.Values(params)

	requestURL := fmt.Sprintf("%s%s", baseURI, p.Encode())

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)

	if err != nil {
		return nil, err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	output := make([]int64, 0, count)

	s := bufio.NewScanner(res.Body)
	for s.Scan() {
		n, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			return nil, err
		}

		output = append(output, n)
	}

	return output, nil
}
