package gobnb_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	bnb "gitlab.com/gobnb"
)

type SimpleProblem struct {
	State        SimpleProblemState
	InitialState SimpleProblemState
}

type SimpleProblemState struct {
	UpperBound float64
	LowerBound float64
}

func (s *SimpleProblem) Sense() bnb.ProblemSense {
	return bnb.Minimize
}
func (s *SimpleProblem) Objective(n *bnb.Node) (bound float64) {
	state := &SimpleProblemState{}
	n.LoadState(state)

	u, l := state.UpperBound, state.LowerBound
	return u - l
}
func (s *SimpleProblem) Bound(n *bnb.Node) float64 {
	state := &SimpleProblemState{}
	n.LoadState(state)

	u, l := state.UpperBound, state.LowerBound
	return math.Pow(l-u, 2)
}

func (s *SimpleProblem) LoadInitialNode() *bnb.Node {
	upper, lower := s.InitialState.UpperBound, s.InitialState.LowerBound
	initialNode := &bnb.Node{State: SimpleProblemState{UpperBound: upper, LowerBound: lower}}
	return initialNode
}

func (s *SimpleProblem) Branch(n *bnb.Node, currentBound float64) []*bnb.Node {
	state := &SimpleProblemState{}
	n.LoadState(state)

	upper, lower := state.UpperBound, state.LowerBound
	mean := (lower + upper) / 2

	// possibly generate go coroutine to generate new state in a
	// parallel way
	var newNodes []*bnb.Node

	lowerState := SimpleProblemState{LowerBound: lower, UpperBound: mean}
	lowerNode := &bnb.Node{State: lowerState}

	if s.Bound(lowerNode) <= currentBound {
		newNodes = append(newNodes, lowerNode)
	}

	upperState := SimpleProblemState{LowerBound: mean, UpperBound: upper}
	upperNode := &bnb.Node{State: upperState}
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

	solver := bnb.Solver{simpleProblem}
	_, _, bound, err := solver.Solve(bnb.SolverConfigs{Mode: bnb.DepthFirst})
	assert.NoError(t, err, "Solver should not raise error")
	assert.LessOrEqual(t, bound, math.Pow10(-6), "Bound to target reached a 10-6")
}
