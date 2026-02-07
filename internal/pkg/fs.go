package pkg

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-extract"
)

func EnsureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteDir(path string) error {
	return os.RemoveAll(path)
}

func Env(name string) (string, error) {
	v := os.Getenv(name)
	if v == "" {
		return "", errors.New("env var not found")
	}
	return v, nil
}

func WriteFile(path string, data []byte) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	if err = os.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func Extract(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	if err = extract.Unpack(context.Background(), filepath.Dir(path), f, extract.NewConfig()); err != nil {
		return err
	}

	if err = os.Remove(path); err != nil {
		return err
	}
	return nil
}

func Execute(path string, args [][]string) error {
	flat := []string{}

	for _, pair := range args {
		if len(pair) != 2 {
			return errors.New("invalid arg")
		}
		flat = append(flat, pair[0], pair[1])
	}

	cmd := exec.Command(
		path,
		flat...,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
