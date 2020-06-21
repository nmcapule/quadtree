package main

import (
	"flag"
	"fmt"

	"github.com/nmcapule/quadtree"
)

var redisConnFlag = flag.String("redis", "redis://localhost:6379", "Redis connection string")

func main() {
	flag.Parse()

	backend, err := quadtree.NewRedisBackend(*redisConnFlag)
	if err != nil {
		panic(err)
	}

	bounds := quadtree.Bounds{0., 0., 256., 256.}
	qt := quadtree.NewQuadtree(bounds)
	qt.Insert(&quadtree.Object{1, 1, 100})

	fmt.Println(qt)
	for _, obj := range qt.FindAllWithin(bounds) {
		fmt.Println(obj)
	}

	if err := backend.SetQuadtree(qt); err != nil {
		panic(err)
	}
	if qqt, err := backend.GetQuadtree(qt.UUID); err != nil {
		panic(err)
	} else {
		fmt.Println(qqt)
	}
}
