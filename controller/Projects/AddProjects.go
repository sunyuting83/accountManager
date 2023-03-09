package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/database"
	"colaAPI/utils"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Projects struct {
	UsersID      string `form:"usersid" json:"usersid" xml:"usersid"  binding:"required"`
	ProjectsName string `form:"ProjectsName" json:"ProjectsName" xml:"ProjectsName"  binding:"required"`
	UserName     string `form:"username" json:"username" xml:"username"`
	Password     string `form:"password" json:"password" xml:"password"`
	AccNumber    int    `form:"AccNumber" json:"AccNumber" xml:"AccNumber"`
	ColaAPI      string `form:"ColaAPI" json:"ColaAPI" xml:"ColaAPI"`
	StatusJSON   string `form:"StatusJSON" json:"StatusJSON" xml:"StatusJSON"`
}

type CacheValue struct {
	UsersID    string `json:"UsersID"`
	ProjectsID string `json:"ProjectsID"`
	ColaAPI    bool   `json:"ColaAPI"`
}

type StatusJSON struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

func AddProjects(c *gin.Context) {
	var form Projects
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	if len(form.UsersID) <= 0 || form.UsersID == "0" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't userid",
		})
		return
	}
	if len(form.ProjectsName) < 6 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  1,
			"message": "haven't projects name",
		})
		return
	}
	var Sjsons []byte
	if len(form.StatusJSON) == 0 {
		var Sjson []*StatusJSON
		Sjson = append(Sjson, &StatusJSON{
			Status: "0",
			Title:  "未注册状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "1",
			Title:  "注册中状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "2",
			Title:  "注册完成状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "3",
			Title:  "游戏中状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "4",
			Title:  "游戏完成状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "5",
			Title:  "封号状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "6",
			Title:  "旧帐号状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "7",
			Title:  "备用状态",
		})
		Sjson = append(Sjson, &StatusJSON{
			Status: "8",
			Title:  "提取状态",
		})
		Sjsons, _ = json.Marshal(&Sjson)
	}
	var ColaAPI1 bool = false
	if form.ColaAPI == "true" {
		ColaAPI1 = true
	}
	if ColaAPI1 {
		if len(form.UserName) < 5 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  1,
				"message": "haven't projects name",
			})
			return
		}
		if len(form.Password) < 6 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  1,
				"message": "haven't Password",
			})
			return
		}
		if form.AccNumber == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  1,
				"message": "haven't AccNumber",
			})
			return
		}
	}
	var projects *database.Projects
	UsersIDInt := StrToUInt(form.UsersID)
	if len(form.StatusJSON) == 0 {
		projects = &database.Projects{
			UsersID:      UsersIDInt,
			ProjectsName: form.ProjectsName,
			UserName:     form.UserName,
			Password:     form.Password,
			AccNumber:    form.AccNumber,
			NewStatus:    0,
			ColaAPI:      ColaAPI1,
			StatusJSON:   string(Sjsons),
		}
	} else {
		projects = &database.Projects{
			UsersID:      UsersIDInt,
			ProjectsName: form.ProjectsName,
			UserName:     form.UserName,
			Password:     form.Password,
			AccNumber:    form.AccNumber,
			NewStatus:    0,
			ColaAPI:      ColaAPI1,
			StatusJSON:   form.StatusJSON,
		}
	}
	err := projects.Insert()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	projectsIDInt := strconv.Itoa(int(projects.ID))
	projectsIDStr := string(projectsIDInt)
	d := time.Now()
	date := d.Format("2006-01-02_15:04:05")
	key := utils.MD5(strings.Join([]string{form.UsersID, date, projectsIDStr}, ""))
	key = key[:12]

	cache := &CacheValue{
		UsersID:    form.UsersID,
		ProjectsID: projectsIDStr,
		ColaAPI:    ColaAPI1,
	}
	CacheValues, _ := json.Marshal(&cache)

	Redis.Set(key, string(CacheValues), 0)

	projects.UpProjectsKey(key)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "添加成功",
		"data":    projects,
	})
}

func StrToUInt(str string) uint {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return uint(i)
}
