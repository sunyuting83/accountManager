package controller

import (
	"colaAPI/database"
	"colaAPI/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type List struct {
	List string `form:"list" json:"list" xml:"list"  binding:"required"`
}
type JsonData struct {
	PID int `json:"pid"`
	AID int `json:"aid"`
}

type SqlData struct {
	PID  int64 `json:"pid"`
	AIDs []int `json:"aid"`
}

type ReturnData struct {
	Projects *database.Projects
	Accounts []*database.Accounts
}

func DrawSelectPull(c *gin.Context) {
	userData := utils.GetTokenUserData(c)
	if userData.UserID == 1 {
		var form List
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		if len(form.List) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": "game id must",
			})
			return
		}
		var j []*JsonData
		json.Unmarshal([]byte(form.List), &j)

		if len(j) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": "game id must",
			})
			return
		}
		pidList := RemoveRepeatedList(j)

		x := IgnoreRepeated(pidList, j)

		var (
			rtData []*ReturnData = make([]*ReturnData, 0)
			errs   error
			Status bool = false
		)

		for _, pids := range x {
			Projects, err := database.ProjectsCheckIDWithJoin(pids.PID)
			if err != nil {
				errs = err
				Status = true
				break
			}
			upData, err := database.PullDataUseIn(pids.AIDs)
			if err != nil {
				errs = err
				Status = true
				break
			}

			upDataJsonStr, err := json.Marshal(&upData)
			if err != nil {
				errs = err
				Status = true
				break
			}

			d := time.Now()
			date := d.Format("2006-01-02_15:04:05")
			draw := &database.DrawLogs{
				ProjectsID: uint(pids.PID),
				Data:       string(upDataJsonStr),
				LogName:    date,
				DrawUser:   "admin",
			}
			draw.AddDrawLogs()
			rt := &ReturnData{
				Projects: Projects,
				Accounts: upData,
			}
			rtData = append(rtData, rt)
		}
		// empty := make([]string, 0)
		Data := gin.H{
			"status":  0,
			"message": "提取成功",
			"data":    rtData,
		}
		if Status {
			Data = gin.H{
				"status":  1,
				"message": errs.Error(),
				"data":    rtData,
			}
			c.JSON(http.StatusOK, Data)
			return
		}
		c.JSON(http.StatusOK, Data)
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  1,
		"message": "haven't power",
	})
}

func RemoveRepeatedList(personList []*JsonData) (result []*JsonData) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i].PID == personList[j].PID {
				repeat = true
				break
			}
		}
		if !repeat && personList[i].PID != 0 {
			result = append(result, personList[i])
		}
	}
	return
}

func IgnoreRepeated(postList []*JsonData, dataList []*JsonData) []*SqlData {

	var x []*SqlData
	for _, item := range postList {
		y := make([]int, 0)
		for _, ig := range dataList {
			if item.PID == ig.PID {
				y = append(y, ig.AID)
			}
		}
		x = append(x, &SqlData{
			PID:  int64(item.PID),
			AIDs: y,
		})
	}
	// fmt.Println(temp)
	return x
}
