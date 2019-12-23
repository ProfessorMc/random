package monobit

import "math"

const pThreshold = 0.01

type monobitTest struct {
}

type ArgError struct {
}

func (e *ArgError) Error() string {
	return "invalid argument"
}

func NewMonoBitTest() *monobitTest {
	return &monobitTest{}
}

func (m *monobitTest) IsRandom(seq []bool) (bool, error) {
	n := len(seq)
	if n == 0 {
		return false, &ArgError{}
	}
	s := computeTestStatistic(seq)
	p := computePValue(s)
	return p >= pThreshold, nil
}

func computeTestStatistic(seq []bool) float64 {
	s := 0
	n := float64(len(seq))
	for _, val := range seq {
		if val {
			s = s + 1
		} else {
			s = s - 1
		}
	}
	return math.Abs(float64(s)) / math.Sqrt(n)
}

func computePValue(s float64) float64 {
	return math.Erfc(s / math.Sqrt(2))
}
