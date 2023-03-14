package controller

import (
	"colaAPI/UsersApi/database"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExportAccount(c *gin.Context) {
	var form DeleteList
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	projectsID := GetProjectsID(c)

	Projects, err := database.ProjectsCheckID(projectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	var statusJson []*StatusJSON
	json.Unmarshal([]byte(Projects.StatusJSON), &statusJson)

	var (
		hasPower bool = false
	)

	for _, item := range statusJson {
		if item.Status == form.Status {
			if item.Export {
				hasPower = true
				break
			}
		}
	}

	if !hasPower {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "状态不支持导出",
		})
		return
	}

	var account *database.Accounts
	data, err := account.ExportAccount(projectsID, form.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "获取数据失败",
		})
		return
	}
	file := MakeExportFile(data)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	// c.Writer.WriteString(string(file)) return image
	c.Data(200, "text/plain", file)
}

func MakeExportFile(data []*database.Accounts) []byte {
	var (
		temp    []string
		tempStr string
	)
	for _, item := range data {
		itemStr := strings.Join([]string{item.UserName}, "")
		if len(item.Password) != 0 {
			itemStr = strings.TrimRight(itemStr, "\t")
			itemStr = strings.Join([]string{itemStr, item.Password}, "\t")
		}
		if len(item.PhoneNumber) != 0 {
			itemStr = strings.TrimRight(itemStr, "\t")
			itemStr = strings.Join([]string{itemStr, item.PhoneNumber}, "\t")
		}
		if len(item.PhonePassword) != 0 {
			itemStr = strings.TrimRight(itemStr, "\t")
			itemStr = strings.Join([]string{itemStr, item.PhonePassword}, "\t")
		}
		if len(item.PhonePassword) != 0 {
			itemStr = strings.TrimRight(itemStr, "\t")
			itemStr = strings.Join([]string{itemStr, item.PhonePassword}, "\t")
		}
		if len(item.Remarks) != 0 {
			itemStr = strings.TrimRight(itemStr, "\t")
			itemStr = strings.Join([]string{itemStr, item.Remarks}, "\t")
		}
		itemStr = strings.TrimRight(itemStr, "\t")
		temp = append(temp, itemStr)
	}
	tempStr = strings.Join(temp, "\r\n")
	return []byte(tempStr)
}
