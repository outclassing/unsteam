package api

import "fmt"

type App struct {
	Depots []Depot `json:"depots"`
}

type Depot struct {
	Id         int    `json:"depotid"`
	ManifestId string `json:"manifestid"`
	Size       int64  `json:"size_bytes"`
	BuildId    int    `json:"buildid"`
	Time       int64  `json:"timeupdated"`
}

type Manifest struct {
	Content string `json:"content"`
}

type Key struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}

func SteamDepotUrl(appId string) string {
	return "https://manifest.steam.run/api/depot/" + appId
}

func SteamManifestUrl(id string) string {
	return "https://manifest.steam.run/api/manifest/" + id
}

func SteamCdnManifestUrl(depotId string, manifestId string, content string) string {
	return fmt.Sprintf("https://steampipe.akamaized.net/depot/%s/manifest/%s/5/%s", depotId, manifestId, content)
}
