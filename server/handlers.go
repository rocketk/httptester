package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *SampleServer) Root(c *gin.Context) {
	c.JSON(200, &gin.H{
		"users_url": fmt.Sprintf("http://localhost:%d/users", s.Port),
	})
}

func (s *SampleServer) List(c *gin.Context) {

	c.JSON(200, s.Data.GetAllOrdered())
}

func (s *SampleServer) Get(c *gin.Context) {
	id := c.Param("id")
	if user, ok := s.Data.Get(id); ok {
		c.JSON(200, user)
		return
	}
	c.JSON(404, GenericError{Status: 404, Message: "User not found"})
}

func (s *SampleServer) Delete(c *gin.Context) {
	id := c.Param("id")
	if user, ok := s.Data.Get(id); ok {
		s.Data.Delete(id)
		c.JSON(200, user)
		return
	}
	c.JSON(404, GenericError{Status: 404, Message: "User not found"})
}

func (s *SampleServer) Put(c *gin.Context) {
	id := c.Param("id")
	var reqUser User
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		c.JSON(400, GenericError{Status: 400, Message: "Invalid User Body"})
		return
	}
	reqUser.ID = id
	if _, ok := s.Data.Get(id); !ok {
		c.JSON(404, GenericError{Status: 404, Message: "User not found"})
		return
	}
	s.Data.Put(id, reqUser)
	c.JSON(200, reqUser)
}

func (s *SampleServer) Post(c *gin.Context) {
	var reqUser User
	if err := c.ShouldBindJSON(&reqUser); err != nil {
		log.Println(err)
		c.JSON(400, GenericError{Status: 400, Message: "Invalid User Body"})
		return
	}
	id := uuid.New().String()
	reqUser.ID = id
	s.Data.Put(id, reqUser)
	c.JSON(200, reqUser)
}

type GenericError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
