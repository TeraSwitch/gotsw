package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/teraswitch/gotsw/v2"
)

func main() {
	ctx := context.Background()

	client := gotsw.New(os.Getenv("TSW_API_KEY")).
		SetLogBodies(true).
		SetPlainLogger(os.Stdout).
		SetLogger(
			slog.New(slog.NewTextHandler(os.Stdout, nil)),
		)

	tiers, err := client.ListMetalTiers(ctx, gotsw.MetalTierTypeCompute)
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(tiers)
}
