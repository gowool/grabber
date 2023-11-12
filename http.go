package grabber

import (
	"compress/gzip"
	"fmt"
	"net/http"
)

var (
	Client         = &http.Client{}
	DefaultHeaders = map[string]string{
		"Accept":          "text/html",
		"Accept-Encoding": "gzip",
		"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/111.0",
	}
)

func NewRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range DefaultHeaders {
		req.Header.Set(k, v)
	}

	return req, nil
}

func Do(req *http.Request) (*Page, error) {
	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: response status code: %d", res.StatusCode)
	}

	body := res.Body

	if res.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(body)
		if err != nil {
			return nil, err
		}
		defer body.Close()
	}

	page := NewPage(req.URL.String())

	if err = page.Parse(body); err != nil {
		return nil, err
	}

	return page, page.ToAbs()
}
