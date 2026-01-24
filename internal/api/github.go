package api

type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
	Digest      string `json:"digest"`
}

type Release struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}
