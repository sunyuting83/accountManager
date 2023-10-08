package controller

import (
	"colaAPI/UsersApi/database"
	"colaAPI/UsersApi/utils"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type accArray []*database.Accounts

func (x accArray) Len() int {
	return len(x)
}
func (x accArray) Less(i, j int) bool {
	return x[i].TodayGold > x[j].TodayGold
}
func (x accArray) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func ExportAccountDrawed(c *gin.Context) {

	var (
		date     string = c.Query("date")
		multiple string = c.DefaultQuery("multiple", "true")
		diamond  string = c.DefaultQuery("diamond", "false")
		crazy    string = c.DefaultQuery("crazy", "false")
		cold     string = c.DefaultQuery("cold", "false")
		precise  string = c.DefaultQuery("precise", "false")
		remarks  string = c.DefaultQuery("remarks", "false")
		excel    string = c.DefaultQuery("excel", "false")

		Multiple   bool = false
		Diamond    bool = false
		Crazy      bool = false
		Cold       bool = false
		Precise    bool = false
		Remarks    bool = false
		Excel      bool = false
		Returnfile []byte
		RetrnExcel *excelize.File
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
	if excel == "true" {
		Excel = true
	}
	fmt.Println(excel)
	fmt.Println(Excel)
	projectsID := GetProjectsID(c)
	startTime, endTime := utils.GetSqlDateTime(date)

	data, err := database.ExportAccountDrawed(projectsID, startTime, endTime)
	var arr accArray = data
	sort.Sort(arr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  1,
			"message": "获取数据失败",
		})
		return
	}
	cType := "text/plain"
	fileName := fmt.Sprintf("%s%s%s.xlsx", time.Now().Format("2016-01-02"), `-`, "user")
	if Excel {

		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
		c.Writer.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		RetrnExcel = MakeDrawedExcelFile(arr, Multiple, Diamond, Crazy, Cold, Precise, Remarks, Excel)
	} else {
		c.Header("Content-Type", cType)
		Returnfile = MakeDrawedExportFile(arr, Multiple, Diamond, Crazy, Cold, Precise, Remarks, Excel)
	}

	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")
	// c.Writer.WriteString(string(file)) return image
	if Excel {
		// _ = RetrnExcel.Write(c.Writer)
		rt, _ := RetrnExcel.WriteToBuffer()
		_, _ = c.Writer.Write(rt.Bytes())
		return
	}
	c.Data(200, cType, Returnfile)
}

func MakeDrawedExcelFile(arrData accArray, Multiple, Diamond, Crazy, Cold, Precise, Remarks, Excel bool) (f *excelize.File) {
	var (
		TimeStr string = time.Now().Format("2006-01-02")
	)
	f = excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet(TimeStr)
	if err != nil {
		fmt.Println(err)
	}
	// Set value of a cell.
	f.SetCellValue(TimeStr, "A1", "帐号")
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	com := "A"
	for index, item := range arrData {
		com = "A"
		index = index + 2
		tableLin := strconv.Itoa(index)

		f.SetCellValue(TimeStr, com+tableLin, item.UserName)

		if len(item.Password) != 0 {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "密码")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.Password)
		}
		Gold := MakeGoldString(item.TodayGold, Excel)

		com = GetCom(com)
		if index == 2 {
			TableName := com + "1"
			f.SetCellValue(TimeStr, TableName, "金币")
		}
		Goldint, _ := strconv.Atoi(Gold)
		f.SetCellValue(TimeStr, com+tableLin, Goldint)

		if Multiple {
			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "炮台")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.Multiple)

		}
		if Diamond {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "钻石")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.Diamond)

		}
		if Crazy {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "狂暴")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.Crazy)

		}
		if Precise {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "瞄准")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.Precise)

		}
		if Cold {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "冰冻")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.Cold)

		}
		if len(item.PhoneNumber) != 0 {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "手机号")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.PhoneNumber)

		}
		if len(item.PhonePassword) != 0 {

			com = GetCom(com)
			if index == 2 {
				TableName := com + "1"
				f.SetCellValue(TimeStr, TableName, "手机密码")
			}
			f.SetCellValue(TimeStr, com+tableLin, item.PhonePassword)

		}
		if Remarks {
			if len(item.Remarks) != 0 {

				com = GetCom(com)
				if index == 2 {
					TableName := com + "1"
					f.SetCellValue(TimeStr, TableName, "其他")
				}
				f.SetCellValue(TimeStr, com+tableLin, item.Remarks)

			}
		}
	}
	return f
}

func MakeDrawedExportFile(data accArray, Multiple, Diamond, Crazy, Cold, Precise, Remarks, Excel bool) []byte {
	var (
		temp   []string
		rtData []byte
	)

	for _, item := range data {
		itemStr := strings.Join([]string{item.UserName}, "")

		if len(item.Password) != 0 {
			itemStr = strings.TrimRight(itemStr, "\t")
			itemStr = strings.Join([]string{itemStr, item.Password}, "\t")
		}
		itemStr = strings.TrimRight(itemStr, "\t")
		Gold := MakeGoldString(item.TodayGold, Excel)
		itemStr = strings.Join([]string{itemStr, Gold}, "\t")
		if Multiple {
			itemStr = strings.TrimRight(itemStr, "\t")
			MultipleStr := MakeGoldString(item.Multiple, Excel)
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
		if Remarks {
			if len(item.Remarks) != 0 {
				itemStr = strings.TrimRight(itemStr, "\t")
				itemStr = strings.Join([]string{itemStr, item.Remarks}, "\t")
			}
		}
		itemStr = strings.TrimRight(itemStr, "\t")
		temp = append(temp, itemStr)
	}
	rtData = []byte(strings.Join(temp, "\r\n"))
	return rtData
}

func GetCom(currency string) (newstr string) {
	var comList []string = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M"}
	for i, item := range comList {
		if item == currency {
			newstr = comList[i+1]
			break
		}
	}
	return
}

func MakeGoldString(gold int64, Excel bool) (goldStr string) {
	if gold >= 100000000 {
		gold = gold / 100000000
		goldStr = strconv.FormatInt(gold, 10)
		if !Excel {
			goldStr = strings.Join([]string{goldStr, "亿"}, "")
		}
		return
	}
	if gold >= 10000 {
		gold = gold / 10000
		goldStr = strconv.FormatInt(gold, 10)
		if !Excel {
			goldStr = strings.Join([]string{goldStr, "万"}, "")
		}
	}
	return
}
