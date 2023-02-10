package tsp

import (
	"math"

	"gitlab.com/gobnb"
	"gonum.org/v1/gonum/mat"
)

type TravellingSalespersonProblem struct {
	State          TravellingSalespersonProblemState
	DistanceMatrix *mat.Dense
	NSalesman      int
}
type TravellingSalespersonProblemState struct {
	Sequence    []int
	CurrentCost float64
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (s *TravellingSalespersonProblem) Sense() gobnb.ProblemSense {
	return gobnb.Minimize
}
func (s *TravellingSalespersonProblem) Objective(n *gobnb.Node) (bound float64) {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(state)

	depth := len(state.Sequence)

	if depth != s.NSalesman {
		return math.Inf(1)
	}

	// compute distances
	distances := s.DistanceMatrix
	objective := state.CurrentCost
	first_point, last_point := state.Sequence[0], state.Sequence[depth-1]
	objective += distances.At(last_point, first_point)
	return objective
}
func (s *TravellingSalespersonProblem) Bound(n *gobnb.Node) float64 {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(state)

	if len(state.Sequence) == s.NSalesman {
		return s.Objective(n)
	}

	// as a first bound, we take the minimal distance the traveler will do
	distances := s.DistanceMatrix
	nRemainingDistances := s.NSalesman - len(state.Sequence) + 1

	var distance float64
	min_distance := math.Inf(+1)
	for distance_from := 0; distance_from < s.NSalesman; distance_from++ {
		for distance_to := 0; distance_to < s.NSalesman; distance_to++ {
			if distance_from == distance_to {
				continue
			}
			if !contains(state.Sequence, distance_to) {
				distance = distances.At(distance_from, distance_to)
				if distance < min_distance {
					min_distance = distance
				}
			}
		}
	}

	if min_distance == math.Inf(+1) {
		return math.Inf(-1)
	}

	return state.CurrentCost + float64(nRemainingDistances)*float64(min_distance)
}

func (s *TravellingSalespersonProblem) LoadInitialNode() *gobnb.Node {
	initialNode := &gobnb.Node{State: TravellingSalespersonProblemState{Sequence: []int{0}}}
	return initialNode
}

func (s *TravellingSalespersonProblem) Branch(n *gobnb.Node, currentBound float64, bestObjectiveReached float64) []*gobnb.Node {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(state)
	distances := s.DistanceMatrix

	activeSequence := state.Sequence
	lastPassage := state.Sequence[len(activeSequence)-1]

	var nextNodes []*gobnb.Node
	for passage := 0; passage < s.NSalesman; passage++ {
		if contains(activeSequence, passage) {
			continue
		}

		newSequence := make([]int, len(activeSequence)+1)
		_ = copy(newSequence, activeSequence)
		newSequence[len(activeSequence)] = passage

		newNode := &gobnb.Node{
			State: TravellingSalespersonProblemState{
				Sequence:    newSequence,
				CurrentCost: state.CurrentCost + distances.At(lastPassage, passage),
			},
		}
		local_bound := s.Bound(newNode)
		if ((local_bound <= currentBound) && (!math.IsInf(bestObjectiveReached, +1))) || math.IsInf(bestObjectiveReached, +1) {
			nextNodes = append(nextNodes, newNode)
		}
	}

	return nextNodes
}
