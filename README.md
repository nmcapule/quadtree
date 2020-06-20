# Quadtree

> Not invented here is alive and well. It is shitty, but it is mine :D

https://pkg.go.dev/github.com/nmcapule/quadtree

## Usage

```go
qt := quadtree.NewQuadtree(quadtree.Bounds{
    X: 0, Y: 0,
    Width: 100, Height: 100,
})
qt.Insert(quadtree.Object{
    X: 20, Y: 25, Data: "hello",
})
qt.Insert(quadtree.Object{
    X: 80, Y: 95, Data: "world",
})

fmt.Println(qt.FindAllWithin(quadtree.Bounds{
    X: 10, Y: 10, Width: 20, Height: 20,
}))
```
