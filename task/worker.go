package task

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type TaskDef struct {
	Loop        int
	Concurrency int
	Timeout     time.Duration
	KeepAlive   bool
	URL         string
	Method      string
	Headers     []string
	Body        string
	TimeUnit    string
	DisableBar  bool
	PrintError  bool
	// AssertStatusCodes    []int
	// AssertJSONExpression string
}

const (
	NanoSecond         string = "ns"
	MicroSecond        string = "mms"
	MilliSecond        string = "ms"
	Second             string = "s"
	NanoSecondDivisor  int64  = 1
	MicroSecondDivisor int64  = 1000 * NanoSecondDivisor
	MilliSecondDivisor int64  = 1000 * MicroSecondDivisor
	SecondDivisor      int64  = 1000 * MilliSecondDivisor
)

// Summary for single http request
type Summary struct {
	// StartTimeOfAll  time.Time
	StartTime time.Time
	EndTime   time.Time
	// EndTimeOfAll    time.Time
	StatusCode      int
	Success         bool
	FailedAssertion string
	FailedCause     string
	HasError        bool
}

type Worker struct {
	ID      int
	TaskDef TaskDef
	// WG      *sync.WaitGroup
	// SummaryChannel chan Summary
	// WorkerStopChannel chan int
	Assertions []Assertion
	// reusedTransport http.RoundTripper
	httpClient *http.Client
}

func (w *Worker) StartLoop(wg *sync.WaitGroup, summaryChannel chan Summary) {
	// t0 := time.Now()
	defer wg.Done()
	w.initClient()
	// var costOfPreSending, costOfSending, costOfPostSending, costOfWritingChannel int64
	for i := 0; i < w.TaskDef.Loop; i++ {
		w.doRequest(wg, summaryChannel)
		// costOfPreSending += c1
		// costOfSending += c2
		// costOfPostSending += c3
		// costOfWritingChannel += c4
	}
	// log.Printf("StartLoop cost %d ms, costOfPreSending: %d ms, costOfSending: %d ms, costOfPostSending: %d ms, costOfWritingChannel: %d ms\n",
	// 	time.Now().Sub(t0).Milliseconds(), costOfPreSending/1000000, costOfSending/1000000, costOfPostSending/1000000, costOfWritingChannel/1000000)
	// w.WorkerStopChannel <- w.ID
}

func (w *Worker) initClient() {
	var dialKeepAlive time.Duration
	if w.TaskDef.KeepAlive {
		dialKeepAlive = 120 * time.Second
	} else {
		dialKeepAlive = 1 * time.Nanosecond
	}
	reusedTransport := &http.Transport{
		// Proxy: ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: dialKeepAlive,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives:   !w.TaskDef.KeepAlive,
		ForceAttemptHTTP2:   false,
		MaxIdleConns:        w.TaskDef.Concurrency * 2,
		MaxIdleConnsPerHost: w.TaskDef.Concurrency,
		MaxConnsPerHost:     0,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 90 * time.Second,
		// ExpectContinueTimeout: 1 * time.Second,
	}
	var timeout time.Duration
	if w.TaskDef.Timeout == 0 {
		timeout = 10 * time.Second
	} else {
		timeout = w.TaskDef.Timeout
	}
	w.httpClient = &http.Client{Timeout: timeout, Transport: reusedTransport}
}

func (w *Worker) doRequest(wg *sync.WaitGroup, summaryChannel chan Summary) {
	summary := Summary{}
	req, err := http.NewRequest(w.TaskDef.Method, w.TaskDef.URL, bytes.NewBuffer([]byte(w.TaskDef.Body)))
	// req.Header.Add("Connection", "keep-alive")
	// log.Printf("%+v", w.TaskDef.Headers)
	if len(w.TaskDef.Headers) > 0 {
		// log.Printf("len of headers: %d, headers: %+v\n", len(w.TaskDef.Headers), w.TaskDef.Headers)
		for _, header := range w.TaskDef.Headers {
			if strings.Trim(header, " ") == "" {
				continue
			}
			pair := strings.Split(header, ":")
			key := strings.Trim(pair[0], " ")
			if key == "" {
				continue
			}
			var value string
			if len(pair) < 2 {
				value = ""
			} else {
				value = pair[1]
			}
			req.Header.Add(key, value)
		}
	}
	summary.StartTime = time.Now()
	resp, err := w.httpClient.Do(req)
	summary.EndTime = time.Now()
	if err != nil {
		// panic(err)
		if w.TaskDef.PrintError {
			log.Printf("error: %s\n", err)
		}
		summary.HasError = true
		summaryChannel <- summary
		return
	}
	wg.Add(1)
	go w.verifyAllAssertions(resp, wg, &summary, summaryChannel)
	return
}

// verifyAllAssertions returns success, assertionName, cause
func (w *Worker) verifyAllAssertions(resp *http.Response, wg *sync.WaitGroup, summary *Summary, summaryChannel chan Summary) {
	// log.Println("111")
	defer wg.Done()
	// log.Printf("%+v\n", resp)
	// t0 := time.Now()
	body, err := ioutil.ReadAll(resp.Body)
	// t1 := time.Now()
	resp.Body.Close()
	if err != nil {
		summary.Success = false
		summary.FailedAssertion = ""
		summary.FailedCause = err.Error()
		summaryChannel <- *summary
		return
	}
	httpResponse := HttpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       body,
	}
	if len(w.Assertions) == 0 {
		summary.Success = true
		summaryChannel <- *summary
		return
	}
	for _, a := range w.Assertions {
		if a == nil {
			continue
		}
		if ok, cause := a.Assert(httpResponse); !ok {
			summary.Success = false
			summary.FailedAssertion = a.Name()
			summary.FailedCause = cause
			summaryChannel <- *summary
			if w.TaskDef.PrintError {
				log.Printf("Assertion Failed, Caused by: %s, %s\n", summary.FailedAssertion, summary.FailedCause)
			}
			return
		}
	}
	// t2 := time.Now()
	// log.Printf("reading body: %d ms, asserting: %d ms", t1.Sub(t0).Milliseconds(), t2.Sub(t1).Milliseconds())
	summary.Success = true
	summaryChannel <- *summary
}

var neverReusedTransport http.RoundTripper = &http.Transport{
	// Proxy: ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 1 * time.Nanosecond,
		DualStack: true,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          1,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

// var reusedTransport http.RoundTripper = &http.Transport{
// 	// Proxy: ProxyFromEnvironment,
// 	DialContext: (&net.Dialer{
// 		Timeout:   60 * time.Second,
// 		KeepAlive: 120 * time.Second,

// 		DualStack: true,
// 	}).DialContext,
// 	DisableKeepAlives:   false,
// 	ForceAttemptHTTP2:   false,
// 	MaxIdleConns:        10000,
// 	MaxIdleConnsPerHost: 1000,
// 	MaxConnsPerHost:     0,
// 	IdleConnTimeout:     90 * time.Second,
// 	TLSHandshakeTimeout: 90 * time.Second,
// 	// ExpectContinueTimeout: 1 * time.Second,
// }

// func transport() http.RoundTripper {
// return &http.Transport{
// 	// Proxy: ProxyFromEnvironment,
// 	DialContext: (&net.Dialer{
// 		Timeout:   60 * time.Second,
// 		KeepAlive: 120 * time.Second,
// 		DualStack: true,
// 	}).DialContext,
// 	ForceAttemptHTTP2:   false,
// 	MaxIdleConns:        10000,
// 	IdleConnTimeout:     90 * time.Second,
// 	TLSHandshakeTimeout: 90 * time.Second,
// 	// ExpectContinueTimeout: 1 * time.Second,
// }
// 	return reusedTransport
// }

func (d TaskDef) PrintToStdOut() {
	fmt.Println("-- Configuration --")
	fmt.Printf("Concurrency: %d\t", d.Concurrency)
	fmt.Printf("Loop: %d\t", d.Loop)
	fmt.Printf("Timeout: %d ms\t", d.Timeout.Milliseconds())
	fmt.Printf("KeepAlive: %t\t", d.KeepAlive)
	fmt.Printf("TimeUnit: %s\t", d.TimeUnit)
	fmt.Printf("Method: %s\t", d.Method)
	fmt.Printf("URL: %s\n", d.URL)
	fmt.Printf("Headers: %s\n", d.Headers)
	fmt.Printf("Body: %s\n\n", d.Body)
}
