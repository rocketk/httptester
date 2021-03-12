package server

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SampleServer struct {
	Port int
	Data OrderedUserMap
}
type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Age       int     `json:"age"`
	Stature   int     `json:"stature"`
	Weight    float32 `json:"weight"`
	Available bool    `json:"available"`
}

//go:embed sample_data.json
var sampleData []byte

func (s *SampleServer) InitData() {
	var users []User
	err := json.Unmarshal(sampleData, &users)
	if err != nil {
		panic(err)
	}
	// log.Printf("%+v", users)
	s.Data = NewOrderedUserMap(32)
	for _, u := range users {
		s.Data.Put(u.ID, u)
	}
}

func (s *SampleServer) Serve() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", s.Root)
	r.GET("/users", s.List)
	r.GET("/users/:id", s.Get)
	r.DELETE("/users/:id", s.Delete)
	r.PUT("/users/:id", s.Put)
	r.POST("/users", s.Post)
	fmt.Printf("sample restful-api server started successfully! try it out: \ncurl http:127.0.0.1:%d\n", s.Port)

	// fmt.Printf("listening on %d\n", s.Port)
	r.Run(fmt.Sprintf(":%d", s.Port))
}
