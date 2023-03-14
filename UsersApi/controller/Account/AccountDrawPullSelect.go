package controller

import (
	"colaAPI/UsersApi/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Filter struct {
	MinGold  int64 `form:"mingold" json:"mingold" xml:"mingold"  binding:"required"`
	MaxGold  int64 `form:"maxgold" json:"maxgold" xml:"maxgold"  binding:"required"`
	Multiple int64 `form:"multiple" json:"multiple" xml:"multiple"  binding:"required"`
	Diamond  int64 `form:"diamond" json:"diamond" xml:"diamond"`
	Crazy    int64 `form:"crazy" json:"crazy" xml:"crazy"`
	Cold     int64 `form:"cold" json:"cold" xml:"cold"`
	Precise  int64 `form:"precise" json:"precise" xml:"precise"`
}

func PullAccountDrawSelect(c *gin.Context) {
	var form Filter
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
		hasStatus []string
	)
	for _, item := range statusJson {
		if item.Pull {
			hasStatus = append(hasStatus, item.Status)
		}
	}
	hasStatusStr := strings.Join(hasStatus, ",")

	SQL := MakeSelectSQL(form, hasStatusStr)
	fmt.Println(SQL)
	var acc *database.Accounts

	acc.PullDataUseSQL(SQL)

	Data := gin.H{
		"status":  0,
		"message": "提取成功",
	}
	c.JSON(http.StatusOK, Data)
}

func MakeSelectSQL(filter Filter, hasStatusStr string) (SQL string) {
	SQL = "UPDATE accounts SET new_status = 108 WHERE new_status IN ("
	SQL = strings.Join([]string{SQL, hasStatusStr, ") AND "}, "")
	if filter.MinGold > 0 {
		MinGold := strconv.FormatInt(filter.MinGold, 10)
		salveSQL := strings.Join([]string{"today_gold >= ", MinGold, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	if filter.MaxGold > 0 {
		MaxGold := strconv.FormatInt(filter.MaxGold, 10)
		salveSQL := strings.Join([]string{"today_gold <= ", MaxGold, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	if filter.Multiple > 0 {
		Multiple := strconv.FormatInt(filter.Multiple, 10)
		salveSQL := strings.Join([]string{"multiple >= ", Multiple, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	if filter.Diamond > 0 {
		Diamond := strconv.FormatInt(filter.Diamond, 10)
		salveSQL := strings.Join([]string{"diamond >= ", Diamond, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	if filter.Crazy > 0 {
		Crazy := strconv.FormatInt(filter.Crazy, 10)
		salveSQL := strings.Join([]string{"crazy >= ", Crazy, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	if filter.Cold > 0 {
		Cold := strconv.FormatInt(filter.Cold, 10)
		salveSQL := strings.Join([]string{"cold >= ", Cold, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	if filter.Precise > 0 {
		Precise := strconv.FormatInt(filter.Precise, 10)
		salveSQL := strings.Join([]string{"precise >= ", Precise, " AND "}, "")
		SQL = strings.Join([]string{SQL, salveSQL}, "")
	}
	SQL = strings.TrimRight(SQL, " AND ")
	return
}
