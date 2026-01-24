package main

import (
	"os"
	"errors"
	"context"
	"path/filepath"
	"github.com/hashicorp/go-extract"
)

func ensureDir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func env(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", errors.New("Environment variable not set: " + name)
	}
	return value, nil
}

func folderSize(path string) (int64, error) {
    var size int64
    err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            size += info.Size()
        }
        return nil
    })

    if err != nil {
        return 0, err
    }
    
    return size, nil
}

func extractArchive(path string, ctx context.Context) error {
	dir := filepath.Dir(path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	
	if err = extract.Unpack(ctx, dir, f, extract.NewConfig()); err != nil {
		return err
	}
	
	if err = os.Remove(path); err != nil {
		return err
	}
	
	if err = os.Rename(filepath.Join(dir, "z"), path); err != nil {
		return err
	}
	return nil
}

func removeAll(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}
