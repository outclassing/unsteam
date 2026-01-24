package app

import (
	"context"
	"errors"

	"github.com/yarlson/tap"
)

type Stream struct {
	*tap.Stream
}

func ShowIntro(ctx context.Context) {
	tap.Intro("unsteam v0.0.1")
	MainMenu(ctx)
}

func MainMenu(ctx context.Context) {
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
		DownloadDepot(ctx)
	case "exit":
		tap.Outro("Exiting")
	}
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

func PromptNumeric(ctx context.Context, label string, def string) string {
	return tap.Text(ctx, tap.TextOptions{
		Message:     label + ":",
		Placeholder: def,
		Validate: func(s string) error {
			if s == "" || !IsNumeric(s) {
				return errors.New("numbers only")
			}
			return nil
		},
	})
}

func Confirmation(ctx context.Context, msg string) bool {
	return tap.Confirm(ctx, tap.ConfirmOptions{
		Message: msg,
	})
}

func Message(msg string) {
	tap.Message(msg)
}

func IsNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
