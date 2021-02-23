package task

import (
	"log"
	"math"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceSorting(t *testing.T) {
	ast := assert.New(t)
	a := make([]int, 3)
	a[0] = 5
	a[1] = 8
	a[2] = 1
	sort.Ints(a)
	b := make([]int, 3)
	b[0] = 1
	b[1] = 5
	b[2] = 8
	ast.Equal(a, b)
}

func TestSliceSortingNotEqualToArray(t *testing.T) {
	ast := assert.New(t)
	a := []int{5, 8, 1}
	sort.Ints(a)
	b := [...]int{1, 5, 8}
	ast.NotEqual(a, b)
}

func TestSliceAppend(t *testing.T) {
	ast := assert.New(t)
	a := make([]int, 3)
	b := append(a, 1)
	ast.Equal(a, b)
}

func TestInt64Devided(t *testing.T) {
	ast := assert.New(t)
	var i int64
	i = 5
	d := i / 2
	log.Printf("%v, %T", d, d)
	ast.IsType(int64(2), d)
	ast.Equal(int64(2), d)
	// e := 5 / 2
	// log.Printf("%v, %T", e, e)
}

func TestCalculate(t *testing.T) {
	ast := assert.New(t)
	l := BuildSimpleListener(10, "ms")
	l.successCount = 10
	l.costs = []int64{5, 10, 4, 3, 9, 8, 1, 23, 12, 2}
	l.calculate()
	ast.Equal(float64(7.7), l.mean)
	ast.Equal(int64(77), l.totalCost)
	ast.Truef(math.Abs(l.stdDev-float64(6.16522505671934)) < 0.0001, "expect: %f, actual: %f", 6.16522505671934, l.stdDev)
	ast.Equal(float64(6.5), l.median)
	ast.Equal(int64(23), l.max)
	ast.Equal(int64(1), l.min)
}

func TestCalculate2(t *testing.T) {
	ast := assert.New(t)
	l := BuildSimpleListener(10, "ms")
	l.successCount = 10
	l.costs = []int64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	l.calculate()
	ast.Equal(float64(2), l.mean)
	ast.Equal(int64(20), l.totalCost)
	ast.Equal(float64(0), l.stdDev)
	ast.Equal(float64(2), l.median)
}
