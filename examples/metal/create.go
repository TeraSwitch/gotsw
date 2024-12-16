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

	resp, err := client.CreateMetalService(ctx, 480, &gotsw.CreateBareMetalRequest{
		DisplayName: "test-metal-service",
		RegionID:    "LAX1",
		TierID:      "2388g",
		MemoryGB:    64,
		ImageID:     "ubuntu-noble",
		SSHKeyIDs:   []int{588},
		Disks: map[string]string{
			"nvme0n1": "960g",
			"nvme1n1": "960g",
		},
		ReservePricing: false,
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
				Type:       gotsw.RaidTypeRaid0,
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
