package main

import (
	"context"
	"unsteam/internal/app"
	"unsteam/internal/pkg"
)

func main() {
	ctx := context.Background()
	pkg.EnsureDir("data")
	app.ShowIntro(ctx)
}
