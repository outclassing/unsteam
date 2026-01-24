package main

import "context"
import "fmt"
import "path/filepath"

func downloadDepot(ctx context.Context) {
	appId := promptNumeric(ctx, "Enter app id", "70")
	depotId := promptNumeric(ctx, "Enter depot id", "71")
	
	stream := NewStream("Downloading manifest")
	manifestPath := "data/manifests/"
	ensureDir(manifestPath)
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
	printDepotInfo(stream, depot)
	
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
	ensureDir(keyPath)

	apiKey, err := env("US_API_KEY")
	if err != nil {
		stream.Error("Failed to get API key, please set US_API_KEY")
		return
	}

	key, err := fetchDepotKey(depotId, apiKey, keyPath)
	if err != nil {
		stream.Error("Failed to fetch depot key: " + err.Error())
	}
	
	stream.WriteLine("Saved decryption key to temporary file")
	stream.Ok("Depot key retrieved successfully, source: " + key.Source)
	stream = NewStream("Downloading depot")
	stream.WriteLine("Running DepotDownloader, please wait")
	
	manifestFile := filepath.Join(manifestPath, depot.ManifestId)
	keyFile := filepath.Join(keyPath, depotId + ".txt")
	if err = runDepotDownloader(appId, keyFile, depotId, depot.ManifestId, manifestFile); err != nil {
		stream.Error("Failed to run DepotDownloader: " + err.Error())
	}
	
	stream.WriteLine("Removing temporary keyfile and manifest")
	clean(manifestPath, keyPath)
	
	depotPath := filepath.Join("depots", depotId)
	localSize, err := folderSize(depotPath)
	if err != nil {
		stream.Error("Failed to get local folder size: " + err.Error())
	}
	stream.WriteLine(fmt.Sprintf("File sizes: %s (remote) %s (local)", fmt.Sprint(depot.Size), fmt.Sprint(localSize)))
	
	stream.Ok("Depot downloaded successfully")
	
	verifyDepotSize(localSize, depot.Size, depotPath, ctx)
	mainMenu(ctx)
}

func fetchApp(id string) (App, bool) {
	url := steamDepotUrl(id)
	app, err := fetchJson[App](url)
	if err != nil {
		return App{}, false
	}
	return app, true
}

func fetchManifest(id string) (Manifest, bool) {
	url := steamManifestUrl(id)
	manifest, err := fetchJson[Manifest](url)
	if err != nil {
		return Manifest{}, false
	}
	return manifest, true
}

func downloadManifest(depotId string, manifestId string, content string, s *Stream, ctx context.Context) error {
	url := steamCdnManifestUrl(depotId, manifestId, content)
	path := filepath.Join("data/manifests", manifestId)
	
	if err := downloadToFile(url, path); err != nil {
		return err
	}
	
	if err := extractArchive(path, ctx); err != nil {
		return err
	}
	s.WriteLine("Extracted archive to manifests")
	return nil
}

func findDepot(app App, id string) (Depot, bool) {
	for _, d := range app.Depots {
		if fmt.Sprint(d.Id) == id {
			return d, true
		}
	}
	return Depot{}, false
}

func verifyDepotSize(local int64, remote int64, depotPath string, ctx context.Context) {
	if local == remote {
		return
	}
	
	if !confirmation(ctx, "Ignore file size divergence?") {
		if err := removeAll(depotPath); err != nil {
			message("Failed to remove file: " + err.Error())
		}
		message("Depot discarded; please try download again")
		return
	}
}

func clean(manifestPath string, keyPath string) error {
	if err := removeAll(manifestPath); err != nil {
		return err
	}
	
	if err := removeAll(keyPath); err != nil {
		return err
	}
	return nil
}
