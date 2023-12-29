package main

import (
	"fmt"
	"time"

	"github.com/raynine/pub-sub/pubSub"
)

func main() {
	start := time.Now()
	fmt.Println("Starting Pub Sub")

	pubSub.Run()

	fmt.Println("Completed in: ", time.Since(start))
}
