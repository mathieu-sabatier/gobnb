package gobnb

import (
	"encoding/json"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

type ProblemSense string

const (
	Minimize ProblemSense = "minize"
	Maximize ProblemSense = "maximze"
)

type Node struct {
	State  any
	Depth  int
	Parent *Node
}

func (*Node) LoadState(a, b interface{}) error {
	js, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, b)
}

type Problem interface {
	Sense() ProblemSense
	Objective(*Node) float64
	Bound(*Node) float64
	Branch(*pq.Queue, *Node, float64) error
	LoadInitialNode() *Node
}
