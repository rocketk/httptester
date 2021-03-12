/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"strconv"
	"strings"
	"time"

	"github.com/rocketk/httptester/task"
	"github.com/spf13/cobra"
)

var (
	loop                  int
	concurrency           int
	timeout               time.Duration
	keepAlive             bool
	url                   string
	method                string
	headers               []string
	body                  string
	assertStatusCodes     string
	assertJSONExpression  string
	assertRegexExpression string
	expectedCode          string
	timeunit              string
	disableBar            bool
	printError            bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the test case",
	Long: `Run the test case. For example:

httptester run --loop 10 --concurrency 10 --timeout 10s
httptester run --loop 10 --concurrency 100 --timeout 500ms --keep-alive false 
`,
	Run: func(cmd *cobra.Command, args []string) {
		if url == "" {
			panic("url is required")
		}
		// fmt.Printf("keepAlive: %t\n", keepAlive)
		assertions := make([]task.Assertion, 0, 8)
		if len(assertStatusCodes) > 0 {
			intAssertStatusCodes := make([]int, 0, 8)
			for _, code := range strings.Split(assertStatusCodes, " ") {
				codeInt, err := strconv.Atoi(code)
				if err != nil {
					panic(err)
				}
				intAssertStatusCodes = append(intAssertStatusCodes, codeInt)
			}
			assertions = append(assertions, &task.StatusCodeAssertion{
				ExpectedCodes: intAssertStatusCodes,
			})
		}
		if len(assertJSONExpression) > 0 {
			assertions = append(assertions, &task.JsonPathAssertion{
				Expression: assertJSONExpression,
			})
		}
		if len(assertRegexExpression) > 0 {
			assertions = append(assertions, &task.RegexAssertion{
				Expression: assertRegexExpression,
			})
		}

		taskDef := task.TaskDef{
			Loop:        loop,
			Concurrency: concurrency,
			Timeout:     timeout,
			KeepAlive:   keepAlive,
			URL:         url,
			Method:      method,
			Headers:     headers,
			Body:        body,
			TimeUnit:    timeunit,
			DisableBar:  disableBar,
			PrintError:  printError,
			// AssertStatusCodes:    intAssertStatusCodes,
			// AssertJSONExpression: assertJSONExpression,
		}
		plan := &task.Plan{
			TaskDef:    taskDef,
			Assertions: assertions,
		}
		// data, _ := json.MarshalIndent(plan, "", "  ")
		// fmt.Printf("%s\n", data)
		plan.Start()

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().BoolVarP(&keepAlive, "keep-alive", "", true, "to use keep-alive for connections")
	runCmd.Flags().BoolVarP(&disableBar, "disable-bar", "", false, "disable the progress bar")
	runCmd.Flags().IntVarP(&loop, "loop", "l", 1, "how many requests would a goroutine send synchronously")
	runCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "how many goroutines would run concurrently")
	runCmd.Flags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "how many goroutines would run concurrently")
	runCmd.Flags().StringVarP(&url, "url", "u", "", "the target url you want to test")
	runCmd.Flags().StringVarP(&body, "body", "b", "", "the request body")
	runCmd.Flags().StringArrayVarP(&headers, "header", "H", []string{}, "the headers")
	runCmd.Flags().StringVarP(&assertStatusCodes, "assert-status-codes", "", "", "assertion: expected http response status codes, use space-splited string")
	runCmd.Flags().StringVarP(&assertJSONExpression, "assert-json-expression", "", "", "assertion: use jsonpath expression to verify a field, e.g. '$.expensive == 10', which '$' means the root of the json body. see https://github.com/oliveagle/jsonpath for more details")
	runCmd.Flags().StringVarP(&assertRegexExpression, "assert-regex-expression", "", "", "assertion: use regex expression to validate the response body, e.g. '$.expensive == 10'")
	runCmd.Flags().StringVarP(&timeunit, "time-unit", "", "ms", "time unit for printing report and calculating the standard deviation. 'ms' for milli-second, 'mms' for micro-second, 'ns' for nano-second, 's' for second")
	runCmd.Flags().StringVarP(&method, "method", "", "GET", "http method")
	runCmd.Flags().BoolVarP(&printError, "print-error", "e", false, "to print the error information")
}
