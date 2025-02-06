package girm

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Compatible with single object or multiple objects
func operation[T any](c *gin.Context, opr func(...*T) error) {
	var err error
	var rsp any
	defer func() {
		if err != nil {
			JsonFail(c, err.Error())
		} else {
			JsonOK(c, rsp)
		}
	}()
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	//write jsonData back to request body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonData))

	if jsonData[0] == '[' { //array
		var es []*T
		if err = c.ShouldBindJSON(&es); err != nil {
			return
		}
		if err = opr(es...); err != nil {
			return
		}
		rsp = es
	} else { //object
		var e T
		if err = c.ShouldBindJSON(&e); err != nil {
			return
		}
		if err = opr(&e); err != nil {
			return
		}
		rsp = e
	}
}

func IndexHandler(c *gin.Context)  {
	c.HTML(http.StatusOK, "index.html", nil)
}

func JsonOK(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": obj,
	})
}

func JsonFail(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, gin.H{
		"code":    500,
		"message": obj,
	})
}
