package quadtree

import (
	"math/rand"
	"time"

	"testing"
)

var globalQt *Quadtree

func init() {
	globalQt = createQuadtree(1000000)
}

func createQuadtree(n int) *Quadtree {
	qt := NewQuadtree(Bounds{0., 0., 100., 100.})

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		err := qt.Insert(&Object{
			X:    rand.Float64() * 100,
			Y:    rand.Float64() * 100,
			Data: i,
		})
		if err != nil {
			panic(err)
		}
	}

	return qt
}

func benchmarkCreate(n int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		createQuadtree(n)
	}
}

func benchmarkFindAllWithin(qt *Quadtree, b *testing.B) {
	for n := 0; n < b.N; n++ {
		x1 := rand.Float64() * 50
		y1 := rand.Float64() * 50
		x2 := x1 + rand.Float64()*50
		y2 := y1 + rand.Float64()*50
		qt.FindAllWithin(Bounds{x1, y1, x2, y2})
	}
}

func BenchmarkQuadtreeFind(b *testing.B) {
	benchmarkFindAllWithin(globalQt, b)
}

func BenchmarkQuadtreeCreate1000(b *testing.B) {
	benchmarkCreate(1000, b)
}
func BenchmarkQuadtreeCreate10000(b *testing.B) {
	benchmarkCreate(10000, b)
}
func BenchmarkQuadtreeCreate100000(b *testing.B) {
	benchmarkCreate(100000, b)
}
func BenchmarkQuadtreeCreate1000000(b *testing.B) {
	benchmarkCreate(1000000, b)
}
