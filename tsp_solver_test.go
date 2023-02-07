package gobnb_test

import (
	"math"
	"testing"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/stretchr/testify/assert"
	bnb "gitlab.com/gobnb"
)

type TravellingSalespersonProblem struct {
	State TravellingSalespersonProblemState
}
type TravellingSalespersonProblemState struct {
	nSalesman int
}

func (s *TravellingSalespersonProblem) Sense() bnb.ProblemSense {
	return bnb.Minimize
}
func (s *TravellingSalespersonProblem) Objective(n *bnb.Node) (bound float64) {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(n.State, state)

	if n.Depth != state.nSalesman {
		return math.Inf(1)
	}

	// distanceMatrix := n.State["distanceMatrix"].(float32)
	// u, l := n.State["upperBound"].(float64), n.State["lowerBound"].(float64)
	return 0
}
func (s *TravellingSalespersonProblem) Bound(n *bnb.Node) float64 {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(n.State, state)
	return float64(state.nSalesman)
}

func (s *TravellingSalespersonProblem) LoadInitialNode() *bnb.Node {
	initialNode := &bnb.Node{State: TravellingSalespersonProblemState{}, Depth: 1}
	return initialNode
}

func (s *TravellingSalespersonProblem) Branch(q *pq.Queue, n *bnb.Node, currentBound float64) error {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(n.State, state)

	lowerNode := &bnb.Node{State: state, Parent: n, Depth: n.Depth + 1}

	if s.Bound(lowerNode) >= currentBound {
		q.Enqueue(lowerNode)
	}

	return nil
}

func TestSolverTSP(t *testing.T) {
	simpleProblem := &SimpleProblem{
		State: map[string]any{
			"lowerBound": 0.0,
			"upperBound": 1.0,
		},
	}

	solver := bnb.Solver{simpleProblem, nil, nil}
	_, bound, err := solver.Solve(bnb.SolverConfigs{Mode: bnb.DepthFirst})
	assert.NoError(t, err, "Solver should not raise error")
	assert.LessOrEqual(t, bound, math.Pow10(-6), "Bound to target reached a 10-6")
}
