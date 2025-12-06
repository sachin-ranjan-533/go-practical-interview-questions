package main

import (
	"cache/cache"
	"fmt"
)

func main() {
	c := cache.NewCache()

	fn := func() any {
		fmt.Println("Computing value...")
		return "Hello, Cache!"
	}

	v1, ok1 := c.GetOrCompute("greeting", fn)
	fmt.Println(v1, ok1)

	v2, ok2 := c.GetOrCompute("greeting", fn)
	fmt.Println(v2, ok2)
}
