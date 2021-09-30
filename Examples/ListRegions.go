package main

import (
	"context"
	"fmt"

	"github.com/teraswitch/gotsw"
)

func main() {
	ctx := context.Background()

	tswClient := gotsw.NewFromToken("id:secret")
	i, err := tswClient.RegionService.List(ctx)
	if err != nil {
		fmt.Println(err)
	}

	for _, v := range i {
		fmt.Println(v)
	}
}
