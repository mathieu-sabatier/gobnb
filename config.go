package gobnb

import "github.com/emirpasic/gods/utils"

type SolverMode uint

const (
	DepthFirst SolverMode = iota
	BreadthFirst
	BestBound
	Custom
)

var SolverModeNames = [...]string{"DepthFirst", "BreathFirst", "BestBound", "Custom"}

func (mode SolverMode) String() string {
	return SolverModeNames[mode]
}

type SolverConfigs struct {
	MaxDeltaBound    float64
	MaxDuration      int
	MaxIterCount     int
	Mode             SolverMode
	customComparator utils.Comparator
}

func newComparatorFromConfig(config *SolverConfigs) utils.Comparator {
	if config.Mode == DepthFirst {
		return func(a, b interface{}) int {
			priorityA := a.(*Node).depth
			priorityB := b.(*Node).depth
			return -utils.IntComparator(priorityA, priorityB)
		}
	} else if config.Mode == BreadthFirst {
		return func(a, b interface{}) int {
			priorityA := a.(*Node).depth
			priorityB := b.(*Node).depth
			return +utils.IntComparator(priorityA, priorityB)
		}
	} else if config.Mode == BestBound {
		return func(a, b interface{}) int {
			priorityA := a.(*Node).depth
			priorityB := b.(*Node).depth
			return -utils.IntComparator(priorityA, priorityB)
		}
	} else if config.Mode == Custom {
		return config.customComparator
	}
	panic(config.Mode)
}
