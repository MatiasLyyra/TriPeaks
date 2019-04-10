package mcts

import "github.com/MatiasLyyra/TriPeaks/deck"

type Deter struct {
	Pos         int
	Card        deck.Card
	Initialized bool
}

type NodeData struct {
	CardsLeft          []deck.Card
	CardsLeftBeginning int
}

type Node struct {
	X        float64
	N        int
	Pos      int
	LeftDet  Deter
	RightDet Deter
	Parent   *Node
	Children []*Node
	Data     *NodeData
}

func (n *Node) GetUnvisitedChild() *Node {
	unvisited := make([]*Node, 0)
	for _, child := range n.Children {
		if child.N == 0 {
			unvisited = append(unvisited, child)
		}
	}
	if len(unvisited) == 0 {
		return nil
	}
	return unvisited[0]
}

func (n *Node) ChildPos(pos int) int {
	for i, child := range n.Children {
		if child.Pos == pos {
			return i
		}
	}
	return -1
}

func (n *Node) GetParentDeterminization(pos int, left bool) bool {
	if n.Parent == nil {
		return false
	}
	assignL := func(p *Node) {
		if left {
			n.LeftDet = p.LeftDet
		} else {
			n.RightDet = p.LeftDet
		}
	}
	assignR := func(p *Node) {
		if left {
			n.LeftDet = p.RightDet
		} else {
			n.RightDet = p.RightDet
		}
	}
	for parent := n.Parent; parent != nil; parent = parent.Parent {
		if parent.LeftDet.Pos == pos {
			assignL(parent)
			return true
		} else if parent.RightDet.Pos == pos {
			assignR(parent)
			return true
		}
	}
	return false
}

func NewNode() *Node {
	return &Node{
		X:        0,
		N:        0,
		Pos:      -2,
		Parent:   nil,
		Children: make([]*Node, 0, 5),
	}
}
