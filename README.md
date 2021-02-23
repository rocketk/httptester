### httptester
A lightweight tool for pressure testing

### Get Started
Download the latest binary file [Releases](https://github.com/rocketk/httptester/releases)

Save the binary file `httptester` into any directory, and open the terminal

For help
```bash
./httptester run -h
```

#### First Example
```bash
./httptester run -u 'https://www.baidu.com/'
```

Output:
```text
 1 / 1 [================================================================================] 100.00% 0s
-- Configuration --
Concurrency: 1  Loop: 1 Timeout: 10000000000    KeepAlive: true TimeUnit: ms    Method: GET     URL: https://www.baidu.com/
Headers: []
Body: 

-- Conclusion --
total count: 1
success count: 1
failed count: 0
error count: 0
nature duration: 207 ms
total cost: 207 ms
max: 207 ms
min: 207 ms
median: 207 ms
mean: 207 ms
standard deviation: 0.000000
throughput: 4 requests/second
```

#### More Examples
```bash
./httptester run -u 'https://www.baidu.com/' -c 10 -l 10
```
This example will start 10 goroutines to send the http request, and each goroutine will repeat sending the same request for 10 times. 


