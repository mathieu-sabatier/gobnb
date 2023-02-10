package gobnb

import (
	"math"
	"time"
)

type convergenceChecker struct {
	maxDeltaBound   *float64
	maxIterCount    *int
	maxDuration     *int // duration in seconds
	iterCount       int
	startedAt       time.Time
	convergenceMode string
}

func newConvergenceCheckerFromConfig(config *SolverConfigs) *convergenceChecker {
	var maxDeltaBound *float64
	var maxIterCount, maxDuration *int
	if config.MaxDuration > 0 {
		maxDuration = &config.MaxDuration
	}
	if config.MaxIterCount > 0 {
		maxIterCount = &config.MaxIterCount
	}
	if config.MaxDeltaBound > 0 {
		maxDeltaBound = &config.MaxDeltaBound
	}
	return &convergenceChecker{
		iterCount:     1,
		startedAt:     time.Now(),
		maxDuration:   maxDuration,
		maxIterCount:  maxIterCount,
		maxDeltaBound: maxDeltaBound,
	}
}

func (c *convergenceChecker) iter(currentBound float64, currentObjective float64) bool {
	c.iterCount += 1

	// stop if iter count reached
	if c.maxIterCount != nil {
		if c.iterCount >= *c.maxIterCount {
			c.convergenceMode = "node_count"
			return true
		}
	}

	// stop if objective and bound are close enouth
	if c.maxDeltaBound != nil {
		if math.Abs(currentBound-currentObjective) <= *c.maxDeltaBound {
			c.convergenceMode = "delta_bound_objective"
			return true
		}
	}

	return false
}
