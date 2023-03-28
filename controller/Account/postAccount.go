package controller

import (
	"bytes"
	"colaAPI/database"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
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
	HasMore  string `form:"hasmore" json:"hasmore" xml:"hasmore" binding:"required"`
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
	hasMore := false
	if form.HasMore == "true" {
		hasMore = true
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

	projectsID := GetProjectsID(c)
	ProjectsID, _ := strconv.ParseInt(projectsID, 10, 64)
	Projects, err := database.ProjectsCheckID(ProjectsID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	var statusJson []*StatusJSON
	json.Unmarshal([]byte(Projects.StatusJSON), &statusJson)

	var hasPower bool = false

	for _, item := range statusJson {
		if item.Status == form.Status {
			if item.Import {
				hasPower = true
				break
			}
		}
	}

	if !hasPower {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "状态不支持导入",
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

	var FirstLen string = data[0]
	dataLen = dataLen - 1
	checkFirstLen := len(data[0])
	if checkFirstLen == 0 {
		hasContent := RandInt(1, dataLen)
		FirstLen = data[hasContent]
	}

	IdIsFirst := true

	if form.Status == "8" {
		checkID := strings.Split(FirstLen, itemSplit)[0]
		if !IsIdCard(checkID) {
			IdIsFirst = false
		}
	}

	var account []database.Accounts

	hasPhone, index := isPhone(FirstLen, itemSplit)

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
					Remarks       string = ""
				)
				if len(itemS) > 1 {
					Password = itemS[1]
				}
				if hasPhone {
					PhoneNumber = itemS[index]
					PhonePassword = itemS[index+1]
				}
				if form.Status == "8" {
					if !IdIsFirst {
						UserName = itemS[1]
						Password = itemS[0]
					}
				}
				if hasMore {
					if len(itemS) < 8 {
						c.JSON(http.StatusOK, gin.H{
							"status":  1,
							"message": "数据格式错误",
						})
						return
					}
					var gold int64
					if strings.Contains(itemS[2], "亿") {
						gx := strings.Split(itemS[2], "亿")
						goldstr := gx[0]
						if strings.Contains(itemS[2], ".") {
							g := strings.Split(goldstr, ".")
							goldstr = strings.Join([]string{g[0], g[1]}, "")
							var x int64 = 10000000
							if len(g[1]) >= 2 {
								x = 1000000
							}
							n, _ := strconv.ParseInt(goldstr, 10, 64)
							gold = n * x
						} else {
							n, _ := strconv.ParseInt(goldstr, 10, 64)
							gold = n * 100000000
						}
					} else if strings.Contains(itemS[2], "万") {
						gx := strings.Split(itemS[2], "万")
						goldstr := gx[0]
						n, _ := strconv.ParseInt(goldstr, 10, 64)
						gold = n * 10000
					} else {
						n, _ := strconv.ParseInt(itemS[2], 10, 64)
						gold = n
					}
					TodayGold = gold
					Multiple, _ = strconv.ParseInt(itemS[3], 10, 64)
					Diamond, _ = strconv.Atoi(itemS[4])
					Crazy, _ = strconv.Atoi(itemS[5])
					Precise, _ = strconv.Atoi(itemS[6])
					Cold, _ = strconv.Atoi(itemS[7])
					if len(itemS) > 8 {
						for i := 7; i < len(itemS); i++ {
							Remarks += strings.Join([]string{itemS[i], itemSplit}, "")
						}
					}
				} else {
					if len(itemS) > 2 {
						for i := 2; i < len(itemS); i++ {
							Remarks = strings.Join([]string{Remarks, itemS[i], itemSplit}, "")
						}
					}
				}
				Remarks = strings.TrimRight(Remarks, itemSplit)
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
					Remarks:       Remarks,
				})
			}
		}
	}
	account = RemoveRepeatedSingle(account, hasPhone)
	if Repeated {
		var (
			hasStatus []string
		)
		for _, item := range statusJson {
			if item.Pull {
				hasStatus = append(hasStatus, item.Status)
			}
		}
		if !in(form.Status, hasStatus) {
			hasStatus = append(hasStatus, form.Status)
		}
		accList, err := database.GetAccountListUseIn(projectsID, hasStatus)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status":  1,
				"message": "获取所有帐号失败，关闭过滤重复功能再试",
			})
			return
		}
		account = IgnoreRepeated(account, accList, hasPhone)
	}
	batchLen := len(account)
	// fmt.Println(batchLen)
	if batchLen > 1000 {
		database.AccountInBatches(account)
	} else {
		if batchLen > 0 {
			database.AccountBatches(account)
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
func RemoveRepeatedSingle(personList []database.Accounts, hasPhone bool) (result []database.Accounts) {
	n := len(personList)
	for i := 0; i < n; i++ {
		repeat := false
		for j := i + 1; j < n; j++ {
			if hasPhone {
				if personList[i].PhoneNumber == personList[j].PhoneNumber {
					repeat = true
					break
				}
			} else {
				if personList[i].UserName == personList[j].UserName {
					repeat = true
					break
				}
			}
		}
		if !repeat {
			result = append(result, personList[i])
		}
	}
	return
}

func IgnoreRepeated(postList, dataList []database.Accounts, hasPhone bool) []database.Accounts {
	if len(dataList) != 0 {
		var temp []database.Accounts
		for _, item := range postList {
			exist := false
			for _, ig := range dataList {
				if hasPhone {
					if item.PhoneNumber == ig.PhoneNumber {
						exist = true
					}
				} else {
					if item.UserName == ig.UserName {
						exist = true
					}
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

func RandInt(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min+1) + min
}

func IsIdCard(idCard string) (res bool) {
	res, _ = regexp.Match("^[1-9]\\d{7}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}$|^[1-9]\\d{5}[1-9]\\d{3}((0\\d)|(1[0-2]))(([0|1|2]\\d)|3[0-1])\\d{3}([0-9]|X)$", []byte(idCard))
	return
}
func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
