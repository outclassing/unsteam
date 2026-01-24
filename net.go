package main

import (
	"os"
	"io"
	"net/http"
	"encoding/json"
)

func fetchJson[T any](url string) (T, error) {
	var v T
	resp, err := http.Get(url)
	if err != nil {
		return v, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&v)
	return v, err
}

func downloadToFile(url, path string) error {
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

func fetchDepotKey(depotId string, token string, path string) (Key, error) {
	req, _ := http.NewRequest("GET",
		"https://unsteam.cloudflare-delivery914.workers.dev/key?id=" + depotId, nil)
	req.Header.Set("Authorization", "Bearer " + token)
	
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	
	var key Key
	if err = json.NewDecoder(resp.Body).Decode(&key); err != nil {
		return Key{}, err
	}
	
	if err = os.WriteFile(path + depotId + ".txt", []byte(depotId+";"+key.Value), 0644); err != nil {
		return key, err
	}
	
	return key, nil
}
