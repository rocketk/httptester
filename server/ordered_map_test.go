package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedUserMap(t *testing.T) {
	ast := assert.New(t)
	usersMap := NewOrderedUserMap(32)
	ast.Equal(0, usersMap.Size())
	usersMap.Put("a", User{Name: "aa"})
	usersMap.Put("b", User{Name: "bb"})
	usersMap.Put("c", User{Name: "cc"})
	ast.Equal(3, usersMap.Size())
	usersMap.Delete("b")
	ast.Equal(2, usersMap.Size())
	a, ok := usersMap.Get("a")
	ast.True(ok)
	ast.Equal(User{Name: "aa"}, a)
	_, ok = usersMap.Get("b")
	ast.False(ok)
	c, ok := usersMap.Get("c")
	ast.True(ok)
	ast.Equal(User{Name: "cc"}, c)

	usersMap.Delete("b")
	ast.Equal(2, usersMap.Size())
	fmt.Printf("%+v", usersMap.GetAllOrdered())
}
