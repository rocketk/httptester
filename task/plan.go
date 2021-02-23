package task

import (
	"log"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

type Plan struct {
	TaskDef    TaskDef
	Assertions []Assertion
	// wg             *sync.WaitGroup
	// summaryChannel chan Summary
	// workerStopChannel chan int
	// barChannel chan int
	listener Listener
	report   Report
}

func (p *Plan) Start() {
	// log.Println(p.TaskDef.TimeUnit)
	// return
	listener := BuildSimpleListener(p.TaskDef.Concurrency*p.TaskDef.Loop, p.TaskDef.TimeUnit)
	p.listener = &listener
	// p.Assertions = []Assertion{
	// 	&StatusCodeAssertion{
	// 		ExpectedCodes: []int{200},
	// 	},
	// }
	summaryChannel := make(chan Summary, 128)
	barChannel := make(chan int, 8)
	// p.workerStopChannel = make(chan int, 1024)
	// p.barChannel = make(chan int, 1024)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go p.startListener(wg, summaryChannel, barChannel)
	for i := 0; i < p.TaskDef.Concurrency; i++ {
		wg.Add(1)
		w := &Worker{
			ID:      i,
			TaskDef: p.TaskDef,
			// WG:             p.wg,
			// SummaryChannel: summaryChannel,
			// WorkerStopChannel: p.workerStopChannel,
			Assertions: p.Assertions,
		}
		go w.StartLoop(wg, summaryChannel)
	}
	// log.Println("all workers started")
	// go log.Fatal(http.ListenAndServe(":8001", nil))

	// log.Printf("all workers stopped")
	if !p.TaskDef.DisableBar {
		wg.Add(1)
		go p.startBar(wg, barChannel)
	}
	wg.Wait()
}

func (p *Plan) startListener(wg *sync.WaitGroup, summaryChannel chan Summary, barChannel chan int) {
	defer wg.Done()
	// t0 := time.Now()
	var readChannelDuration int64
	p.listener.OnStart()
	count := p.TaskDef.Concurrency * p.TaskDef.Loop
	finished := 0
	// log.Printf("pre select: %d ns\n", time.Now().Sub(t0).Nanoseconds())
	for finished < count {
		// select {
		// case workerID := <-p.workerStopChannel:
		// 	// log.Printf("worker stopped, id: %d", workerID)
		// 	stoppedWorker++
		// 	if p.listener != nil {
		// 		p.listener.OnWorkerFinished(workerID)
		// 	} else {
		// 		log.Println("listener is nil")
		// 	}
		t1 := time.Now()
		summ := <-summaryChannel
		readChannelDuration += time.Now().Sub(t1).Milliseconds()

		finished++
		if p.listener != nil {
			if !p.TaskDef.DisableBar {
				barChannel <- 1
			}
			// log.Printf("OnRequestFinished: %+v\n", summ)
			p.listener.OnRequestFinished(summ)
			// fmt.Println("11")
		} else {
			log.Println("listener is nil")
		}
		// }
	}
	// natureDuration := time.Now().Sub(t0).Milliseconds()
	p.report = p.listener.OnPlanFinished()
	time.Sleep(101 * time.Millisecond)
	// log.Printf("natureDuration: %d ms, readChannelDuration: %d ms, %f\n", natureDuration, readChannelDuration, float64(readChannelDuration)/float64(natureDuration))
	p.TaskDef.PrintToStdOut()
	// fmt.Printf("%+v", p.report)
	p.report.PrintToStdOut()
}

func (p *Plan) startBar(wg *sync.WaitGroup, barChannel chan int) {
	defer wg.Done()
	if p.TaskDef.DisableBar {
		return
	}
	count := p.TaskDef.Loop * p.TaskDef.Concurrency
	// create and start new bar
	bar := pb.New(count).SetMaxWidth(100)
	// bar.AlwaysUpdate = true
	bar.SetRefreshRate(100 * time.Millisecond)
	bar.TimeBoxWidth = 0
	bar.SetWidth(100)
	bar.Start()

	for i := 0; i < count; i++ {
		select {
		case step := <-barChannel:
			bar.Add(step)
		}
	}

	bar.Finish()

}
