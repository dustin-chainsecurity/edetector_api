package login

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Init(c *gin.Context) {
	//msg := c.PostForm("input")
	msg := c.DefaultPostForm("input", "表單沒有input。") // 沒有輸入參數時 可設定預設值
	c.String(http.StatusOK, "您輸入的文字為: \n%s", msg)
}