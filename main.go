package main

import "context"

func main() {
	ctx := context.Background()
	ensureDir("data")
	showIntro(ctx)
}
