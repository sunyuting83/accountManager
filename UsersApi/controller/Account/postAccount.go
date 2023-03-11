package controller

import (
	"colaAPI/UsersApi/database"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

// Node node
type NodeList struct {
	Data     string `form:"data" json:"data" xml:"data"  binding:"required"`
	SplitStr string `form:"splitstr" json:"splitstr" xml:"splitstr"  binding:"required"`
	Status   string `form:"status" json:"status" xml:"status" binding:"required"`
	HasMore  string `form:"hasmore" json:"hasmore" xml:"hasmore" binding:"required"`
}

func PostAccount(c *gin.Context) {
	var form NodeList
	// This will infer what binder to use depending on the content-type header.
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	StatusInt, err := strconv.Atoi(form.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	hasMore, err := strconv.Atoi(form.HasMore)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	linSplit := "\r\n"
	if !strings.Contains(form.Data, "\r") {
		linSplit = "\n"
	}
	if !strings.Contains(form.Data, "\n") {
		linSplit = "\r"
	}

	itemSplit := makeSplitStr(form.SplitStr)

	data := strings.Split(form.Data, linSplit)

	var account []database.Accounts
	projectsID := GetProjectsID(c)
	ProjectsID, _ := strconv.Atoi(projectsID)

	hasPhone, index := isPhone(data[0], itemSplit)

	for _, item := range data {
		itemS := strings.Split(item, itemSplit)
		if len(item) != 0 {
			if len(itemS) >= 1 {
				var (
					UserName      string = itemS[0]
					Password      string = ""
					PhoneNumber   string = ""
					PhonePassword string = ""
					TodayGold     int64  = 0
					Multiple      int64  = 0
					Diamond       int    = 0
					Crazy         int    = 0
					Precise       int    = 0
					Cold          int    = 0
				)
				if len(itemS) >= 1 {
					Password = itemS[1]
				}
				if hasPhone {
					PhoneNumber = itemS[index]
					PhonePassword = itemS[index+1]
				}
				if hasMore == 1 {
					TodayGold, _ = strconv.ParseInt(itemS[2], 10, 64)
					Multiple, _ = strconv.ParseInt(itemS[3], 10, 64)
					Diamond, _ = strconv.Atoi(itemS[4])
					Crazy, _ = strconv.Atoi(itemS[5])
					Precise, _ = strconv.Atoi(itemS[6])
					Cold, _ = strconv.Atoi(itemS[7])
				}
				account = append(account, database.Accounts{
					ProjectsID:    uint(ProjectsID),
					ComputID:      0,
					PhoneNumber:   PhoneNumber,
					PhonePassword: PhonePassword,
					UserName:      UserName,
					Password:      Password,
					Cover:         "",
					NewStatus:     StatusInt,
					TodayGold:     TodayGold,
					YesterdayGold: 0,
					Multiple:      Multiple,
					Diamond:       Diamond,
					Crazy:         Crazy,
					Precise:       Precise,
					Cold:          Cold,
					Exptime:       0,
					Price:         0,
					Remarks:       "",
				})
			}
		}
	}
	batchLen := len(account)
	if batchLen > 1000 {
		database.AccountInBatches(account)
	} else {
		database.AccountBatches(account)
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "导入成功",
	})
}

func makeSplitStr(s string) string {
	var x string = "\t"
	switch s {
	case "0":
		x = "\t"
	case "1":
		x = "----"
	case "2":
		x = " "
	default:
		x = s
	}
	return x
}

func isDigit(str string) bool {
	for _, x := range str {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}

func isPhone(s, itemSplit string) (has bool, i int) {
	var x bool = false
	var index int = 0
	itemS := strings.Split(s, itemSplit)
	for i, item := range itemS {
		if len(item) == 11 {
			if isDigit(item) {
				x = true
				index = i
				break
			}
		}
	}
	return x, index
}
