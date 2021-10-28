package main

import (
	"fmt"
	"github.com/cryptowilliam/goutil/webapi/github"
)

func main() {

	gd, err := github.NewGitDrive("", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gd.CreatePartition("test2"))
	fmt.Println(gd.Partitions())
	fmt.Println(gd.Partition("test2").LastUpdateTime())
}
