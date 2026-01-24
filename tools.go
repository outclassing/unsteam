package main

import (
	"context"
	"errors"
	"bytes"
	"os/exec"
)

func toolsMenu(ctx context.Context) {

}

func installDepotDownloader(ctx context.Context) {
	// os aware; select
}

func runDepotDownloader(appId string, keyFile string, depotId string, manifestId string, manifestFile string) error {
	cmd := exec.Command(
		"tools/depotdownloader/linux/DepotDownloaderMod",
		"-app", appId,
		"-depotkeys", keyFile,
		"-depot", depotId,
		"-manifest", manifestId,
		"-manifestfile", manifestFile,
	)
	manifestFile = manifestFile

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	if err := cmd.Run(); err != nil {
		return errors.New(err.Error() + stderr.String())
	}
	return nil
}
