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

	sshKey, err := client.CreateSshKey(ctx, 480, gotsw.CreateSshKeyRequest{
		DisplayName: "test",
		Key:         "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIA56EsAhXB6bArh1gqN3sQwYGMPEJ5Y94mq+OSzYldVz colin@coder",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(sshKey)
}
