package controllers

import (
	"Week4/db"
	"Week4/forms"
	"Week4/helper"
	"Week4/helper/validator"
	"Week4/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ItemController struct {
}

func (h ItemController) CreateItem(c *gin.Context){
	var itemForm forms.CreateItem
	if err := c.ShouldBindJSON(&itemForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	merchantId := c.Param("merchantId")
	if !validator.StringCheck(itemForm.Name, 2, 30){
		c.JSON(400, gin.H{
			"message":"invalid name"})
		return
	}
	if !helper.Includes(itemForm.ProductCategory,models.ItemCategory[:] ){
		c.JSON(400, gin.H{
			"message":"invalid category"})
		return
	}
	if itemForm.Price < 1 {
		c.JSON(400, gin.H{
			"message":"invalid price"})
		return
	}
	if !validator.IsURL(itemForm.ImageUrl){
		c.JSON(400, gin.H{"message":"invalid url"})
		return
	}
	// check merchant
	conn := db.CreateConn()
	query := "SELECT COUNT(1) FROM merchant WHERE id::text = $1 LIMIT 1"
	res, err := conn.Exec(query, merchantId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0{
		c.JSON(404, gin.H{"message":"merchant not found"})
		return
	}
	query = "INSERT INTO item (name, \"productCategory\", \"imageUrl\", price, \"merchantId\") VALUES ($1,$2,$3,$4,$5) RETURNING id"
	var itemId string
	err = conn.QueryRow(query, itemForm.Name, itemForm.ProductCategory, itemForm.ImageUrl, itemForm.Price, merchantId).Scan(&itemId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(201, gin.H{"itemId":itemId})
}