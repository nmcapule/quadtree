package quadtree

import (
	"errors"
	"fmt"
)

var (
	errAlreadySplit      = errors.New("quadtree is already split")
	errInsertOutOfBounds = errors.New("object to be inserted is out-of-bounds")
	errInsertFailure     = errors.New("object insert failure")
)

// Object encapsulates any object inserted into the quadtree.
type Object struct {
	X, Y float64
	Data interface{}
}

// Bounds is a rectangle.
type Bounds struct {
	X, Y          float64
	Width, Height float64
}

// Within checks if the input coordinate is within the bounds.
func (b *Bounds) Within(x, y float64) bool {
	return x >= b.X && x < b.X+b.Width && y >= b.Y && y < b.Y+b.Height
}

// Intersects checks if two bounds intersect.
func (b *Bounds) Intersects(other Bounds) bool {
	return b.X+b.Width > other.X &&
		b.Y+b.Height > other.Y &&
		b.X < other.X+other.Width &&
		b.Y < other.Y+other.Height
}

// Quadtree is a naive implementation of a quadtree.
type Quadtree struct {
	Bounds     Bounds
	Nodes      []*Quadtree
	Objects    []*Object
	Level      int
	MaxObjects int
}

// NewQuadtree instantiates a new quadtree from a given bounds.
func NewQuadtree(bounds Bounds) *Quadtree {
	return &Quadtree{
		Bounds:     bounds,
		MaxObjects: 5,
	}
}

// Insert inserts an object into the quadtree.
func (qt *Quadtree) Insert(object *Object) error {
	// If tried to insert out of bounds, immediately return an error.
	if !qt.Bounds.Within(object.X, object.Y) {
		return errInsertOutOfBounds
	}
	// If there's no children yet, try to insert simple.
	if len(qt.Nodes) == 0 {
		if len(qt.Objects) < qt.MaxObjects {
			qt.Objects = append(qt.Objects, object)
			return nil
		}
		// If can't insert, subdivide.
		err := qt.subdivide()
		if err != nil {
			return err
		}
	}
	for _, node := range qt.Nodes {
		err := node.Insert(object)
		// If insert is successful, terminate and return nil error.
		if err == nil {
			return nil
		}
		if err == errInsertOutOfBounds {
			continue
		}
		return err
	}
	return errInsertFailure
}

// FindAllWithin retrieves all objects that intersects the given bound within the quadtree.
func (qt *Quadtree) FindAllWithin(bounds Bounds) []*Object {
	// If this is a leaf node, do a simple iterate.
	if qt.Nodes == nil {
		var res []*Object
		for i, object := range qt.Objects {
			if !bounds.Within(object.X, object.Y) {
				continue
			}
			res = append(res, qt.Objects[i])
		}
		return res
	}

	var res []*Object
	for _, node := range qt.Nodes {
		if !bounds.Intersects(node.Bounds) {
			continue
		}
		res = append(res, node.FindAllWithin(bounds)...)
	}
	return res
}

func (qt *Quadtree) subdivide() error {
	if qt.Nodes != nil {
		return errAlreadySplit
	}

	qt.Nodes = make([]*Quadtree, 4)

	bx, by := qt.Bounds.X, qt.Bounds.Y
	hw, hh := qt.Bounds.Width/2, qt.Bounds.Height/2

	// Nodes index by quadrant: NW, NE, SE, SW
	qt.Nodes[0] = &Quadtree{
		Bounds:     Bounds{X: bx, Y: by, Width: hw, Height: hh},
		Level:      qt.Level + 1,
		MaxObjects: qt.MaxObjects,
	}
	qt.Nodes[1] = &Quadtree{
		Bounds:     Bounds{X: bx + hw, Y: by, Width: hw, Height: hh},
		Level:      qt.Level + 1,
		MaxObjects: qt.MaxObjects,
	}
	qt.Nodes[2] = &Quadtree{
		Bounds:     Bounds{X: bx + hw, Y: by + hh, Width: hw, Height: hh},
		Level:      qt.Level + 1,
		MaxObjects: qt.MaxObjects,
	}
	qt.Nodes[3] = &Quadtree{
		Bounds:     Bounds{X: bx, Y: by + hh, Width: hw, Height: hh},
		Level:      qt.Level + 1,
		MaxObjects: qt.MaxObjects,
	}

	objects := qt.Objects
	qt.Objects = nil
	for _, object := range objects {
		err := qt.Insert(object)
		if err != nil {
			return err
		}
	}

	return nil
}

// DebugPrint is a naive debug printing for the whole quadtree.
func (qt *Quadtree) DebugPrint() {
	var tabs string
	for i := 0; i < qt.Level; i++ {
		tabs += " "
	}
	fmt.Printf("%sNode (%05.2f, %05.2f): %d objects --- Level: %d\n", tabs, qt.Bounds.X, qt.Bounds.Y, len(qt.Objects), qt.Level)
	for _, node := range qt.Nodes {
		node.DebugPrint()
	}
}
