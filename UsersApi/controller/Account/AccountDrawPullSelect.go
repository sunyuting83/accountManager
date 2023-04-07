package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PullAccountDrawSelect(c *gin.Context) {
	var form utils.Filter
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	projectsID, ColaAPI := GetProjects(c)
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)

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

	rows, _ := database.GetDataUseScopes(form, hasStatus, projectsID)

	if len(rows) != 0 {
		IDs := make([]int, len(rows))
		for _, item := range rows {
			IDs = append(IDs, int(item.ID))
		}
		AdminID := utils.GetCurrentUserID(c)
		user, err := database.UserCheckID(int64(AdminID))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "get admin failed",
			})
			return
		}
		upData, err := database.PullDataUseIn(IDs)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err,
			})
			return
		}

		upDataJsonStr, _ := json.Marshal(&upData)

		d := time.Now()
		date := d.Format("2006-01-02_15:04:05")
		draw := &database.DrawLogs{
			ProjectsID: uint(ProjectsID),
			Data:       string(upDataJsonStr),
			LogName:    date,
			DrawUser:   user.UserName,
		}
		draw.AddDrawLogs()
		if ColaAPI {
			Projects = &database.Projects{
				UserName:  Projects.UserName,
				Password:  Projects.Password,
				AccNumber: Projects.AccNumber - int(len(rows)),
			}
			Projects.UpdateProjects(projectsID)
		}

		Data := gin.H{
			"status":  0,
			"message": "提取成功",
			"data":    rows,
		}
		c.JSON(http.StatusOK, Data)
		return
	}

	empty := make([]string, 0)
	Data := gin.H{
		"status":  1,
		"message": "没有符合条件的帐号",
		"data":    empty,
	}
	c.JSON(http.StatusOK, Data)
}

/*
func MakeSelectSQL(filter Filter, hasStatusStr, projectsID string) (SQL string) {
	nowTime := time.Now().Unix() * 1000
	TimeStr := strconv.FormatInt(nowTime, 10)
	SQL = "UPDATE accounts SET updated_at = " + TimeStr + ", new_status = 108 WHERE projects_id = " + projectsID + " AND new_status IN ("
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
*/
