package server

import (
	"fmt"
	"github.com/labstack/echo"
	"my-go-db/storage"
	"net/http"
	"time"
)

type Server struct {
	storage  *storage.Storage
	bindAddr string
	echo     *echo.Echo
}

func New(bindAddr string) *Server {
	s := &Server{
		storage:  storage.New(),
		bindAddr: bindAddr,
		echo:     echo.New(),
	}
	g := s.echo.Group("/storage")
	g.GET("/", s.getKeys)
	g.GET("/:key", s.getValue)
	g.POST("/:key", s.setValue)
	g.DELETE("/:key", s.deleteValue)
	return s
}

func (s *Server) Start() {
	go func() {
		s.echo.Logger.Fatal(s.echo.Start(s.bindAddr))
	}()
	go func() {
		for range time.Tick(1) {
			s.storage.DeleteExpired()
		}
	}()

}

// GET /storage/:key
func (s *Server) getValue(c echo.Context) error {
	key := c.Param("key")

	resp := new(ResponseBody)
	status := http.StatusNotFound
	item := s.storage.GetItem(key)
	if item != nil {
		if item.ValueInt > 0 {
			status = http.StatusOK
			resp.ValueInt = item.ValueInt
		} else if item.ValueStr != "" {
			status = http.StatusOK
			resp.ValueStr = item.ValueStr
		}
	}
	resp.Message = "Not found"
	return c.JSON(status, resp)
}

// POST /storage/:key
func (s *Server) setValue(c echo.Context) error {
	req := RequestBody{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, &ResponseBody{
			Success: false,
			Message: fmt.Sprintf("Could not set value: %v", err.Error()),
		})
	}

	key := c.Param("key")

	var value interface{}
	fmt.Println(req.ValueStr)
	if req.ValueStr != "" {
		value = req.ValueStr
		s.storage.SetString(key, req.ValueStr, req.TTL)
	} else if req.ValueInt > 0 {
		value = req.ValueInt
		s.storage.SetInt(key, req.ValueInt, req.TTL)
	}
	//if req.ItemText != "" {
	//	value = req.ItemText
	//} else if req.ItemTextArray != nil {
	//	value = req.ItemTextArray
	//} else if req.ItemTextDict != nil {
	//	value = req.ItemTextDict
	//}

	//s.storage.Set(key, value, req.TTL)

	return c.JSON(http.StatusOK, &ResponseBody{
		Success: true,
		Message: fmt.Sprintf("Set value: %v", value),
	})
}

func (s *Server) deleteValue(c echo.Context) error {
	key := c.Param("key")
	if err := s.storage.Remove(key); err != nil {
		return c.JSON(http.StatusNotFound, &ResponseBody{
			Success: false,
			Message: fmt.Sprintf("Key: %s does not exist", key),
		})
	}
	return c.JSON(http.StatusOK, &ResponseBody{
		Success: true,
	})
}

func (s *Server) getKeys(c echo.Context) error {
	keys := s.storage.Keys()
	return c.JSON(http.StatusOK, &ResponseBody{
		Success: true,
		Value:   keys,
	})
}

// todo getFromList
// todo getFromDict
