## gobnb

A branch-and-bound starter for Go.

This software is copyright (c) by Mathieu Sabatier (mathieu.sbt@gmail.com).

This software is released under the MIT software license.
This license, including disclaimer, is available in the 'LICENSE' file.

### Quick Start

A problem implements the *Problem* interface. 

```go
type Problem interface {
    Sense() ProblemSense
    Objective(*Node) float64
    Bound(*Node) float64
    Branch(*Node, float64) []*Node
    LoadInitialNode() *Node
}
```

A node contains a state, that you can define and later use within : 

```go

type Node struct {
    State  any
}

type TSPState struct {
    Sequence    []int
    CurrentCost float64
}

```

### Solve

```go
data = []float64{99, 99, 1, 99, 99, 99, 99, 1, 99, 1, 99, 99, 1, 99, 99, 99}
distances = mat.NewDense(size, size, data)
state = TravellingSalespersonProblemState{}

tsp = &TravellingSalespersonProblem{
    DistanceMatrix: distances,
    State:          state,
    NSalesman:      distances.RawMatrix 
    Rows,
}

solver = Solver{tsp}
solution, _, _, _ = solver.Solve(&SolverConfigs{Mode: DepthFirst, MaxIterCount:100})

```

## TO DO

- [ ] Badge & Coverage
- [X] Statistics outputs
- [ ] MPI implementation
- [ ] Benchmark for TSP