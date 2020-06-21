package quadtree

import (
	"math/rand"
	"time"

	"testing"
)

func createQuadtree(n int, bound Bounds) *Quadtree {
	qt := NewQuadtree(bound)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < n; i++ {
		err := qt.Insert(&Object{
			X:    bound.X + rand.Float64()*bound.Width,
			Y:    bound.Y + rand.Float64()*bound.Height,
			Data: i,
		})
		if err != nil {
			panic(err)
		}
	}

	return qt
}

func benchmarkFindAllWithin(qt *Quadtree, b *testing.B) {
	for n := 0; n < b.N; n++ {
		// x, y := 0., 0.
		// w, h := 100., 100.
		x := qt.Bounds.X + rand.Float64()*qt.Bounds.Width
		y := qt.Bounds.Y + rand.Float64()*qt.Bounds.Height
		w := 1 * qt.Bounds.Width / 100
		h := 1 * qt.Bounds.Height / 100
		qt.FindAllWithin(Bounds{x, y, w, h})
	}
}

func BenchmarkQuadtreeFind(b *testing.B) {
	qt := createQuadtree(1000, Bounds{0, 0, 4096, 4096})
	qt.MaxLevel = 5
	b.ResetTimer()
	benchmarkFindAllWithin(qt, b)
}

func benchmarkCreate(n int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		createQuadtree(n, Bounds{0, 0, 100, 100})
	}
}

func BenchmarkQuadtreeCreate100(b *testing.B) {
	benchmarkCreate(100, b)
}

func benchmarkMove(qt *Quadtree, objects []*Object, bounds Bounds, b *testing.B) {
	for n := 0; n < b.N; n++ {
		x := bounds.X + rand.Float64()*bounds.Width
		y := bounds.Y + rand.Float64()*bounds.Height
		qt.Move(objects[n%qt.Total], x, y)
	}
}

func BenchmarkQuadtreeMove(b *testing.B) {
	bounds := Bounds{0, 0, 4096, 4096}
	qt := createQuadtree(1000000, bounds)
	objects := qt.FindAllWithin(bounds)
	b.ResetTimer()
	benchmarkMove(qt, objects, bounds, b)
}

func TestInsertRemove(t *testing.T) {
	bounds := Bounds{0, 0, 100, 100}
	n := 5000
	qt := createQuadtree(n, bounds)
	if qt.Total != n {
		t.Errorf("Wrong number of total objects: want %d, got %d", n, qt.Total)
	}
	objects := qt.FindAllWithin(bounds)
	if len(objects) != n {
		t.Errorf("Wrong number of retrieved objects: want %d, got %d", n, len(objects))
	}
	for _, obj := range objects {
		x := bounds.X + rand.Float64()*bounds.Width
		y := bounds.Y + rand.Float64()*bounds.Height
		qt.Move(obj, x, y)
	}
	objects = qt.FindAllWithin(bounds)
	if len(objects) != n {
		t.Errorf("Wrong number of retrieved objects after move: want %d, got %d", n, len(objects))
	}
	for i, obj := range objects {
		err := qt.Remove(obj)
		if err != nil {
			t.Errorf("Unexpected error on object Remove: %v", err)
		}
		if qt.Total != n-i-1 {
			t.Errorf("Wrong number of total objects after removal: want %d, got %d", n-i-1, qt.Total)
		}
	}
	if qt.Nodes != nil {
		t.Errorf("Expecting child nodes to be nil, got %+v", qt.Nodes)
	}
}
