package app

import (
	"context"
	"fmt"
	"path/filepath"
	"time"
	"unsteam/internal/api"
	"unsteam/internal/pkg"
)

func DownloadDepot(ctx context.Context) {
	appId := PromptNumeric(ctx, "Enter app id", "70")
	depotId := PromptNumeric(ctx, "Enter depot id", "71")

	stream := NewStream("Downloading manifest")
	manifestPath := "data/manifests/"
	pkg.EnsureDir(manifestPath)
	app, ok := fetchApp(appId)
	if !ok {
		stream.Error("Failed to fetch app info")
		return
	}

	depot, ok := findDepot(app, depotId)
	if !ok {
		stream.Error("Failed to find depot")
		return
	}

	stream.WriteLine(fmt.Sprintf("Manifest %s selected", depot.ManifestId))
	stream.WriteLine(fmt.Sprintf("Build ID: %s", fmt.Sprint(depot.BuildId)))
	stream.WriteLine(fmt.Sprintf("Time updated: %s", time.Unix(depot.Time, 0).Format(time.RFC3339)))

	manifest, ok := fetchManifest(depot.ManifestId)
	if !ok {
		stream.Error("Failed to fetch manifest")
		return
	}

	if err := downloadManifest(depotId, depot.ManifestId, manifest.Content, stream, ctx); err != nil {
		stream.Error(fmt.Sprintf("Failed to download manifest file: %s", err))
		return
	}

	stream.Ok("Manifest downloaded successfully")
	stream = NewStream("Fetching depot key")
	keyPath := "data/keys/"
	pkg.EnsureDir(keyPath)

	apiKey, err := pkg.Env("US_API_KEY")
	if err != nil {
		stream.Error("Failed to get API key, please set US_API_KEY")
		return
	}

	key, err := pkg.FetchDepotKey(depotId, apiKey, keyPath)
	if err != nil {
		stream.Error("Failed to fetch depot key: " + err.Error())
	}

	stream.WriteLine("Saved decryption key to temporary file")
	stream.Ok("Depot key retrieved successfully, source: " + key.Source)
	stream = NewStream("Downloading depot")
	stream.WriteLine("Running DepotDownloader, please wait")

	manifestFile := filepath.Join(manifestPath, depot.ManifestId)
	keyFile := filepath.Join(keyPath, depotId+".txt")
	if err = runDepotDownloader(appId, keyFile, depotId, depot.ManifestId, manifestFile); err != nil {
		stream.Error("Failed to run DepotDownloader: " + err.Error())
	}

	stream.WriteLine("Removing temporary keyfile and manifest")
	clean(manifestPath, keyPath)

	depotPath := filepath.Join("depots", depotId)
	localSize, err := pkg.FolderSize(depotPath)
	if err != nil {
		stream.Error("Failed to get local folder size: " + err.Error())
	}
	stream.WriteLine(fmt.Sprintf("File sizes: %s (remote) %s (local)", fmt.Sprint(depot.Size), fmt.Sprint(localSize)))

	stream.Ok("Depot downloaded successfully")

	verifyDepotSize(localSize, depot.Size, depotPath, ctx)
	MainMenu(ctx)
}

func fetchApp(id string) (api.App, bool) {
	url := api.SteamDepotUrl(id)
	app, err := pkg.FetchJson[api.App](url)
	if err != nil {
		return api.App{}, false
	}
	return app, true
}

func fetchManifest(id string) (api.Manifest, bool) {
	url := api.SteamManifestUrl(id)
	manifest, err := pkg.FetchJson[api.Manifest](url)
	if err != nil {
		return api.Manifest{}, false
	}
	return manifest, true
}

func downloadManifest(depotId string, manifestId string, content string, s *Stream, ctx context.Context) error {
	url := api.SteamCdnManifestUrl(depotId, manifestId, content)
	path := filepath.Join("data/manifests", manifestId)

	if err := pkg.DownloadToFile(url, path); err != nil {
		return err
	}

	if err := pkg.ExtractArchive(path, ctx); err != nil {
		return err
	}
	s.WriteLine("Extracted archive to manifests")
	return nil
}

func findDepot(app api.App, id string) (api.Depot, bool) {
	for _, d := range app.Depots {
		if fmt.Sprint(d.Id) == id {
			return d, true
		}
	}
	return api.Depot{}, false
}

func verifyDepotSize(local int64, remote int64, depotPath string, ctx context.Context) {
	if local == remote {
		return
	}

	if !Confirmation(ctx, "Ignore file size divergence?") {
		if err := pkg.RemoveAll(depotPath); err != nil {
			Message("Failed to remove file: " + err.Error())
		}
		Message("Depot discarded; please try download again")
		return
	}
}

func clean(manifestPath string, keyPath string) error {
	if err := pkg.RemoveAll(manifestPath); err != nil {
		return err
	}

	if err := pkg.RemoveAll(keyPath); err != nil {
		return err
	}
	return nil
}
