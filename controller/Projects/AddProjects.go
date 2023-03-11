package controller

import (
	Redis "colaAPI/Redis"
	"colaAPI/database"
	"colaAPI/utils"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

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
	Status   string `json:"status"`
	Title    string `json:"title"`
	Delete   bool   `json:"delete"`
	CallBack bool   `json:"callback"`
	BackTo   string `json:"backto"`
	Export   bool   `json:"export"`
	Import   bool   `json:"import"`
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

	StatusJSON := MakeStatusJSON(form.StatusJSON)

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
	UsersIDInt := StrToUInt(form.UsersID)

	projects := &database.Projects{
		UsersID:      UsersIDInt,
		ProjectsName: form.ProjectsName,
		UserName:     form.UserName,
		Password:     form.Password,
		AccNumber:    form.AccNumber,
		NewStatus:    0,
		ColaAPI:      ColaAPI1,
		StatusJSON:   StatusJSON,
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
	key = MakeKey(key)

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

func MakeKey(a string) string {
	b := strings.Split(a, "")
	c := make([]string, 12)
	for _, v := range b {
		if unicode.IsLetter([]rune(v)[0]) {
			n, _ := rand.Int(rand.Reader, big.NewInt(2))
			if n.Int64() == 0 {
				v = strings.ToUpper(v)
			}
		}
		c = append(c, v)
	}
	return strings.Join(c, "")
}

func MakeStatusJSON(JSON string) (statusJSON string) {
	var Sjsons []byte
	statusJSON = JSON
	if len(JSON) == 0 {
		Sjson := &[]StatusJSON{
			{
				Status:   "0",
				Title:    "未注册状态",
				Delete:   true,
				CallBack: false,
				BackTo:   "",
				Export:   true,
				Import:   true,
			},
			{
				Status:   "1",
				Title:    "注册中状态",
				Delete:   false,
				CallBack: true,
				BackTo:   "0",
				Export:   false,
				Import:   false,
			},
			{
				Status:   "2",
				Title:    "注册完成状态",
				Delete:   false,
				CallBack: false,
				BackTo:   "",
				Export:   false,
				Import:   true,
			},
			{
				Status:   "3",
				Title:    "游戏中状态",
				Delete:   false,
				CallBack: true,
				BackTo:   "2",
				Export:   false,
				Import:   false,
			},
			{
				Status:   "4",
				Title:    "游戏完成状态",
				Delete:   false,
				CallBack: true,
				BackTo:   "2",
				Export:   false,
				Import:   false,
			},
			{
				Status:   "5",
				Title:    "封号状态",
				Delete:   false,
				CallBack: false,
				BackTo:   "0",
				Export:   true,
				Import:   false,
			},
			{
				Status:   "6",
				Title:    "旧帐号状态",
				Delete:   false,
				CallBack: false,
				BackTo:   "",
				Export:   true,
				Import:   false,
			},
			{
				Status:   "7",
				Title:    "备用状态",
				Delete:   false,
				CallBack: false,
				BackTo:   "",
				Export:   true,
				Import:   false,
			},
			{
				Status:   "8",
				Title:    "未使用身份证",
				Delete:   true,
				CallBack: false,
				BackTo:   "",
				Export:   true,
				Import:   true,
			},
			{
				Status:   "9",
				Title:    "使用中身份证",
				Delete:   false,
				CallBack: true,
				BackTo:   "8",
				Export:   false,
				Import:   false,
			},
			{
				Status:   "10",
				Title:    "已使用身份证",
				Delete:   true,
				CallBack: true,
				BackTo:   "8",
				Export:   true,
				Import:   false,
			},
			{
				Status:   "108",
				Title:    "提取状态",
				Delete:   false,
				CallBack: false,
				BackTo:   "",
				Export:   true,
				Import:   false,
			},
		}
		Sjsons, _ = json.Marshal(&Sjson)
		statusJSON = string(Sjsons)
	}
	return
}
