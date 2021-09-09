package main

import (
	"ConcurrentSearch/routers"
	"fmt"
)

func main() {
	r := routers.SetRouter()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}
