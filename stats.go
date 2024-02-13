package gobnb

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"
	"time"
)

type statsWriter struct {
	iterCount  int
	printEvery *int
	writer     *tabwriter.Writer
	startedAt  time.Time
}

func newStatsWriterFromConfig(config *SolverConfigs) *statsWriter {
	printEvery := 20000
	if config.PrintStatsEvery > 0 {
		printEvery = config.PrintStatsEvery
	}
	w := tabwriter.NewWriter(os.Stdout, 15, 8, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Node explored\tNode count\tBest objective\tBest bound\tTotal time\tNode /s\t")
	return &statsWriter{
		iterCount:  1,
		printEvery: &printEvery,
		writer:     w,
		startedAt:  time.Now(),
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (w *statsWriter) iter(nodeCount int, bestObjective float64, bestBound float64) {
	w.iterCount += 1

	if (w.iterCount % *w.printEvery) == 0 {
		now := time.Now()
		elapsedTime := now.Sub(w.startedAt).Seconds()
		fmt.Fprintln(w.writer, w.iterCount, "\t", nodeCount, "\t", roundFloat(bestObjective, 4), "\t", roundFloat(bestBound, 4), "\t", roundFloat(elapsedTime, 2), "\t", roundFloat(float64(w.iterCount)/elapsedTime, 2), "\t")
	}
	w.writer.Flush()
}

func (w *statsWriter) inform(nodeCount int, bestObjective float64, bestBound float64) {
	now := time.Now()
	elapsedTime := now.Sub(w.startedAt).Seconds()
	fmt.Fprintln(w.writer, "*", w.iterCount, "\t", nodeCount, "\t", bestObjective, "\t", bestBound, "\t", roundFloat(elapsedTime, 2), "\t", roundFloat(float64(w.iterCount)/elapsedTime, 2), "\t")
	err := w.writer.Flush()
	if err != nil {
		fmt.Println("Error", err)
	}
}

func (w *statsWriter) terminate() {
	w.writer.Flush()
}
