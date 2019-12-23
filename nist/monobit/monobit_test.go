package monobit

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestMonobit(t *testing.T) {
	e := []bool{true, false, true, true, false, true, false, true, false, true}
	s := computeTestStatistic(e)
	expectedS := .632455532
	log.Printf("S Value: %f", s)
	assert.InEpsilon(t, expectedS, s, 0.00001)

	p := computePValue(s)
	expectedP := 0.527089
	log.Printf("P Value: %f", p)
	assert.InEpsilon(t, expectedP, p, 0.00001)

	mb := NewMonoBitTest()
	result, err := mb.IsRandom(e)
	assert.True(t, result)
	assert.NoError(t, err)
}

func TestMonobitFull(t *testing.T) {
	seq := "1100100100001111110110101010001000100001011010001100001000110100110001001100011001100010100010111000"
	e := make([]bool, 0)
	for _, val := range seq {
		e = append(e, val == rune('1'))
	}
	s := computeTestStatistic(e)
	expectedS := 1.6
	log.Printf("S Value: %f", s)
	assert.InEpsilon(t, expectedS, s, 0.00001)

	p := computePValue(s)
	expectedP := 0.109599
	log.Printf("P Value: %f", p)
	assert.InEpsilon(t, expectedP, p, 0.00001)

	mb := NewMonoBitTest()
	result, err := mb.IsRandom(e)
	assert.True(t, result)
	assert.NoError(t, err)
}
