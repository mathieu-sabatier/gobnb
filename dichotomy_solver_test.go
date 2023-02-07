package gobnb_test

import (
	"math"
	"testing"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/stretchr/testify/assert"
	bnb "gitlab.com/gobnb"
)

type SimpleProblem struct {
	State map[string]any
}

func (s *SimpleProblem) Sense() bnb.ProblemSense {
	return bnb.Minimize
}
func (s *SimpleProblem) Objective(n *bnb.Node) (bound float64) {
	u, l := n.State["upperBound"].(float64), n.State["lowerBound"].(float64)
	return u - l
}
func (s *SimpleProblem) Bound(n *bnb.Node) float64 {
	u, l := n.State["upperBound"].(float64), n.State["lowerBound"].(float64)
	return -math.Pow(l-u, 2)
}

func (s *SimpleProblem) LoadInitialNode() *bnb.Node {
	upper, lower := s.State["upperBound"].(float64), s.State["lowerBound"].(float64)

	newState := make(map[string]any)

	newState["lowerBound"] = lower
	newState["upperBound"] = upper
	initialNode := &bnb.Node{State: newState, Depth: 1}
	return initialNode
}

func (s *SimpleProblem) Branch(q *pq.Queue, n *bnb.Node, currentBound float64) error {

	upper, lower := n.State["upperBound"].(float64), n.State["lowerBound"].(float64)
	mean := (lower + upper) / 2

	// possibly generate go coroutine to generate new state in a
	// parallel way

	lowerState := make(map[string]any)
	lowerState["lowerBound"] = lower
	lowerState["upperBound"] = mean
	lowerNode := &bnb.Node{State: lowerState, Parent: n, Depth: n.Depth + 1}

	if s.Bound(lowerNode) >= currentBound {
		q.Enqueue(lowerNode)
	}

	upperState := make(map[string]any)
	upperState["lowerBound"] = mean
	upperState["upperBound"] = upper
	upperNode := &bnb.Node{State: upperState, Parent: n, Depth: n.Depth + 1}
	if s.Bound(upperNode) >= currentBound {
		q.Enqueue(upperNode)
	}

	return nil
}

func TestSolver(t *testing.T) {
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
