package server

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"my-go-db/storage"
	"net/http"
	"time"
	"sync"
)

type Server struct {
	storage  *storage.Storage
	bindAddr string
	echo     *echo.Echo
	wg       *sync.WaitGroup
}

func New(bindAddr string) *Server {
	s := &Server{
		storage:  storage.New(),
		bindAddr: bindAddr,
		echo:     echo.New(),
		wg:       new(sync.WaitGroup),
	}

	g := s.echo.Group("/storage")
	g.GET("/", s.getKeys)
	g.GET("/:key", s.getValue)
	g.POST("/:key", s.setValue)
	g.DELETE("/:key", s.deleteValue)

	s.echo.Logger.SetLevel(log.DEBUG)
	return s
}

func (s *Server) Start() {
	s.wg.Add(1)
	go func() {
		err := s.echo.Start(s.bindAddr)
		fmt.Println("Server was stopped:", err.Error())
		s.wg.Done()
	}()
	go func() {
		for range time.Tick(1) {
			s.storage.DeleteExpired()
		}
	}()

}

func (s *Server) WaitStop() {
	s.wg.Wait()
}

// GET /storage/:key
func (s *Server) getValue(c echo.Context) error {
	key := c.Param("key")

	resp := new(ResponseBody)
	item := s.storage.GetItem(key)
	if item != nil {
		if item.String != "" {
			resp.Success = true
			resp.String = item.String
			return c.JSON(http.StatusOK, resp)
		} else if item.Int > 0 {
			resp.Success = true
			resp.Int = item.Int
			return c.JSON(http.StatusOK, resp)
		} else if item.StringSlice != nil {
			resp.Success = true
			resp.StringList = item.StringSlice
			return c.JSON(http.StatusOK, resp)
		} else if item.IntSlice != nil {
			resp.Success = true
			resp.IntList = item.IntSlice
			return c.JSON(http.StatusOK, resp)
		} else if item.StringMap != nil {
			resp.Success = true
			resp.StringDict = item.StringMap
			return c.JSON(http.StatusOK, resp)
		} else if item.IntMap != nil {
			resp.Success = true
			resp.IntDict = item.IntMap
			return c.JSON(http.StatusOK, resp)
		}
	}
	resp.Message = "Not found"
	return c.JSON(http.StatusNotFound, resp)
}

// POST /storage/:key
func (s *Server) setValue(c echo.Context) error {
	reqBody := RequestBody{}

	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, &ResponseBody{
			Success: false,
			Message: fmt.Sprintf("Could not set value: %v", err.Error()),
		})
	}

	key := c.Param("key")

	if reqBody.String != "" {
		s.storage.SetString(key, reqBody.String, reqBody.TTL)
	} else if reqBody.Int > 0 {
		s.storage.SetInt(key, reqBody.Int, reqBody.TTL)
	} else if reqBody.StringList != nil {
		s.storage.SetStringSlice(key, reqBody.StringList, reqBody.TTL)
	} else if reqBody.IntList != nil {
		s.storage.SetIntSlice(key, reqBody.IntList, reqBody.TTL)
	} else if reqBody.StringDict != nil {
		s.storage.SetStringMap(key, reqBody.StringDict, reqBody.TTL)
	} else {
		return c.JSON(http.StatusBadRequest, &ResponseBody{
			Success: false,
			Message: "Unsupported type",
		})
	}

	return c.JSON(http.StatusOK, &ResponseBody{
		Success: true,
		Message: "Done",
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

// GET /storage/
func (s *Server) getKeys(c echo.Context) error {
	keys := s.storage.Keys()
	return c.JSON(http.StatusOK, &ResponseBody{
		Success: true,
		Keys:   keys,
	})
}

// todo getFromList
// todo getFromDict
