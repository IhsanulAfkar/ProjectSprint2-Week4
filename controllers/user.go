package controllers

import (
	"Week4/db"
	"Week4/forms"
	"Week4/helper"
	"Week4/helper/jwt"
	"Week4/helper/validator"
	"Week4/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct{}

func (h UserController)Register(c *gin.Context){
	var userRegister forms.AdminRegister
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if !validator.StringCheck(userRegister.Username, 5,30){
		c.JSON(400, gin.H{
			"message":"invalid username"})
		return
	}
	if !validator.StringCheck(userRegister.Password, 5,30){
		c.JSON(400, gin.H{
			"message":"invalid password"})
		return
	}
	if !validator.IsEmail(userRegister.Email){
		c.JSON(400, gin.H{
			"message":"invalid email"})
		return
	}
	conn := db.CreateConn()
	query := "SELECT EXISTS (SELECT 1 FROM public.user WHERE username = $1 OR email = $2 UNION SELECT 1 FROM admin WHERE username = $3) AS is_exist"
	var isExist bool
	err := conn.QueryRow(query, userRegister.Username,userRegister.Email, userRegister.Username).Scan(&isExist)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if isExist {
		c.JSON(409, gin.H{"message":"username or email already exists"})
		return
	}
	
	query = "INSERT INTO public.user (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	var userId string
	err = conn.QueryRow(query,userRegister.Username,userRegister.Email, userRegister.Password).Scan(&userId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	accessToken := jwt.SignJWT(userId,"user")
	c.JSON(201, gin.H{"token":accessToken})
}


func (h UserController) Login(c *gin.Context){
	var userLogin forms.AdminLogin
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
	}
	conn:=db.CreateConn()
	var user models.User
	query:= "SELECT * FROM public.user WHERE username = $1"
	err := conn.QueryRowx(query, userLogin.Username).StructScan(&user)
	if err != nil {
		if err ==sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"no user found"})
			return
		}
		fmt.Println(err.Error())
		return
	}
	if userLogin.Password != user.Password {
		c.JSON(400, gin.H{"message":"invalid password"})
		return
	}
	accessToken := jwt.SignJWT(user.Id,"user")
	c.JSON(200, gin.H{"token":accessToken})
}

func (h UserController)Order(c *gin.Context){
	userId := c.MustGet("userId")
	// lazy approach
	var orderForm struct {
		CalculatedEstimateId string `json:"calculatedEstimateId"`
	}
	if err := c.ShouldBindJSON(&orderForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if orderForm.CalculatedEstimateId == ""{
		c.JSON(400, gin.H{
			"message":"empty estimateId"})
		return
	}
	conn := db.CreateConn()
	// var estimateId string
	query := "UPDATE estimate SET \"orderId\" = $1 WHERE \"userId\" = $2 AND id = $3 AND \"orderId\" IS NULL"
	orderId := uuid.New()
	res, err := conn.Exec(query, orderId, userId, orderForm.CalculatedEstimateId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if rows, _:=res.RowsAffected(); rows == 0{
		c.JSON(404, gin.H{"message":"estimate not found"})
		return
	}
	c.JSON(201, gin.H{"orderId":orderId})
}

func (h UserController)OrderHistory(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	merchants := make([]models.NearbyMerchant, 0)
	merchantId := c.Query("merchantId")
	merchantCategory := c.Query("merchantCategory")
	name := strings.ToLower(c.Query("name"))
	if !helper.Includes(merchantCategory, models.MerchantCategory[:]) && merchantCategory != "" {
		c.JSON(200, gin.H{"data":merchants})
		return
	}
	baseQuery := `SELECT e.id
	FROM estimate e
		LEFT JOIN "estimateOrder" eo ON e.id = eo."estimateId"
		LEFT JOIN "estimateOrderItem" eoi ON eo.id = eoi."estimateOrderId"
		LEFT JOIN item i ON i.id = eoi."itemId"
		LEFT JOIN merchant m ON m.id = eo."merchantId"
	WHERE e."orderId" IS NOT NULL `
	var args []interface{}
	var queryParams []string
	argIdx := 1
	if name != ""{
		nameWildcard := "%" + name +"%"
		// m.name ILIKE '%' || 'its' || '%' OR
		// p.name ILIKE '%' || 'its' || '%'
		queryParams = append(queryParams," (m.name ILIKE $"+strconv.Itoa(argIdx) +" OR i.name ILIKE $"+strconv.Itoa(argIdx+1) +") ") 
		args = append(args, nameWildcard)
		args = append(args, nameWildcard)
		argIdx += 2
	}
	
	if merchantId != ""{
		queryParams = append(queryParams, " m.id::text = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, merchantId)
		argIdx += 1
	}
	if merchantCategory != ""{
		queryParams = append(queryParams, " m.\"merchantCategory\" = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, merchantCategory)
		argIdx += 1
	}
	if len(queryParams) > 0 {
		allQuery := strings.Join(queryParams, " AND")
		baseQuery +=" AND "+ allQuery
	}
	baseQuery += "GROUP BY e.id LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	conn := db.CreateConn()
	var orderList []string
	err := conn.Select(&orderList, baseQuery, args...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var orders []models.Order
	for _, orderId := range orderList {
		query := `SELECT 
		m.id as "merchantId",
		m.name as "name",
		m."merchantCategory" as "merchantCategory",
		m."imageUrl" as "imageUrl",
		m.lat as lat,
		m.lon as lon,
		m."createdAt" as "createdAt",
		COALESCE(json_agg(json_build_object('itemId', i.id, 'name', i.name, 'productCategory', i."productCategory", 'price', i.price, 'imageUrl', 'quantity', eoi.quantity ,i."imageUrl", 'createdAt', i."createdAt")) FILTER (WHERE i.id IS NOT NULL), '[]')::text AS products  
		FROM "estimateOrder" e 
			LEFT JOIN "estimateOrderItem" eoi ON e.id = eoi."estimateOrderId"
			LEFT JOIN item i ON i.id = eoi."itemId"
			LEFT JOIN merchant m ON m.id = e."merchantId" 
		WHERE e."estimateId" = $1
		GROUP BY m.id`
		var getOrderDetail []models.GetOrderDetail
		var orderDetails []models.OrderDetail
		var order models.Order
		err := conn.Select(&getOrderDetail, query, orderId)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		for _, detail := range getOrderDetail {
			var orderDetail models.OrderDetail
			orderDetail.Merchant.MerchantId = detail.MerchantId
			orderDetail.Merchant.Name = detail.Name
			orderDetail.Merchant.MerchantCategory = detail.MerchantCategory
			orderDetail.Merchant.Location.Lat = detail.Lat
			orderDetail.Merchant.Location.Lon = detail.Lon
			orderDetail.Merchant.CreatedAt = detail.CreatedAt
			var parsedItems []models.GetItemDetail
			err = json.Unmarshal([]byte(detail.Products), &parsedItems)
			if err != nil {
				fmt.Println("json")
				fmt.Println(err.Error())
				c.JSON(500, gin.H{"message":"server error"})
				return
			}
			orderDetail.Items = parsedItems
			orderDetails = append(orderDetails, orderDetail)
		}
		order.OrderId = orderId
		order.Orders = orderDetails
		orders = append(orders, order)

	}
	c.JSON(200, orders)
}