package gobnb_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	bnb "gitlab.com/gobnb"
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

func (s *TravellingSalespersonProblem) Sense() bnb.ProblemSense {
	return bnb.Minimize
}
func (s *TravellingSalespersonProblem) Objective(n *bnb.Node) (bound float64) {
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
func (s *TravellingSalespersonProblem) Bound(n *bnb.Node) float64 {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(state)

	if len(state.Sequence) == 1 {
		return math.Inf(1)
	} else if len(state.Sequence) == s.NSalesman {
		return s.Objective(n)
	}

	// as a first bound, we take the minimal distance the traveler will do
	distances := s.DistanceMatrix
	nRemainingDistances := s.NSalesman - len(state.Sequence) + 1

	var max_distance, distance float64
	max_distance = math.Inf(-1)
	for distance_from := 0; distance_from < s.NSalesman; distance_from++ {
		for distance_to := 0; distance_to < s.NSalesman; distance_to++ {
			if distance_from == distance_to {
				continue
			}
			if !contains(state.Sequence, distance_to) {
				distance = distances.At(distance_from, distance_to)
				if max_distance < distance {
					max_distance = distance
				}
			}
		}
	}

	return state.CurrentCost + float64(nRemainingDistances)*float64(max_distance)
}

func (s *TravellingSalespersonProblem) LoadInitialNode() *bnb.Node {
	initialNode := &bnb.Node{State: TravellingSalespersonProblemState{Sequence: []int{0}}}
	return initialNode
}

func (s *TravellingSalespersonProblem) Branch(n *bnb.Node, currentBound float64) []*bnb.Node {
	state := &TravellingSalespersonProblemState{}
	n.LoadState(state)
	distances := s.DistanceMatrix

	activeSequence := state.Sequence
	lastPassage := state.Sequence[len(activeSequence)-1]

	var nextNodes []*bnb.Node
	for passage := 0; passage < s.NSalesman; passage++ {
		if contains(activeSequence, passage) {
			continue
		}

		newSequence := make([]int, len(activeSequence)+1)
		_ = copy(newSequence, activeSequence)
		newSequence[len(activeSequence)] = passage

		newNode := &bnb.Node{
			State: TravellingSalespersonProblemState{
				Sequence:    newSequence,
				CurrentCost: state.CurrentCost + distances.At(lastPassage, passage),
			},
		}
		if s.Bound(newNode) <= currentBound {
			nextNodes = append(nextNodes, newNode)
		}
	}

	return nextNodes
}

func TestSolverTSP(t *testing.T) {
	size := 4
	data := []float64{99, 1, 99, 99, 99, 99, 1, 99, 99, 99, 99, 1, 1, 99, 99, 99}
	distances := mat.NewDense(size, size, data)
	state := TravellingSalespersonProblemState{}

	tsp := &TravellingSalespersonProblem{
		DistanceMatrix: distances,
		State:          state,
		NSalesman:      distances.RawMatrix().Rows,
	}

	solver := bnb.Solver{tsp}
	solution, _, _, err := solver.Solve(bnb.SolverConfigs{Mode: bnb.DepthFirst})
	assert.NoError(t, err, "Solver should not raise error")

	bestState := &TravellingSalespersonProblemState{}
	solution.LoadState(bestState)

	assert.Equal(t, []int{0, 1, 2, 3}, bestState.Sequence, "should be 0/1/2/3 as best path")

	// no monotonous

	data = []float64{99, 99, 1, 99, 99, 99, 99, 1, 99, 1, 99, 99, 1, 99, 99, 99}
	distances = mat.NewDense(size, size, data)
	state = TravellingSalespersonProblemState{}

	tsp = &TravellingSalespersonProblem{
		DistanceMatrix: distances,
		State:          state,
		NSalesman:      distances.RawMatrix().Rows,
	}

	solver = bnb.Solver{tsp}
	solution, _, _, _ = solver.Solve(bnb.SolverConfigs{Mode: bnb.DepthFirst})
	solution.LoadState(bestState)
	assert.Equal(t, []int{0, 2, 1, 3}, bestState.Sequence, "should be 0/1/2/3 as best path")
}
