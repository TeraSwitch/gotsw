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

	sshKeys, err := client.ListSshKeys(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(sshKeys)

	if len(sshKeys) == 0 {
		return
	}

	sshKey, err := client.GetSshKey(ctx, sshKeys[0].ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(sshKey)
}
