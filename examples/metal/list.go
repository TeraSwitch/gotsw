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

	resp, err := client.ListMetal(ctx, gotsw.ListMetalOptions{
		Limit:  1,
		Region: "SLC1",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(resp)
}
