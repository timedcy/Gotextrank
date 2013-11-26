package textRank

import (
	"fmt"
)

type Vertex int

type Edge struct {
	index  Vertex
	weight float64
}

// MyError 实现Error的接口
type MyError struct {
	What string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%v", e.What)
}

type Graph struct {
	vertexCount   Vertex
	maxVertex     int
	OutBoundEdges map[Vertex][]Edge //Out(V) queries
	InBoundEdges  map[Vertex][]Edge //In(V) queries
}

func NewGraph(maxVertex int) *Graph {
	return &Graph{0, maxVertex, make(map[Vertex][]Edge), make(map[Vertex][]Edge)}
}

func (g *Graph) AddVertex() Vertex {
	o := g.vertexCount
	g.vertexCount += 1
	return o
}

func (g *Graph) VertexCount() int {
	return int(g.vertexCount)
}

//这是一个有向图的加边方法
func (g *Graph) AddEdge(from, to Vertex, weight float64) {
	if g.OutBoundEdges[from] == nil {
		g.OutBoundEdges[from] = make([]Edge, 1, g.maxVertex+1)
		g.OutBoundEdges[from][0] = Edge{to, weight}
	} else {
		g.OutBoundEdges[from] = append(g.OutBoundEdges[from], Edge{to, weight})
	}
	if g.InBoundEdges[to] == nil {
		g.InBoundEdges[to] = make([]Edge, 1, g.maxVertex+1)
		g.InBoundEdges[to][0] = Edge{from, weight}
	} else {
		g.InBoundEdges[to] = append(g.InBoundEdges[to], Edge{from, weight})
	}
}

//给边加权
func (g *Graph) AddEdgeWeight(from, to Vertex, weight float64) {
	if g.OutBoundEdges[from] == nil {
		panic("from Vertex has no edge.")
	}
	if g.InBoundEdges[to] == nil {
		panic("in Vertex has no edge.")
	}
	//only outBoundEges are used , so just add the out weight
	for i := range g.OutBoundEdges[from] {
		if g.OutBoundEdges[from][i].index == to {
			g.OutBoundEdges[from][i].weight += weight
		}
	}

}

//判断有无边,构造的时候防止重复构造边.有重复的边就一直加权
func (g *Graph) HasEdge(from, to Vertex) bool {
	if g.OutBoundEdges[from] == nil || g.InBoundEdges[to] == nil {
		return false
	}
	for _, k := range g.OutBoundEdges[from] {
		if k.index == to {
			return true
		}
	}
	return false

}

func (g *Graph) In(v Vertex) []Edge {
	return g.InBoundEdges[v]
}

func (g *Graph) Out(v Vertex) []Edge {
	return g.OutBoundEdges[v]
}

func (g *Graph) Weight(from Vertex, to Vertex) (weight float64, e error) {
	for _, edge := range g.OutBoundEdges[from] {
		if edge.index == to {
			return edge.weight, nil
		}
	}
	return float64(0), MyError{"Weight not found"}
}
