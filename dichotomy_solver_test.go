package gobnb_test

import (
	"math"
	"testing"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
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
	n.LoadState(n.State, state)

	u, l := state.UpperBound, state.LowerBound
	return u - l
}
func (s *SimpleProblem) Bound(n *bnb.Node) float64 {
	state := &SimpleProblemState{}
	n.LoadState(n.State, state)

	u, l := state.UpperBound, state.LowerBound
	return math.Pow(l-u, 2)
}

func (s *SimpleProblem) LoadInitialNode() *bnb.Node {
	upper, lower := s.InitialState.UpperBound, s.InitialState.LowerBound
	initialNode := &bnb.Node{State: SimpleProblemState{UpperBound: upper, LowerBound: lower}, Depth: 1}
	return initialNode
}

func (s *SimpleProblem) Branch(q *pq.Queue, n *bnb.Node, currentBound float64) error {
	state := &SimpleProblemState{}
	n.LoadState(n.State, state)

	upper, lower := state.UpperBound, state.LowerBound
	mean := (lower + upper) / 2

	// possibly generate go coroutine to generate new state in a
	// parallel way

	lowerState := SimpleProblemState{LowerBound: lower, UpperBound: mean}
	lowerNode := &bnb.Node{State: lowerState, Parent: n, Depth: n.Depth + 1}

	if s.Bound(lowerNode) <= currentBound {
		q.Enqueue(lowerNode)
	}

	upperState := SimpleProblemState{LowerBound: mean, UpperBound: upper}
	upperNode := &bnb.Node{State: upperState, Parent: n, Depth: n.Depth + 1}
	if s.Bound(upperNode) <= currentBound {
		q.Enqueue(upperNode)
	}

	return nil
}

func TestSolver(t *testing.T) {
	simpleProblem := &SimpleProblem{
		InitialState: SimpleProblemState{
			LowerBound: 0.0,
			UpperBound: 1.0,
		},
	}

	solver := bnb.Solver{simpleProblem, nil, nil}
	_, _, bound, err := solver.Solve(bnb.SolverConfigs{Mode: bnb.DepthFirst})
	assert.NoError(t, err, "Solver should not raise error")
	assert.LessOrEqual(t, bound, math.Pow10(-6), "Bound to target reached a 10-6")
}
