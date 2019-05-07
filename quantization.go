package random

import (
	"fmt"
	"log"
	"math"
	"sync"
)

type Quantizer struct {
	NumInterval uint64
	IntervalMap map[uint64]uint64
	Sample [math.MaxUint8][math.MaxUint8]bool
}

func NewQuantizer() *Quantizer {
	var sample [math.MaxUint8][math.MaxUint8]bool

	for i := 0; i < math.MaxUint8; i++{
		for j:= 0; j < math.MaxUint8; j++ {
			sample[i][j] = true
		}
	}
	return &Quantizer{
		IntervalMap:make(map[uint64]uint64),
		Sample:sample,
	}
}

func (q *Quantizer) SetIntervals(NumInterval uint64) {
	q.NumInterval = NumInterval
}

func (q *Quantizer) CreateIntervals() {
	var intervalType uint8
	interval := GetInterval(intervalType, q.NumInterval)
	for i := uint64(0); i < math.MaxUint8; i += interval {
		q.IntervalMap[i + interval] = 0
		//log.Printf("Initialized Interval: (%d - %d] ", i, i + interval)
	}
	log.Printf("Initialized Map, Length: %d", len(q.IntervalMap))
}

func (q *Quantizer) EvaluateMatches(length uint32) {
	//Create Array
	dimension := length * 2 + 1
	fmt.Printf("Dimension: %d\n", dimension)
	eval := make([][]bool, dimension)
	for column := range eval {
		eval[column] = make([]bool, dimension)
	}
	for i := uint32(0); i < dimension; i++{
		for j:= uint32(0); j < dimension; j++ {
			x := float64(i)
			y := float64(j)
			l := float64(length)
			lengthFromCenter := math.Pow(x-l, 2)+math.Pow(y-l, 2)
			if lengthFromCenter <= math.Pow(l, 2) {
				eval[i][j] = true
				//fmt.Println(i-length)
			}
			//fmt.Printf("%v\t", eval[i][j])
			//fmt.Printf("Length of (%d, %d): %v, \n", i, j, lengthFromCenter)
		}
		//fmt.Printf("\n")
	}
	matches := uint64(0)
	var wg sync.WaitGroup


	matchesChan := make(chan uint64)
	doneChan := make(chan uint64)

	go func() {
		defer close(doneChan)
		wg.Wait()
	}()

	for xOffset := length; xOffset > 0; xOffset -- {
		for yOffset := length; yOffset > 0; yOffset -- {
			wg.Add(1)
			go q.evaluate(eval, xOffset, yOffset, matchesChan, &wg)

		}
	}
	//go q.evaluate(eval, length, length, matchesChan, &wg)


	listener: for {
		select {
		case <- doneChan:
			break listener
		case match := <- matchesChan:
			matches += match

		}
	}

	//offset := length
	//for i := uint32(0); i + offset < dimension  && i <= uint32(len(q.Sample)); i++ {
	//	for j := uint32(0); j + offset < dimension  && j <= uint32(len(q.Sample)); j++ {
	//		//fmt.Printf("Eval: [%d,%d]", i, j)
	//		//fmt.Printf("I, J, O: [%d,%d, %d]", i, j, offset)
	//		if eval[i+offset][j+offset] && q.Sample[i][j] {
	//			matches++
	//		}
	//	}
	//}
	fmt.Printf("Matched: %d\n", matches)

}

func (q *Quantizer) evaluate(eval [][]bool, xOffset uint32, yOffset uint32, matches chan uint64,wg *sync.WaitGroup){
	dimension := uint32(len(eval))

	defer wg.Done()
	matchTotal := 0
	for i := uint32(0); i + xOffset < dimension  && i <= uint32(len(q.Sample)); i++ {
		for j := uint32(0); j + yOffset < dimension  && j <= uint32(len(q.Sample)); j++ {
			//fmt.Printf("Eval: [%d,%d]", i, j)
			//fmt.Printf("I, J, O: [%d,%d, %d]", i, j, offset)
			if eval[i+xOffset][j+yOffset] && q.Sample[i][j] {
				matches <- 1
				matchTotal ++
			}
		}
	}


}

func GetInterval(t interface{}, interval uint64) uint64{
	switch tType := t.(type) {
	case uint8:
		if interval > math.MaxUint8{
			log.Fatalf("Received interval too large for type: %v", tType)
		}
		if interval == 0 {
			interval = math.MaxUint8
		}
		return math.MaxUint8 / interval
	case uint16:
		if interval > math.MaxUint16{
			log.Fatalf("Received interval too large for type: %v", tType)
		}
		if interval == 0 {
			interval = math.MaxUint16
		}
		return math.MaxUint16 / interval
	case uint32:
		if interval > math.MaxUint32{
			log.Fatalf("Received interval too large for type: %v", tType)
		}
		if interval == 0 {
			interval = math.MaxUint32
		}
		return math.MaxUint32 / interval
	case uint64:
		if interval > math.MaxUint64{
			log.Fatalf("Received interval too large for type: %v", tType)
		}
		if interval == 0 {
			interval = math.MaxUint64
		}
		return math.MaxUint64 / interval
	default:
		log.Fatalf("Only uint types are supported, requested: %v", tType)
		return 0
	}
}