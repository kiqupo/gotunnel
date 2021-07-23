package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"
)

var testpath = "D:/goland-2021.1.3.exe"

func Download(c *gin.Context) {
	log.Println("[用户请求下载文件]")
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(testpath)))
	c.File(testpath)
}
