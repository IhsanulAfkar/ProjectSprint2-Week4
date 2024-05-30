package controllers

import (
	"Week4/db"
	"Week4/forms"
	"Week4/helper"
	"Week4/helper/validator"
	"Week4/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type MerchantController struct{}

func (h MerchantController) Create(c *gin.Context){
	var merchantForm forms.CreateMerchant
	if err := c.ShouldBindJSON(&merchantForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	fmt.Println(merchantForm)
	if !validator.StringCheck(merchantForm.Name, 2, 30){
		c.JSON(400, gin.H{
			"message":"invalid name"})
		return
	}
	if !helper.Includes(merchantForm.MerchantCategory,models.MerchantCategory[:] ){
		c.JSON(400, gin.H{
			"message":"invalid category"})
		return
	}
	if !validator.IsURL(merchantForm.ImageUrl){
		c.JSON(400, gin.H{
			"message":"invalid url"})
		return
	}
	conn := db.CreateConn()
	query := "INSERT INTO merchant (name, \"merchantCategory\", \"imageUrl\", lat, lon) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	var merchantId string
	err := conn.QueryRow(query,merchantForm.Name, merchantForm.MerchantCategory, merchantForm.ImageUrl, merchantForm.Location.Lat,merchantForm.Location.Lon).Scan(&merchantId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(201, gin.H{"merchantId":merchantId})
}

func (h MerchantController)GetAllMerchant(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	merchantId := c.Query("merchantId")
	merchantCategory := c.Query("merchantCategory")
	name := strings.ToLower(c.Query("name"))
	createdAt := c.Query("createdAt")
	if !helper.Includes(merchantCategory, models.MerchantCategory[:]){
		merchantCategory = ""
	}
	if createdAt != "asc" && createdAt != "desc"{
		createdAt = ""
	}
	baseQuery := "SELECT id, name, \"merchantCategory\", \"imageUrl\", lat, lon, \"createdAt\" FROM merchant"
	var args []interface{}
	var queryParams []string
	argIdx := 1
	if name != ""{
		nameWildcard := "%" + name +"%"
		queryParams = append(queryParams," name ILIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nameWildcard)
		argIdx += 1
	}
	if merchantId != ""{
		queryParams = append(queryParams, " id::text = $" + strconv.Itoa(argIdx)+ " ")
		args = append(args, merchantId)
		argIdx += 1
	}
	if merchantCategory != ""{
		queryParams = append(queryParams, " \"merchantCategory\" = $" + strconv.Itoa(argIdx)+ " ")
		args = append(args, merchantCategory)
		argIdx += 1
	}
	if len(queryParams) > 0 {
		allQuery := strings.Join(queryParams, " AND")
		baseQuery += " WHERE " + allQuery
	}
	baseQuery += " ORDER BY "
	if createdAt == "" {
		baseQuery += " \"createdAt\" DESC"
	} else {
		if createdAt == "asc"{
			baseQuery += " \"createdAt\" ASC"
		}
	}
	baseQuery +=  " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	conn := db.CreateConn()
	merchants := make([]models.GetMerchant, 0)
	// err := conn.Select(&merchants, baseQuery, args...)
	rows, err := conn.Query(baseQuery, args...)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var merchant models.GetMerchant
		err = rows.Scan(&merchant.MerchantId, &merchant.Name, &merchant.MerchantCategory, &merchant.ImageUrl, &merchant.Location.Lat, &merchant.Location.Lon, &merchant.CreatedAt)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		merchants = append(merchants, merchant)
	}
	var totalData int
	err = conn.QueryRow("SELECT COUNT(id)::int FROM merchant").Scan(&totalData)
	if err != nil {
		if err != sql.ErrNoRows{
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
	}
	c.JSON(200, gin.H{"data": merchants, "meta":gin.H{"limit":limit,"offset":offset, "total":totalData}})
}

func (h MerchantController)GetNearbyMerchant(c *gin.Context){
	coordinate := c.Param("coordinate")
	coor := strings.Split(coordinate,",")
	if len(coor) != 2 {
		c.JSON(400, gin.H{"message":"bad params"})
		return
	}
	lat, err := strconv.ParseFloat(coor[0], 64)
	if err != nil {
		c.JSON(400, gin.H{"message":"bad params"})
		return
	}
	lon, err := strconv.ParseFloat(coor[1], 64)

	if err != nil {
		c.JSON(400, gin.H{"message":"bad params"})
		return
	}
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
	name := c.Query("name")
	merchantCategory := c.Query("merchantCategory")
	if !helper.Includes(merchantCategory, models.MerchantCategory[:]) && merchantCategory != "" {
		c.JSON(200, gin.H{"data":merchants})
		return
	}
	// query := "SELECT id, name, \"merchantCategory\", \"imageUrl\", lat, lon, \"createdAt\", FROM merchant"
	// distQuery := "(6371 * acos(cos(radians(INPUTLAT)) * cos(radians(lat)) * cos(radians(lon) - radians(INPUTLON)) + sin(radians(INPUTLAT)) * sin(radians(lat)))) AS distance"
	baseQuery := "SELECT m.id as id, m.name as name, m.\"merchantCategory\" as \"merchantCategory\", m.\"imageUrl\", m.lat, m.lon,m.\"createdAt\" as \"createdAt\", COALESCE(json_agg(json_build_object('id', p.id, 'name', p.name, 'productCategory', p.\"productCategory\", 'price', p.price, 'imageUrl',p.\"imageUrl\", 'createdAt', p.\"createdAt\")) FILTER (WHERE p.id IS NOT NULL), '[]')::text AS products, (6371 * acos(cos(radians("+fmt.Sprintf("%.4f", lat)+")) * cos(radians(m.lat)) * cos(radians(m.lon) - radians("+fmt.Sprintf("%.4f", lon)+")) + sin(radians("+fmt.Sprintf("%.4f", lat)+")) * sin(radians(m.lat)))) AS distance FROM merchant m LEFT JOIN item p ON m.id = p.\"merchantId\""
	fmt.Println(baseQuery)
	var args []interface{}
	var queryParams []string
	argIdx := 1
	if name != ""{
		nameWildcard := "%" + name +"%"
		// m.name ILIKE '%' || 'its' || '%' OR
		// p.name ILIKE '%' || 'its' || '%'
		queryParams = append(queryParams," (m.name ILIKE $"+strconv.Itoa(argIdx) +" OR p.name ILIKE $"+strconv.Itoa(argIdx+1) +" ") 
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
		baseQuery += " WHERE " + allQuery
	}
	baseQuery += " GROUP BY m.id, m.name ORDER BY distance LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	conn := db.CreateConn()
	rows, err := conn.Query(baseQuery, args...)
	fmt.Println(baseQuery)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	defer rows.Close()
	for rows.Next(){ 
		var merchant models.DistanceMerchant
		var items string
		var distance float64
		err = rows.Scan(&merchant.MerchantId, &merchant.Name, &merchant.MerchantCategory, &merchant.ImageUrl, &merchant.Location.Lat, &merchant.Location.Lon, &merchant.CreatedAt, &items, &distance)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		fmt.Println(items)
		// unmarshal
		var jsonItems []models.GetItem
		err = json.Unmarshal([]byte(items), &jsonItems)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		var nearbyMerch models.NearbyMerchant
		nearbyMerch.Merchant = merchant
		nearbyMerch.Items = jsonItems
		merchants = append(merchants, nearbyMerch)
		// count distance
		// distance := helper.CountHaversine(lat, lon, float64(merchant.Location.Lat), float64(merchant.Location.Lon))
		fmt.Println(items)
		fmt.Println(merchant)
		fmt.Println(distance)
	}
	// c.JSON(200, gin.H{"data":merchants})

	// fmt.Println(baseQuery)
	c.JSON(200, gin.H{"data":merchants})
}
// coret
// SELECT 
// m.id,
// m.name,
// m."merchantCategory",
// m."imageUrl",
// m.lat,
// m.lon,
// COALESCE(json_agg(json_build_object('id', p.id, 'name', p.name)) FILTER (WHERE p.id IS NOT NULL), '[]') AS products
// FROM 
// merchant m
// LEFT JOIN 
// item p ON m.id = p."merchantId"
// WHERE 
// m.name ILIKE '%' || 'Handmade' || '%' OR
// p.name ILIKE '%' || 'Handmade' || '%'
// GROUP BY 
// m.id, m.name;