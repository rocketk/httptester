package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
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
		// log.Printf("Actual: %d, Expected: %d, response body: %s", resp.StatusCode, code, resp.Body)
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
	return a.handler(arg, a.value)
	// var argInStr string
	// switch arg.(type) {
	// case float32, float64:
	// 	argInStr = fmt.Sprintf("%f", arg)
	// case int, int32, int64:
	// 	argInStr = fmt.Sprintf("%d", arg)
	// default:
	// 	argInStr = fmt.Sprintf("%v", arg)
	// }
	// return a.handler(argInStr, a.value)
}

func (a JsonPathAssertion) Name() string {
	return "JsonPathAssertion"
}

type verbHandler func(arg interface{}, target string) (bool, string)

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
	// log.Printf("%s\n", a.value)
	handler, ok := verbFuncMap[a.verb]
	if !ok {
		return fmt.Errorf("No handler was found for the verb '%s'", a.verb)
	}
	a.handler = handler
	a.initialized = true
	return nil
}

func gt(arg interface{}, target string) (bool, string) {
	switch arg.(type) {
	case int:
		targetInt, err := strconv.Atoi(target)
		argInt, _ := arg.(int)
		if err != nil {
			return false, err.Error()
		}
		return argInt > targetInt, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float32:
		targetFloat64, err := strconv.ParseFloat(target, 32)
		argFloat32, _ := arg.(float32)
		targetFloat32 := float32(targetFloat64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat32 > targetFloat32, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float64:
		targetFloat64, err := strconv.ParseFloat(target, 64)
		argFloat64, _ := arg.(float64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat64 > targetFloat64, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	}
	return false, fmt.Sprintf("Assertion failed: missmatched type: arg %T, target %T", arg, target)
}
func ge(arg interface{}, target string) (bool, string) {
	switch arg.(type) {
	case int:
		targetInt, err := strconv.Atoi(target)
		argInt, _ := arg.(int)
		if err != nil {
			return false, err.Error()
		}
		return argInt >= targetInt, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float32:
		targetFloat64, err := strconv.ParseFloat(target, 32)
		argFloat32, _ := arg.(float32)
		targetFloat32 := float32(targetFloat64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat32 >= targetFloat32, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float64:
		targetFloat64, err := strconv.ParseFloat(target, 64)
		argFloat64, _ := arg.(float64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat64 >= targetFloat64, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	}
	return false, fmt.Sprintf("Assertion failed: missmatched type: arg %T, target %T", arg, target)
}
func lt(arg interface{}, target string) (bool, string) {
	switch arg.(type) {
	case int:
		targetInt, err := strconv.Atoi(target)
		argInt, _ := arg.(int)
		if err != nil {
			return false, err.Error()
		}
		return argInt < targetInt, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float32:
		targetFloat64, err := strconv.ParseFloat(target, 32)
		argFloat32, _ := arg.(float32)
		targetFloat32 := float32(targetFloat64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat32 < targetFloat32, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float64:
		targetFloat64, err := strconv.ParseFloat(target, 64)
		argFloat64, _ := arg.(float64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat64 < targetFloat64, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	}
	return false, fmt.Sprintf("Assertion failed: missmatched type: arg %T, target %T", arg, target)
}
func le(arg interface{}, target string) (bool, string) {
	switch arg.(type) {
	case int:
		targetInt, err := strconv.Atoi(target)
		argInt, _ := arg.(int)
		if err != nil {
			return false, err.Error()
		}
		return argInt <= targetInt, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float32:
		targetFloat64, err := strconv.ParseFloat(target, 32)
		argFloat32, _ := arg.(float32)
		targetFloat32 := float32(targetFloat64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat32 <= targetFloat32, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float64:
		targetFloat64, err := strconv.ParseFloat(target, 64)
		argFloat64, _ := arg.(float64)
		if err != nil {
			return false, err.Error()
		}
		return argFloat64 <= targetFloat64, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	}
	return false, fmt.Sprintf("Assertion failed: missmatched type: arg %T, target %T", arg, target)
}
func eq(arg interface{}, target string) (bool, string) {
	// log.Printf("arg: %s, target: %s\n", arg, target)
	switch arg.(type) {
	case int:
		targetInt, err := strconv.Atoi(target)
		argInt, _ := arg.(int)
		if err != nil {
			return false, err.Error()
		}
		return argInt == targetInt, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float32:
		targetFloat32, err := strconv.ParseFloat(target, 32)
		if err != nil {
			return false, err.Error()
		}
		return arg == targetFloat32, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case float64:
		targetFloat64, err := strconv.ParseFloat(target, 64)
		if err != nil {
			return false, err.Error()
		}
		return arg == targetFloat64, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case bool:
		targetBool, err := strconv.ParseBool(target)
		if err != nil {
			return false, err.Error()
		}
		return arg == targetBool, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	case string:
		argStr := arg.(string)
		if argStr == target {
			return true, ""
		}
		return false, fmt.Sprintf("Assertion failed: Actual: %s, Expected: %s", arg, target)
	}
	return false, fmt.Sprintf("Assertion failed: missmatched type: arg %T (%v), target %T (%v)", arg, arg, target, target)
}

type RegexAssertion struct {
	Expression string
}

func (r RegexAssertion) Assert(resp HttpResponse) (bool, string) {
	reg := regexp.MustCompile(r.Expression)
	match := reg.Match(resp.Body)
	if match {
		return true, ""
	}
	return false, fmt.Sprintf("Assertion failed: response body does not match the regex excpression: %s", r.Expression)
}

func (r RegexAssertion) Name() string {
	return "RegexAssertion"
}

func (r RegexAssertion) Validate() error {
	if len(r.Expression) == 0 {
		return errors.New("Expression is required")
	}
	return nil
}
