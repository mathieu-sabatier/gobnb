package gobnb

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SimpleProblem struct {
	State        SimpleProblemState
	InitialState SimpleProblemState
}

type SimpleProblemState struct {
	UpperBound float64
	LowerBound float64
}

func (s *SimpleProblem) Sense() ProblemSense {
	return Minimize
}
func (s *SimpleProblem) Objective(n *Node) (bound float64) {
	state := &SimpleProblemState{}
	n.LoadState(state)

	u, l := state.UpperBound, state.LowerBound
	return u - l
}
func (s *SimpleProblem) Bound(n *Node) float64 {
	state := &SimpleProblemState{}
	n.LoadState(state)

	u, l := state.UpperBound, state.LowerBound
	return math.Pow(l-u, 2)
}

func (s *SimpleProblem) LoadInitialNode() *Node {
	upper, lower := s.InitialState.UpperBound, s.InitialState.LowerBound
	initialNode := &Node{State: SimpleProblemState{UpperBound: upper, LowerBound: lower}}
	return initialNode
}

func (s *SimpleProblem) Branch(n *Node, currentBound float64) []*Node {
	state := &SimpleProblemState{}
	n.LoadState(state)

	upper, lower := state.UpperBound, state.LowerBound
	mean := (lower + upper) / 2

	// possibly generate go coroutine to generate new state in a
	// parallel way
	var newNodes []*Node

	lowerState := SimpleProblemState{LowerBound: lower, UpperBound: mean}
	lowerNode := &Node{State: lowerState}

	if s.Bound(lowerNode) <= currentBound {
		newNodes = append(newNodes, lowerNode)
	}

	upperState := SimpleProblemState{LowerBound: mean, UpperBound: upper}
	upperNode := &Node{State: upperState}
	if s.Bound(upperNode) <= currentBound {
		newNodes = append(newNodes, upperNode)
	}

	return newNodes
}

func TestSolver(t *testing.T) {
	simpleProblem := &SimpleProblem{
		InitialState: SimpleProblemState{
			LowerBound: 0.0,
			UpperBound: 1.0,
		},
	}

	solver := Solver{simpleProblem}
	_, _, bound, err := solver.Solve(&SolverConfigs{Mode: DepthFirst, MaxIterCount: 100})
	assert.NoError(t, err, "Solver should not raise error")
	assert.LessOrEqual(t, bound, math.Pow10(-6), "Bound to target reached a 10-6")
}
