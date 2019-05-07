package random

import (
	"math"
	"sync"
	"testing"
)

func TestNewQuantizer(t *testing.T) {
	q := NewQuantizer()
	q.SetIntervals(0)
	q.CreateIntervals()
	t.Log(len(q.IntervalMap))
	q.EvaluateMatches(4)
}

func TestRaw(t *testing.T) {
	minx := 0
	miny := 0
	currX := 0
	currY := 0
	dimension := 255
	runs := 10
	length := 4

	matches := uint64(0)
	matchChan := make(chan uint64)

	go func() {
		for{
			select {
			case match := <- matchChan:
				matches += match
			}
		}
	}()

	var wg sync.WaitGroup
	for currX = minx; currX < runs; currX ++ {
		for currY = miny; currY < runs; currY ++ {
			wg.Add(1)
			go testPoint(currX,currY,dimension,length, matchChan, &wg)
		}
	}

	wg.Wait()
	close(matchChan)
	t.Logf("Received %d: matches", matches)
}

func testPoint(xval int, yval int, dimension int, length int, matchchan chan uint64, wg *sync.WaitGroup){
	defer wg.Done()
	for currX := 0; currX < dimension; currX ++ {
		for currY := 0; currY < dimension; currY ++ {
			dist := math.Pow(float64(xval-currX), 2) + math.Pow(float64(yval-currY), 2)
			lengthSq := math.Pow(float64(length), 2)
			if dist <= lengthSq {
				matchchan <- 1
			}
		}
	}
}