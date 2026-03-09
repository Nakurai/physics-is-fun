package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Solid struct {
	Edges   []Edge
	Corners [4]Vector[float32] // list the outside box of the solid
}

func (o *Solid) Update() error {
	return nil
}

func (o *Solid) Draw(screen *ebiten.Image) {
	for _, edge := range o.Edges {
		vector.StrokeLine(screen, edge.A.X, edge.A.Y, edge.B.X, edge.B.Y, 1, Red, true)
	}
}

func (o Solid) GetLastEdge() *Edge {
	nbEdges := len(o.Edges)
	if nbEdges == 0 {
		return nil
	}
	return &o.Edges[nbEdges-1]
}

func (o Solid) ComputeEdges() {
	for i := range len(o.Edges) {
		o.Edges[i].Compute()
	}
}
