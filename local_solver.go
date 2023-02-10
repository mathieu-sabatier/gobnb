package gobnb

import (
	"fmt"
	"math"

	pq "github.com/emirpasic/gods/queues/priorityqueue"
)

type Solver struct {
	Problem Problem
}

func (s *Solver) Solve(config *SolverConfigs) (bestNode *Node, objective float64, bound float64, err error) {

	var bestbound, bestObjective float64
	bestObjective = math.Inf(1)

	initialNode := s.Problem.LoadInitialNode()
	bestbound = s.Problem.Bound(initialNode)

	comparator := newComparatorFromConfig(config)
	nodes := pq.NewWith(comparator)

	newNodes := s.Problem.Branch(initialNode, bestbound, bestObjective)
	for _, node := range newNodes {
		node = initialNode.iter(node)
		nodes.Enqueue(node)
	}

	if nodes.Size() == 0 {
		fmt.Println("could not find initial node")
		fmt.Println("early stop")
		return
	}

	stats := newStatsWriterFromConfig(config)
	checker := newConvergenceCheckerFromConfig(config)
	for {
		n, optimalReached := nodes.Dequeue()
		if !optimalReached {
			fmt.Println("optimality condition reached - no more nodes.")
			break
		}

		nextNode := n.(*Node)
		bound, objective = s.Problem.Bound(nextNode), s.Problem.Objective(nextNode)

		// update bound if objective is reached or when first objective is reached
		if ((bound < bestbound) && (!math.IsInf(objective, +1))) || ((math.IsInf(bestObjective, +1)) && (!math.IsInf(objective, +1))) {
			bestbound = bound
		}
		if objective < bestObjective {
			bestObjective = objective
			bestNode = nextNode
			stats.inform(nodes.Size(), bestObjective, bound)
		}

		convergenceReached := checker.iter(bound, objective)
		if convergenceReached {
			fmt.Println("external convergence reached.")
			fmt.Println("stopping from", checker.convergenceMode)
			break
		}

		newNodes := s.Problem.Branch(nextNode, bestbound, bestObjective)
		for _, node := range newNodes {
			node = nextNode.iter(node)
			nodes.Enqueue(node)
		}
		stats.iter(nodes.Size(), bestObjective, bestbound)
	}
	stats.terminate()
	return bestNode, bestObjective, bound, err
}
