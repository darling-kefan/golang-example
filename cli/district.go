package main

import (
	"fmt"
	"districts"
)

func main() {
	dists, err := districts.Source()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := districts.ToRedis(dists); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Generate successfully!")
}
