package quadtree

import (
	"math/rand"
	"time"

	"testing"
)

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
		// x, y := 0., 0.
		// w, h := 100., 100.
		x := rand.Float64() * 90
		y := rand.Float64() * 90
		w := 10. //rand.Float64() * 5
		h := 10. //rand.Float64() * 5
		qt.FindAllWithin(Bounds{x, y, w, h})
	}
}

func BenchmarkQuadtreeFind(b *testing.B) {
	qt := createQuadtree(10000)
	b.ResetTimer()
	benchmarkFindAllWithin(qt, b)
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
