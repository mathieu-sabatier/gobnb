package gobnb

import (
	"fmt"
	"math"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
	"github.com/emirpasic/gods/utils"
)

type Solver struct {
	Problem  Problem
	LastNode *Node
	BestNode *Node
}

type SolverMode uint

const (
	DepthFirst SolverMode = iota
	BreadthFirst
)

var SolverModeNames = [...]string{"DepthFirst", "BreathFirst"}

func (mode SolverMode) String() string {
	return SolverModeNames[mode]
}

type SolverConfigs struct {
	AbsoluteGap  float64
	MaxSpentTime int64
	Mode         SolverMode
}

func newComparatorFromConfig(mode SolverMode) utils.Comparator {
	if mode == DepthFirst {
		return func(a, b interface{}) int {
			priorityA := a.(*Node).Depth
			priorityB := b.(*Node).Depth
			return -utils.IntComparator(priorityA, priorityB)
		}
	} else if mode == BreadthFirst {
		return func(a, b interface{}) int {
			priorityA := a.(*Node).Depth
			priorityB := b.(*Node).Depth
			return -utils.IntComparator(priorityA, priorityB)
		}
	}
	panic(mode)
}

func (s *Solver) Solve(config SolverConfigs) (objective float64, bound float64, err error) {

	var bestbound, bestObjective float64
	bestObjective = math.Inf(1)

	initialNode := s.Problem.LoadInitialNode()
	bestbound = s.Problem.Bound(initialNode)

	comparator := newComparatorFromConfig(config.Mode)
	nodes := pq.NewWith(comparator)

	s.Problem.Branch(nodes, initialNode, bestbound)

	termination := 0
	for termination <= 100 {
		termination += 1

		n, ok := nodes.Dequeue()
		if !ok {
			fmt.Println("queue is empty after", termination, "iteration.")
			break
		}

		nextNode := n.(*Node)
		bound, objective = s.Problem.Bound(nextNode), s.Problem.Objective(nextNode)

		// check if termination is reached
		if math.Abs(bound) <= math.Pow10(-6) {
			fmt.Println("early stop: bound condition reached after", termination, "iteration.")
			fmt.Println("current bound: ", bound)
			break
		}

		// update bound and objective
		if bound > bestbound {
			bestbound = bound
		}
		if objective < bestObjective {
			bestObjective = objective
			fmt.Println("improving objective to", bestObjective)
		}

		s.Problem.Branch(nodes, nextNode, bestbound)
	}
	return bestObjective, bound, err
}
