package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitData(t *testing.T) {
	ast := assert.New(t)
	s := &SampleServer{
		Port: 1234,
	}
	s.InitData()
	ast.NotNil(s.Data)
	ast.Len(s.Data, 4)
}
