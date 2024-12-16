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

	resp, err := client.ReinstallMetalService(ctx, 10346, &gotsw.ReinstallMetalRequest{
		DisplayName: "test-metal-service",
		ImageID:     "ubuntu-noble",
		SSHKeyIDs:   []int{588},
		Partitions: []gotsw.Partition{
			{
				Name:   "nvme0n1-part1",
				Device: "nvme0n1",
			},
			{
				Name:   "nvme1n1-part1",
				Device: "nvme1n1",
			},
		},
		RaidArrays: []gotsw.RaidArray{
			{
				Name:       "md0",
				Type:       gotsw.RaidTypeRaid1,
				Members:    []string{"nvme0n1-part1", "nvme1n1-part1"},
				FileSystem: gotsw.FileSystemExt4,
				MountPoint: "/",
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	spew.Dump(resp)
}
