package response

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

func Success(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Code:    statusCode,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Code:    statusCode,
		Message: message,
	})
}

func Paginated(c *gin.Context, statusCode int, data interface{}, total int64, page, limit int) {
	c.JSON(statusCode, PaginatedResponse{
		Code:    statusCode,
		Message: "success",
		Data:    data,
		Total:   total,
		Page:    page,
		Limit:   limit,
	})
}
