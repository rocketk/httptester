package task

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"
)

type Listener interface {
	OnStart()
	OnRequestFinished(summary Summary)
	// OnWorkerFinished(workerID int)
	OnPlanFinished() Report
}
type Report interface {
	PrintToStdOut()
}

type SimpleListener struct {
	costs []int64
	// costsOfPreSending      []int64
	// costsOfPostSending     []int64
	index        int
	successCount int
	failedCount  int
	errorCount   int
	totalCost    int64
	// totalCostOfPreSending  int64
	// totalCostOfPostSending int64
	natureDuration  time.Duration
	start           time.Time
	end             time.Time
	mean            float64
	max             int64
	min             int64
	median          float64
	stdDev          float64
	calculated      bool
	timeunit        string
	timeunitDivisor int64
	throughput      int64
}

func BuildSimpleListener(capacity int, timeunit string) SimpleListener {
	var d int64
	if timeunit == NanoSecond {
		d = 1
	} else if timeunit == MicroSecond {
		d = MicroSecondDivisor
	} else if timeunit == MilliSecond {
		d = MilliSecondDivisor
	} else if timeunit == Second {
		d = SecondDivisor
	} else {
		// log.Panicf("unknown timeunit %s\n", timeunit)
		timeunit = "ms"
		d = MilliSecondDivisor
	}
	return SimpleListener{
		costs: make([]int64, capacity),
		// costsOfPreSending:  make([]int64, capacity),
		// costsOfPostSending: make([]int64, capacity),
		timeunit:        timeunit,
		timeunitDivisor: d,
	}
}
func (s *SimpleListener) OnStart() {
	s.start = time.Now()
}
func (s *SimpleListener) OnRequestFinished(summary Summary) {
	if summary.HasError {
		s.errorCount++
	} else if summary.Success {
		s.successCount++
	} else {
		s.failedCount++
	}
	// costPreSending := summary.StartTime.Sub(summary.StartTimeOfAll).Nanoseconds()
	cost := summary.EndTime.Sub(summary.StartTime).Nanoseconds()
	// costPostSending := summary.EndTimeOfAll.Sub(summary.EndTime).Nanoseconds()
	// log.Printf("costPreSending: %d, costPreSending: %d, cost: %d, costPostSending: %d, timeunitDivisor: %d\n", costPreSending, costPreSending/s.timeunitDivisor, cost/s.timeunitDivisor, costPostSending/s.timeunitDivisor, s.timeunitDivisor)
	// s.total += cost
	// append(s.costsOfPreSending, costPreSending/s.timeunitDivisor, s.index)
	appendToSlice(s.costs, cost/s.timeunitDivisor, s.index)
	// append(s.costsOfPostSending, costPostSending/s.timeunitDivisor, s.index)
	s.index++
}

func appendToSlice(arr []int64, item int64, index int) {
	if index >= len(arr) {
		log.Printf("out of slice capacity, capacity: %d", len(arr))
		return
	}
	arr[index] = item
}

// func (s *SimpleListener) OnWorkerFinished(workerID int) {
// }
func (s *SimpleListener) OnPlanFinished() Report {
	if s.calculated {
		return s
	}
	s.end = time.Now()
	s.natureDuration = s.end.Sub(s.start)
	s.calculate()
	return s
}

func (s *SimpleListener) PrintToStdOut() {
	fmt.Println("-- Conclusion --")
	fmt.Printf("total count: %d\n", s.successCount+s.failedCount+s.errorCount)
	fmt.Printf("success count: %d\n", s.successCount)
	fmt.Printf("failed count: %d\n", s.failedCount)
	fmt.Printf("error count: %d\n", s.errorCount)
	fmt.Printf("nature duration: %d %s\n", s.natureDuration/time.Duration(s.timeunitDivisor), s.timeunit)
	fmt.Printf("total cost: %d %s\n", s.totalCost, s.timeunit)
	fmt.Printf("max: %d %s\n", s.max, s.timeunit)
	fmt.Printf("min: %d %s\n", s.min, s.timeunit)
	fmt.Printf("median: %d %s\n", int64(s.median), s.timeunit)
	fmt.Printf("mean: %d %s\n", int64(s.mean), s.timeunit)
	fmt.Printf("standard deviation: %f\n", s.stdDev)
	fmt.Printf("throughput: %d requests/second\n", s.throughput)
	// fmt.Printf("len: %d, costs: %+v\n", len(s.costs), s.costs)
}

func (s *SimpleListener) calculate() {
	if s.calculated {
		return
	}

	totalCount := s.successCount + s.failedCount + s.errorCount
	s.totalCost = 0
	// s.totalCostOfPreSending = 0
	// s.totalCostOfPostSending = 0
	// fmt.Printf("s.costs: %+v\n", s.costs)
	// fmt.Printf("s.costsOfPreSending: %+v\n", s.costsOfPreSending)

	for _, cost := range s.costs {
		s.totalCost += cost
	}
	// for _, cost := range s.costsOfPreSending {
	// 	s.totalCostOfPreSending += cost
	// }
	// for _, cost := range s.costsOfPostSending {
	// 	s.totalCostOfPostSending += cost
	// }
	s.mean = float64(s.totalCost) / float64(totalCount)

	sort.Slice(s.costs, func(i, j int) bool { return s.costs[i] < s.costs[j] })
	s.min = s.costs[0]

	l := len(s.costs)
	s.max = s.costs[l-1]

	if l > 1 {
		if l%2 == 0 {
			// even
			left := s.costs[l/2-1]
			right := s.costs[l/2]
			s.median = float64(left+right) / 2
		} else {
			// odd
			s.median = float64(s.costs[(l-1)/2])
		}
	} else if l == 1 {
		s.median = float64(s.costs[0])
	} else {
		s.median = 0
	}

	// log.Printf("s.successCount: %d ms, s.costDuration: %d ms\n", s.successCount, s.costDuration.Milliseconds())
	s.throughput = int64(float64(s.successCount*1000*1000*1000) / float64(s.natureDuration.Nanoseconds()))

	var sd float64
	for _, cost := range s.costs {
		sd += math.Pow(math.Abs(float64(cost)-s.mean), 2)
	}
	sd = math.Sqrt(sd / float64(l))
	s.stdDev = sd
	s.calculated = true
}
