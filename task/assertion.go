package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/oliveagle/jsonpath"
)

type Assertion interface {
	Assert(resp HttpResponse) (bool, string)
	Name() string
	Validate() error
}

type HttpResponse struct {
	Status     string // e.g. "200 OK"
	StatusCode int    // e.g. 200
	Header     http.Header
	Body       []byte
}

type StatusCodeAssertion struct {
	ExpectedCodes []int
}

func (a StatusCodeAssertion) Assert(resp HttpResponse) (bool, string) {
	if len(a.ExpectedCodes) == 0 {
		return true, ""
	}
	for _, code := range a.ExpectedCodes {
		// log.Printf("%d", code)
		// log.Printf("%+v\n", resp)
		if code == resp.StatusCode {
			return true, ""
		}
	}
	return false, "Invalid Status Code: " + strconv.Itoa(resp.StatusCode)
}

func (a StatusCodeAssertion) Name() string {
	return "StatusCodeAssertion"
}

func (a StatusCodeAssertion) Validate() error {
	if len(a.ExpectedCodes) == 0 {
		return errors.New("At least 1 item is required for ExpectedCodes")
	}
	return nil
}

type JsonPathAssertion struct {
	Expression  string
	initialized bool
	query       string
	verb        string
	value       string
	handler     verbHandler
}

func (a JsonPathAssertion) Assert(resp HttpResponse) (bool, string) {
	if !a.initialized {
		if err := a.Validate(); err != nil {
			return false, err.Error()
		}
	}
	var jsonData interface{}
	json.Unmarshal(resp.Body, &jsonData)
	arg, err := jsonpath.JsonPathLookup(jsonData, a.query)
	if err != nil {
		return false, err.Error()
	}
	return a.handler(fmt.Sprintf("%v", arg), a.value)
}

func (a JsonPathAssertion) Name() string {
	return "JsonPathAssertion"
}

type verbHandler func(arg string, target string) (bool, string)

var verbFuncMap = map[string]verbHandler{
	">":  gt,
	">=": ge,
	"<":  lt,
	"<=": le,
	"==": eq,
}

func (a *JsonPathAssertion) Validate() error {
	if len(a.Expression) == 0 {
		return errors.New("Expression is required")
	}
	args := strings.Split(a.Expression, " ")
	if len(args) < 3 {
		return errors.New("Invalid Expression")
	}
	a.query = args[0]
	a.verb = args[1]
	a.value = args[2]
	handler, ok := verbFuncMap[a.verb]
	if !ok {
		return fmt.Errorf("No handler was found for the verb '%s'", a.verb)
	}
	a.handler = handler
	a.initialized = true
	return nil
}

func gt(arg string, target string) (bool, string) {
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return false, err.Error()
	}
	targetInt, err := strconv.Atoi(target)
	if err != nil {
		return false, err.Error()
	}
	if argInt > targetInt {
		return true, ""
	}
	return false, fmt.Sprintf("Assertion failed: %d > %d", argInt, targetInt)
}
func ge(arg string, target string) (bool, string) {
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return false, err.Error()
	}
	targetInt, err := strconv.Atoi(target)
	if err != nil {
		return false, err.Error()
	}
	if argInt >= targetInt {
		return true, ""
	}
	return false, fmt.Sprintf("Assertion failed: %d >= %d", argInt, targetInt)
}
func lt(arg string, target string) (bool, string) {
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return false, err.Error()
	}
	targetInt, err := strconv.Atoi(target)
	if err != nil {
		return false, err.Error()
	}
	if argInt < targetInt {
		return true, ""
	}
	return false, fmt.Sprintf("Assertion failed: %d < %d", argInt, targetInt)
}
func le(arg string, target string) (bool, string) {
	argInt, err := strconv.Atoi(arg)
	if err != nil {
		return false, err.Error()
	}
	targetInt, err := strconv.Atoi(target)
	if err != nil {
		return false, err.Error()
	}
	if argInt < targetInt {
		return true, ""
	}
	return false, fmt.Sprintf("Assertion failed: %d >= %d", argInt, targetInt)
}
func eq(arg string, target string) (bool, string) {
	if arg == target {
		return true, ""
	}
	return false, fmt.Sprintf("Assertion failed: %s == %s", arg, target)
}
