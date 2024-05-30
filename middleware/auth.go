package middleware

import (
	"Week4/db"
	"Week4/helper/jwt"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
func AdminAuthMiddleware(c *gin.Context) {
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	admin, err := jwt.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	if admin.Role != "admin" {
		c.AbortWithStatusJSON(401, gin.H{"message": "access not allowed"})
		return
	}
	// find user
	conn := db.CreateConn()
	res, err := conn.Exec("SELECT 1 FROM admin WHERE id = $1 LIMIT 1", admin.Id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"message": "server error"})
		return
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "admin not found"})
		return
	}
	c.Set("adminId", admin.Id)
	c.Next()
}

func UserAuthMiddleware(c *gin.Context) {
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	user, err := jwt.ParseToken(token)
	if err != nil {

		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}

	// find user
	conn := db.CreateConn()
	var isExists bool
	fmt.Println(user)
	err = conn.QueryRow("SELECT EXISTS (SELECT 1 FROM public.user WHERE id = $1) AS is_exists", user.Id).Scan(&isExists)
	fmt.Println(isExists)
	if err != nil || !isExists {
		// fmt.Println(err.Error())
		c.AbortWithStatusJSON(401, gin.H{"message": "invalid token"})
		return
	}

	c.Set("userId", user.Id)
	c.Next()
}