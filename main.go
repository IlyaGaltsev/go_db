package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name     string
	Temp     uint
	AnalysID uint
}

type Analysis struct {
	gorm.Model
	Name    string
	Cost    uint
	Price   uint
	Groups  []Group `gorm:"foreignKey:AnalysID"`
	OrderID uint
}

type Orders struct {
	gorm.Model
	AnalysID Analysis `gorm:"foreignKey:OrderID"`
}

func PrintAll(db *gorm.DB) {
	var orders []Orders
	db.Model(&Orders{}).Preload("AnalysID").Preload("AnalysID.Groups").Find(&orders)

	for _, order := range orders {
		jsonOutput, err := json.MarshalIndent(order, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonOutput))
	}
}

func PrintWhere(db *gorm.DB, condition func(order Orders) bool) {
	var orders []Orders
	db.Model(&Orders{}).Preload("AnalysID").Preload("AnalysID.Groups").Find(&orders)

	for _, order := range orders {
		if condition(order) {
			jsonOutput, err := json.MarshalIndent(order, "", " ")

			if err != nil {
				panic(err)
			}

			fmt.Print((string(jsonOutput)))
		}

	}
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("aaaaa pomogiiiiiitee, failed to connect database")
	}

	db.AutoMigrate(&Group{})
	db.AutoMigrate(&Analysis{})
	db.AutoMigrate(&Orders{})

	group1 := Group{Name: "Group A", Temp: 25}
	group2 := Group{Name: "Group B", Temp: 30}
	group3 := Group{Name: "Grou33333", Temp: 12}
	group4 := Group{Name: "Grou44444", Temp: 0}
	db.Create(&group1)
	db.Create(&group2)
	db.Create(&group3)
	db.Create(&group4)

	analysis1 := Analysis{Name: "Analysis 1", Cost: 100, Price: 150, Groups: []Group{group1, group2}}
	analysis2 := Analysis{Name: "Analysis 2", Cost: 120, Price: 180, Groups: []Group{group3, group4}}
	db.Create(&analysis1)
	db.Create(&analysis2)

	order1 := Orders{AnalysID: analysis1}
	order2 := Orders{AnalysID: analysis2}
	db.Create(&order1)
	db.Create(&order2)

	PrintAll(db)
	PrintWhere(db, func(order Orders) bool {
		return order.CreatedAt.After(time.Date(2024, 2, 3, 0, 0, 0, 0, time.UTC))
	})
}
