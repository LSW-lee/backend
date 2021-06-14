package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"backend/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	AddProducts()
	UpdatePromos()
}

func UpdatePromos() {
	jsonfile, _ := ioutil.ReadFile("./promos.json")
	var products []models.Product
	json.Unmarshal(jsonfile, &products)

	db, _ := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1)/gomusic?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	for _, product := range products {
		fmt.Println(product)
		//db.Table("customers").Where("id=?", id).Update("loggedin", 0).Error
		err := db.Table("products").Where("id=?", product.ID).Update("promotion", product.Price).Error
		if err != nil {
			panic("INSERT ERROR")
		}
	}
}

func AddProducts() {
	jsonfile, _ := ioutil.ReadFile("./cards.json")
	var products []models.Product
	json.Unmarshal(jsonfile, &products)

	db, _ := gorm.Open(mysql.Open("root:1234@tcp(127.0.0.1)/gomusic?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	for _, product := range products {
		fmt.Println(product)
		err := db.Create(&product).Error
		if product.Promotion == 0 {
			db.Model(&product).Select("promotion").Updates(map[string]interface{}{"promotion": nil})
		}
		if err != nil {
			panic("INSERT ERROR")
		}
	}
}
