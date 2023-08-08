package controller

import (
	"bytes"
	Accounts "colaAPI/controller/Account"
	"colaAPI/database"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Node node
type NodeList struct {
	Data     string `form:"data" json:"data" xml:"data"`
	SplitStr string `form:"splitstr" json:"splitstr" xml:"splitstr"  binding:"required"`
	Status   string `form:"status" json:"status" xml:"status" binding:"required"`
	HasFile  string `form:"hasfile" json:"hasfile" xml:"hasfile" binding:"required"`
	Repeated string `form:"repeated" json:"repeated" xml:"repeated" binding:"required"`
}

type StatusJSON struct {
	Status   string `json:"status"`
	Title    string `json:"title"`
	Delete   bool   `json:"delete"`
	CallBack bool   `json:"callback"`
	BackTo   string `json:"backto"`
	Export   bool   `json:"export"`
	Import   bool   `json:"import"`
	Pull     bool   `json:"pull"`
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
	HasFile := false
	if form.HasFile == "true" {
		HasFile = true
	}
	Repeated := false
	if form.Repeated == "true" {
		Repeated = true
	}

	if !HasFile {
		if len(form.Data) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  1,
				"message": "数据不能为空",
			})
			return
		}
	}

	projectsID := Accounts.GetProjectsID(c)
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
	_, err = database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}

	Data := form.Data

	if HasFile {
		file, handler, err := c.Request.FormFile("files")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "上传文件失败" + err.Error(),
			})
			return
		}
		extList := strings.Split(handler.Filename, ".")
		extLen := len(extList) - 1
		ext := extList[extLen]
		if ext != "txt" {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "必须是.txt文件",
			})
			return
		}

		b, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": err.Error(),
			})
			return
		}
		if len(b) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "别传空文件",
			})
			return
		}
		_, what, certain := charset.DetermineEncoding(b, "txt")

		if !certain && what != "utf-8" {
			// fmt.Println(what)
			a, _ := GbkToUtf8(b)
			Data = string(a)
		} else {
			Data = string(b)
		}
	}

	linSplit := "\r\n"
	if !strings.Contains(Data, "\r") {
		linSplit = "\n"
	}
	if !strings.Contains(Data, "\n") {
		linSplit = "\r"
	}

	itemSplit := makeSplitStr(form.SplitStr)

	data := strings.Split(Data, linSplit)
	dataLen := len(data)
	if dataLen == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "别传空文件",
		})
		return
	}
	if dataLen > 30000 {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": "数据超过30000行",
		})
		return
	}

	var others []database.Others

	for _, item := range data {
		itemS := strings.Split(item, itemSplit)
		if len(item) != 0 {
			if len(itemS) >= 1 {
				var (
					Account  string = itemS[0]
					Password string = ""
					Remarks  string = ""
				)
				if len(itemS) > 1 {
					Password = itemS[1]
				}
				if len(itemS) > 2 {
					for i := 2; i < len(itemS); i++ {
						Remarks = strings.Join([]string{Remarks, itemS[i], itemSplit}, "")
					}
				}
				Remarks = strings.TrimRight(Remarks, itemSplit)
				others = append(others, database.Others{
					ProjectsID: uint(ProjectsID),
					Account:    Account,
					Password:   Password,
					NewStatus:  StatusInt,
					Exptime:    0,
					Remarks:    Remarks,
				})
			}
		}
	}
	others = RemoveRepeatedSingle(others)
	if Repeated {
		accList, err := database.GetOthersList(projectsID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "获取所有帐号失败，关闭过滤重复功能再试",
			})
			return
		}
		others = IgnoreRepeated(others, accList)
	}
	batchLen := len(others)
	// fmt.Println(batchLen)
	if batchLen > 1000 {
		database.OthersInBatches(others)
	} else {
		if batchLen > 0 {
			database.OthersBatches(others)
		}
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

// RemoveRepeatedSingle Remove Repeated Element
func RemoveRepeatedSingle(personList []database.Others) (result []database.Others) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if personList[i].Account == personList[j].Account {
				repeat = true
				break
			}
		}
		if !repeat {
			result = append(result, personList[i])
		}
	}
	return
}

func IgnoreRepeated(postList, dataList []database.Others) []database.Others {
	if len(dataList) != 0 {
		var temp []database.Others
		for _, item := range postList {
			exist := false
			for _, ig := range dataList {
				if item.Account == ig.Account {
					exist = true
				}
			}
			if !exist {
				temp = append(temp, item)
			}
		}
		// fmt.Println(temp)
		return temp
	}
	return postList
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
