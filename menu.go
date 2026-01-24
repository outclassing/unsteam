package main

import (
	"context"
	"fmt"
	"time"
	"errors"
	"github.com/yarlson/tap"
)

type Stream struct {
	*tap.Stream
}

func showIntro(ctx context.Context) {
	tap.Intro("untsteam v0.0.1")
	mainMenu(ctx)
}

func mainMenu(ctx context.Context) {
	switch tap.Select(ctx, tap.SelectOptions[string]{
		Message: "Choose an action:",
		Options: []tap.SelectOption[string]{
			{Value: "download", Label: "Download depot"},
			{Value: "patch", Label: "Install patch"},
			{Value: "update", Label: "Check for updates"},
			{Value: "tools", Label: "Install tool"},
			{Value: "exit", Label: "Exit"},
		},
	}) {
	case "download":
		downloadDepot(ctx)
	case "exit":
		tap.Outro("Exiting")
	}
}

func printDepotInfo(s *Stream, d Depot) {
	s.WriteLine(fmt.Sprintf("Manifest %s selected", d.ManifestId))
	s.WriteLine(fmt.Sprintf("Build ID: %s", fmt.Sprint(d.BuildId)))
	s.WriteLine(fmt.Sprintf("Time updated: %s", time.Unix(d.Time, 0).Format(time.RFC3339)))
}

func NewStream(msg string) *Stream {
	s := tap.NewStream(tap.StreamOptions{ShowTimer: true})
	s.Start(msg)
	return &Stream{s}
}

func (s *Stream) Ok(msg string) {
	s.Stop(msg, 0)
}

func (s *Stream) Error(msg string) {
	s.Stop(msg, -1)
}

func promptNumeric(ctx context.Context, label string, def string) string {
	return tap.Text(ctx, tap.TextOptions{
		Message:      label + ":",
		Placeholder: def,
		Validate: func(s string) error {
			if s == "" || !isNumeric(s) {
				return errors.New("numbers only")
			}
			return nil
		},
	})
}

func confirmation(ctx context.Context, msg string) bool {
	return tap.Confirm(ctx, tap.ConfirmOptions{
		Message: msg,
	})
}

func message(msg string) {
	tap.Message(msg)
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
