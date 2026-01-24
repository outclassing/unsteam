package pkg

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"unsteam/internal/api"
)

func FetchJson[T any](url string) (T, error) {
	var v T
	resp, err := http.Get(url)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}

func DownloadToFile(url, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func FetchDepotKey(depotId string, token string, path string) (api.Key, error) {
	req, _ := http.NewRequest("GET",
		"https://unsteam.cloudflare-delivery914.workers.dev/key?id="+depotId, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return api.Key{}, err
	}
	defer resp.Body.Close()

	var key api.Key
	if err = json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return api.Key{}, err
	}

	if err = os.WriteFile(path+depotId+".txt", []byte(depotId+";"+key.Value), 0644); err != nil {
		return key, err
	}

	return key, nil
}
