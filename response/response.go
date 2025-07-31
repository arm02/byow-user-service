package response

import (
	"github.com/buildyow/byow-user-service/constants"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, gin.H{
		"status": constants.SUCCESS,
		"code":   code,
		"data":   data,
	})
}

func SuccessWithPagination(c *gin.Context, code int, data interface{}, total int64) {
	c.JSON(code, gin.H{
		"status":    constants.SUCCESS,
		"code":      code,
		"data":      data,
		"row_count": total,
	})
}

func Error(c *gin.Context, code int, message interface{}) {
	c.JSON(code, gin.H{
		"status": constants.ERROR,
		"code":   code,
		"data": gin.H{
			"message": message,
		},
	})
}
