package pkg

import (
	"encoding/json"
	"io"
	"net/http"
)

type Header struct {
	Key   string
	Value func() string
}

func RequestBytes(url string, header *Header) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if header != nil {
		req.Header.Set(header.Key, header.Value())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func RequestJSON[T any](url string, header *Header) (T, error) {
	var out T
	data, err := RequestBytes(url, header)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}
