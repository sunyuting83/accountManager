package main

import (
	orm "colaAPI/UpdatesAccount/database"
	"fmt"
)

func main() {
	orm.InitDB()
	err := orm.UpdateAccountGamesID()
	if err != nil {
		fmt.Println(err)
	}
	// c, _ := orm.GetCount()
	// fmt.Println(c)
	defer orm.Eloquent.Close()
}
