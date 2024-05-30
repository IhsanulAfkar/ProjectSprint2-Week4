package controllers

import (
	"Week4/db"
	"Week4/forms"
	"Week4/helper/jwt"
	"Week4/helper/validator"
	"Week4/models"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func (h AdminController) Register(c *gin.Context){
	var adminRegister forms.AdminRegister
	if err := c.ShouldBindJSON(&adminRegister); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if !validator.StringCheck(adminRegister.Username, 5,30){
		c.JSON(400, gin.H{
			"message":"invalid username"})
		return
	}
	if !validator.StringCheck(adminRegister.Password, 5,30){
		c.JSON(400, gin.H{
			"message":"invalid password"})
		return
	}
	if !validator.IsEmail(adminRegister.Email){
		c.JSON(400, gin.H{
			"message":"invalid email"})
		return
	}
	conn := db.CreateConn()
	query := "SELECT EXISTS (SELECT 1 FROM admin WHERE username = $1 OR email = $2 UNION SELECT 1 FROM public.user WHERE username = $3) AS is_exist"
	var isExist bool
	err := conn.QueryRow(query, adminRegister.Username,adminRegister.Email, adminRegister.Username).Scan(&isExist)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if isExist {
		c.JSON(409, gin.H{"message":"username or email already exists"})
		return
	}
	// hashed_password, _ := hash.HashPassword(adminRegister.Password)
	query = "INSERT INTO admin (username, email, password) VALUES ($1, $2, $3) RETURNING id"
	var adminId string
	err = conn.QueryRow(query,adminRegister.Username,adminRegister.Email, adminRegister.Password).Scan(&adminId)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	accessToken := jwt.SignJWT(adminId,"admin")
	c.JSON(201, gin.H{"token":accessToken})
}

func (h AdminController) Login(c *gin.Context){
	var adminLogin forms.AdminLogin
	if err := c.ShouldBindJSON(&adminLogin); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
	}
	conn:=db.CreateConn()
	var admin models.Admin
	query:= "SELECT * FROM admin WHERE username = $1"
	err := conn.QueryRowx(query, adminLogin.Username).StructScan(&admin)
	if err != nil {
		if err ==sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"no admin found"})
			return
		}
		fmt.Println(err.Error())
		return
	}
	if adminLogin.Password != admin.Password {
		c.JSON(400, gin.H{"message":"invalid password"})
		return
	}
	accessToken := jwt.SignJWT(admin.Id,"admin")
	c.JSON(200, gin.H{"token":accessToken})
}
