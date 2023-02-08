package gobnb

import (
	"fmt"
	"math"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

type Solver struct {
	Problem Problem
}

func (s *Solver) Solve(config SolverConfigs) (bestNode *Node, objective float64, bound float64, err error) {

	var bestbound, bestObjective float64
	bestObjective = math.Inf(1)

	initialNode := s.Problem.LoadInitialNode()
	bestbound = s.Problem.Bound(initialNode)

	comparator := newComparatorFromConfig(config)
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
		if bound < bestbound {
			bestbound = bound
		}
		if objective < bestObjective {
			fmt.Println("improving objective from", bestObjective, "to", objective)
			bestObjective = objective
			bestNode = nextNode
		}

		s.Problem.Branch(nodes, nextNode, bestbound)
	}
	return bestNode, bestObjective, bound, err
}
