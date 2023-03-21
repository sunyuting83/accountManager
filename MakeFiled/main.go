package main

import (
	orm "colaAPI/MakeFiled/database"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func main() {
	orm.InitDB()
	checkDate := orm.GetDateTime()
	checkDate = checkDate * 1000
	checkDateStr := strconv.FormatInt(checkDate, 10)
	ProList, err := orm.GetDateHasProject("pgsql", checkDateStr)
	if err != nil {
		fmt.Println("what?")
		return
	}
	if len(ProList) == 0 {
		fmt.Println("no projects")
		return
	}
	// var account []*orm.Accounts
	for _, item := range ProList {
		projectIDStr := strconv.FormatInt(item, 10)
		dateList, err := orm.GetDateTimeDataDraw(projectIDStr, "pgsql", checkDateStr)
		if err != nil {
			fmt.Println("what?")
			return
		}
		if len(dateList) == 0 {
			fmt.Println("no projects")
			return
		}
		for _, date := range dateList {
			var deleteID []int
			starTime, endTime := orm.GetSqlDateTime(date)
			data, err := orm.GetDateInData(projectIDStr, starTime, endTime)
			for _, Account := range data {
				deleteID = append(deleteID, int(Account.ID))
			}
			if err != nil {
				fmt.Println("what?")
				return
			}
			buff, _ := json.Marshal(&data)
			filed := &orm.Filed{
				ProjectsID: uint(item),
				FiledName:  date,
				Data:       string(buff),
			}
			filed.AddFiled()
			orm.DeleteAtDate(deleteID)
		}

		time.Sleep(1000)
	}

	defer orm.Eloquent.Close()
}
