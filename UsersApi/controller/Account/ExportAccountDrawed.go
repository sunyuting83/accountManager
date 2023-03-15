package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExportAccountDrawed(c *gin.Context) {
	var (
		date     string = c.Query("date")
		multiple string = c.DefaultQuery("multiple", "true")
		diamond  string = c.DefaultQuery("diamond", "false")
		crazy    string = c.DefaultQuery("crazy", "false")
		cold     string = c.DefaultQuery("cold", "false")
		precise  string = c.DefaultQuery("precise", "false")
		remarks  string = c.DefaultQuery("remarks", "false")

		Multiple bool = false
		Diamond  bool = false
		Crazy    bool = false
		Cold     bool = false
		Precise  bool = false
		Remarks  bool = false
	)
	if len(date) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "错误的日期格式",
		})
		return
	}

	if multiple == "true" {
		Multiple = true
	}
	if diamond == "true" {
		Diamond = true
	}
	if crazy == "true" {
		Crazy = true
	}
	if cold == "true" {
		Cold = true
	}
	if precise == "true" {
		Precise = true
	}
	if remarks == "true" {
		Remarks = true
	}

	projectsID := GetProjectsID(c)
	startTime, endTime := utils.GetSqlDateTime(date)

	data, err := database.ExportAccountDrawed(projectsID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "获取数据失败",
		})
		return
	}
	file := MakeDrawedExportFile(data, Multiple, Diamond, Crazy, Cold, Precise, Remarks)
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	// c.Writer.WriteString(string(file)) return image
	c.Data(200, "text/plain", file)
}

func MakeDrawedExportFile(data []*database.Accounts, Multiple, Diamond, Crazy, Cold, Precise, Remarks bool) []byte {
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
		itemStr = strings.TrimRight(itemStr, "\t")
		Gold := strconv.FormatInt(item.TodayGold, 10)
		itemStr = strings.Join([]string{itemStr, Gold}, "\t")
		if Multiple {
			itemStr = strings.TrimRight(itemStr, "\t")
			MultipleStr := strconv.FormatInt(item.Multiple, 10)
			itemStr = strings.Join([]string{itemStr, MultipleStr}, "\t")
		}
		if Diamond {
			itemStr = strings.TrimRight(itemStr, "\t")
			DiamondStr := strconv.Itoa(item.Diamond)
			itemStr = strings.Join([]string{itemStr, DiamondStr}, "\t")
		}
		if Crazy {
			itemStr = strings.TrimRight(itemStr, "\t")
			CrazyStr := strconv.Itoa(item.Crazy)
			itemStr = strings.Join([]string{itemStr, CrazyStr}, "\t")
		}
		if Precise {
			itemStr = strings.TrimRight(itemStr, "\t")
			PreciseStr := strconv.Itoa(item.Precise)
			itemStr = strings.Join([]string{itemStr, PreciseStr}, "\t")
		}
		if Cold {
			itemStr = strings.TrimRight(itemStr, "\t")
			ColdStr := strconv.Itoa(item.Cold)
			itemStr = strings.Join([]string{itemStr, ColdStr}, "\t")
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
		if Remarks {
			if len(item.Remarks) != 0 {
				itemStr = strings.TrimRight(itemStr, "\t")
				itemStr = strings.Join([]string{itemStr, item.Remarks}, "\t")
			}
		}
		itemStr = strings.TrimRight(itemStr, "\t")
		temp = append(temp, itemStr)
	}
	tempStr = strings.Join(temp, "\r\n")
	return []byte(tempStr)
}
