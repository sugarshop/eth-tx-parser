package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sugarshop/token-gateway/model"
	"net/http"
)

// JSONWrapper Encapsulate the data processing function as a JSON API return;
// pay attention to using PureJSON to avoid Gin performing HTML escaping during serialization of data.
func JSONWrapper(fn func(*gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		data, err := fn(c)

		// if write is been writen at here, do not write again or you will get panic
		if c.Writer.Written() {
			return
		}

		if err != nil {
			//c.Set(tracing.CtxRespCodeKey, base.FAILED)
			c.PureJSON(http.StatusOK, &ErrResp{
				Code:   model.RESPONSE_FAILD,
				Msg:    err.Error(),
				Detail: err.Error(),
			})
			return
		}

		c.PureJSON(http.StatusOK, &DataResp{
			Code: model.RESPONSE_OK,
			Data: data,
		})
	}
}

type ErrResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
}

type DataResp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}
