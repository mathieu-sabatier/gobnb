package gobnb

import (
	"encoding/json"
)

type ProblemSense string

const (
	Minimize ProblemSense = "minize"
	Maximize ProblemSense = "maximze"
)

type Node struct {
	State  any
	depth  int
	parent *Node
}

func (n *Node) LoadState(target interface{}) error {
	js, err := json.Marshal(n.State)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, target)
}

func (n *Node) iter(nextNode *Node) *Node {
	if n.parent == nil {
		n.depth = 1
	} else {
		n.depth = n.parent.depth + 1
	}
	nextNode.parent = n
	return nextNode
}

type Problem interface {
	Sense() ProblemSense
	Objective(*Node) float64
	Bound(*Node) float64
	Branch(*Node, float64) []*Node
	LoadInitialNode() *Node
}
