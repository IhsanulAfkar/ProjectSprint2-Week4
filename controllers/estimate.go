package controllers

import (
	"Week4/db"
	"Week4/forms"
	"Week4/helper"
	"Week4/models"
	"database/sql"
	"fmt"
	"math"

	"github.com/gin-gonic/gin"
)

type EstimateController struct{}

func (h EstimateController) EstimatePrice(c *gin.Context){
	userId := c.MustGet("userId")
	fmt.Println("userId", userId)
	var estimateForm forms.EstimatePrice
	if err := c.ShouldBindJSON(&estimateForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if len(estimateForm.Orders) <= 0 {
		c.JSON(400, gin.H{"message":"order cannot be empty"})
		return
	}
	var startingMerchant models.Merchant
	listMerchants := make([]models.Merchant,0)
	totalPrice := 0
	conn := db.CreateConn()
	count := 0
	for _, order := range estimateForm.Orders {
	
		if len(order.Items) == 0 {
			c.JSON(400, gin.H{"message":"item cannot be empty"})
			return
		}
		// check merchant
		var merchant models.Merchant
		err := conn.QueryRowx("SELECT * FROM merchant WHERE id::text = $1", order.MerchantId).StructScan(&merchant)
		if err != nil {
			if err ==sql.ErrNoRows{
				c.JSON(404, gin.H{"message":"merchant not found"})
				return
			}
			fmt.Println(err.Error())
			c.JSON(500,"server error")
			return
		}
		if order.IsStartingPoint {
			if count != 0 {
				c.JSON(400, gin.H{"message":"invalid starting point"})
				return
			}
			count = 1
			startingMerchant = merchant
		} else {
			listMerchants = append(listMerchants, merchant)
		}
		// count item price
		for _,item := range order.Items{
			// check item
			var product models.ItemPrice
			err := conn.QueryRowx("SELECT id, price FROM item WHERE id::text = $1 AND \"merchantId\" = $2", item.ItemId, merchant.Id).StructScan(&product)
			fmt.Println(item.ItemId, merchant.Id)
			if err != nil {
				if err ==sql.ErrNoRows{
					c.JSON(404, gin.H{"message":"item not found"})
					return
				}
				fmt.Println(err.Error())
				c.JSON(500,"server error")
				return
			}
			totalPrice += item.Quantity * product.Price
		}
	}
	totalDistance := 0.0
	// Time to TSP!!!
	n := len(listMerchants)
	isVisited := make([]bool,n)
	tour := make([]models.Merchant,0)
	currentMerchant := startingMerchant
	// make first tour from isStartingPoint
	tour = append(tour, currentMerchant)
	for len(tour) <= n {
		next := -1
		minDist := math.Inf(1)

		for i := 0; i < n; i++ {
			if !isVisited[i]{
				// fmt.Println(listMerchants[i])
				dist := helper.CountHaversine(currentMerchant.Lat, currentMerchant.Lon, listMerchants[i].Lat,listMerchants[i].Lon)
				if dist <= minDist {
					minDist = dist
					next = i
					fmt.Println("minDist: ",minDist, i)
				}
			}
		}
		if next == -1 {
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		// mark as visited and continue from the nearest merchant
		isVisited[next] = true
		tour = append(tour, listMerchants[next])
		totalDistance += minDist
		currentMerchant = listMerchants[next]
		fmt.Println("total :", totalDistance)
	}
	lastMerchant := tour[len(tour)-1]
	totalDistance += helper.CountHaversine(lastMerchant.Lat, lastMerchant.Lon,estimateForm.UserLocation.Lat, estimateForm.UserLocation.Long)
	fmt.Println("total :", totalDistance)
	// count delivery in minutes
	deliveryTime := math.Round(totalDistance / 40 * 60)
	// insert into database
	query := "INSERT INTO estimate (\"userId\",\"userLat\", \"userLon\", \"totalPrice\", \"estimateDeliveryTime\") VALUES ($1,$2,$3,$4,$5) RETURNING id"
	var estimateId string
	err := conn.QueryRow(query, userId, estimateForm.UserLocation.Lat, estimateForm.UserLocation.Long, totalPrice, deliveryTime).Scan(&estimateId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	for _, order := range estimateForm.Orders {
		query = "INSERT INTO \"estimateOrder\" (\"estimateId\", \"isStarting\", \"merchantId\") VALUES ($1,$2,$3) RETURNING id"
		var estimateOrderId string
		err := conn.QueryRow(query, estimateId, order.IsStartingPoint, order.MerchantId).Scan(&estimateOrderId)
		if err != nil{
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		for _, item := range order.Items {
			query = "INSERT INTO \"estimateOrderItem\" (\"estimateOrderId\", \"itemId\", quantity) VALUES ($1,$2,$3) RETURNING id"
			var estimateOrderItemId string
			err := conn.QueryRow(query, estimateOrderId, item.ItemId, item.Quantity).Scan(&estimateOrderItemId)
			if err != nil{
				fmt.Println(err.Error())
				c.JSON(500, gin.H{"message":"server error"})
				return
			}
		}
	}
	c.JSON(200, gin.H{"estimatedDeliveryTimeInMinutes": deliveryTime, "totalPrice": totalPrice, "calculatedEstimateId":estimateId})
}
