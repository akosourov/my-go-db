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

func (s *Server) getValue(c echo.Context) error {
	key := c.Param("key")
	value := s.storage.Get(key)
	if value == nil {
		return c.JSON(http.StatusNotFound, &Response{
			Success: false,
			Message: fmt.Sprintf("Key: %s does not exist", key),
		})
	}
	return c.JSON(http.StatusOK, &Response{
		Success: true,
		Value:   value,
	})
}

func (s *Server) setValue(c echo.Context) error {
	p := Payload{}
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: fmt.Sprintf("Could not set value: %v", err.Error()),
		})
	}

	var value interface{}
	if p.ItemText != "" {
		value = p.ItemText
	} else if p.ItemTextArray != nil {
		value = p.ItemTextArray
	} else if p.ItemTextDict != nil {
		value = p.ItemTextDict
	}

	key := c.Param("key")
	s.storage.Set(key, value, p.TTL)

	return c.JSON(http.StatusOK, &Response{
		Success: true,
		Message: fmt.Sprintf("Set value: %v", value),
	})
}

func (s *Server) deleteValue(c echo.Context) error {
	key := c.Param("key")
	if err := s.storage.Remove(key); err != nil {
		return c.JSON(http.StatusNotFound, &Response{
			Success: false,
			Message: fmt.Sprintf("Key: %s does not exist", key),
		})
	}
	return c.JSON(http.StatusOK, &Response{
		Success: true,
	})
}

func (s *Server) getKeys(c echo.Context) error {
	keys := s.storage.Keys()
	return c.JSON(http.StatusOK, &Response{
		Success: true,
		Value:   keys,
	})
}

// todo getFromList
// todo getFromDict
