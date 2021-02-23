package task

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatusCodeAssertionSuccess(t *testing.T) {
	ast := assert.New(t)
	ass := StatusCodeAssertion{
		ExpectedCodes: []int{200, 302},
	}
	resp := HttpResponse{
		Status:     "success",
		StatusCode: 200,
	}
	success, _ := ass.Assert(resp)
	ast.True(success)
	resp = HttpResponse{
		Status:     "success",
		StatusCode: 302,
	}
	success, _ = ass.Assert(resp)
	ast.True(success)
}

func TestStatusCodeAssertionFailed(t *testing.T) {
	ast := assert.New(t)
	ass := StatusCodeAssertion{
		ExpectedCodes: []int{200},
	}
	resp := HttpResponse{
		Status:     "success",
		StatusCode: 400,
	}
	success, _ := ass.Assert(resp)
	ast.False(success)
}

var responseJsonBody = `{
    "store": {
        "book": [
            {
                "category": "reference",
                "author": "Nigel Rees",
                "title": "Sayings of the Century",
                "price": 8.95
            },
            {
                "category": "fiction",
                "author": "Evelyn Waugh",
                "title": "Sword of Honour",
                "price": 12.99
            },
            {
                "category": "fiction",
                "author": "Herman Melville",
                "title": "Moby Dick",
                "isbn": "0-553-21311-3",
                "price": 8.99
            },
            {
                "category": "fiction",
                "author": "J. R. R. Tolkien",
                "title": "The Lord of the Rings",
                "isbn": "0-395-19395-8",
                "price": 22.99
            }
        ],
        "bicycle": {
            "color": "red",
            "price": 19.95
        }
    },
    "expensive": 10
}`
var responseJsonBodyBytes = []byte(responseJsonBody)

func TestJsonPathAssertionSuccessForString(t *testing.T) {
	ast := assert.New(t)
	ass := JsonPathAssertion{
		Expression: "$.store.book[-1].isbn == 0-395-19395-8",
	}
	resp := HttpResponse{
		Status:     "success",
		StatusCode: 200,
		Body:       responseJsonBodyBytes,
	}
	success, _ := ass.Assert(resp)
	ast.True(success)
}
func TestJsonPathAssertionSuccessForNumber(t *testing.T) {
	ast := assert.New(t)
	ass := JsonPathAssertion{
		Expression: "$.expensive == 10",
	}
	resp := HttpResponse{
		Status:     "success",
		StatusCode: 200,
		Body:       responseJsonBodyBytes,
	}
	success, _ := ass.Assert(resp)
	ast.True(success)
}
