package controller

import (
	"colaAPI/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectsModify struct {
	ID         string `form:"id" json:"id" xml:"id"  binding:"required"`
	UserName   string `form:"username" json:"username" xml:"username"`
	Password   string `form:"password" json:"password" xml:"password"`
	AccNumber  int    `form:"AccNumber" json:"AccNumber" xml:"AccNumber"`
	ColaAPI    string `form:"ColaAPI" json:"ColaAPI" xml:"ColaAPI"`
	StatusJSON string `form:"StatusJSON" json:"StatusJSON" xml:"StatusJSON"`
}

func ModifyProjects(c *gin.Context) {
	var form ProjectsModify
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
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
	fmt.Println(ColaAPI1)
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
	if len(form.StatusJSON) == 0 {
		projects = &database.Projects{
			UserName:   form.UserName,
			Password:   form.Password,
			AccNumber:  form.AccNumber,
			ColaAPI:    ColaAPI1,
			StatusJSON: string(Sjsons),
		}
	} else {
		projects = &database.Projects{
			UserName:   form.UserName,
			Password:   form.Password,
			AccNumber:  form.AccNumber,
			ColaAPI:    ColaAPI1,
			StatusJSON: form.StatusJSON,
		}
	}
	projects.UpdateProjects(form.ID)
	ID, _ := strconv.ParseInt(form.ID, 10, 64)
	data, _ := database.ProjectsCheckID(ID)

	c.JSON(http.StatusOK, gin.H{
		"status":  0,
		"message": "修改成功",
		"data":    data,
	})
}
