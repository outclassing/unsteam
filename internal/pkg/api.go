package pkg

type Manifest struct {
	Id          string `json:"manifestid"`
	DepotId     int    `json:"depotid"`
	Size        int64  `json:"size_bytes"`
	BuildId     int    `json:"buildid"`
	Time        int64  `json:"timeupdated"`
	RequestCode string
}

type App struct {
	Manifests []Manifest `json:"depots"`
}

type RequestCode struct {
	Content string `json:"content"`
}

type DepotKey struct {
	Value  string `json:"value"`
	Source string `json:"source"`
}
