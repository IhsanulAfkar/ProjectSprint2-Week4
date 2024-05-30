package routes

import (
	"Week4/controllers"
	"Week4/middleware"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	adminController := new(controllers.AdminController)
	userController := new(controllers.UserController)
	merchantController := new(controllers.MerchantController)
	itemController := new(controllers.ItemController)
	estimateController := new(controllers.EstimateController)
	mediaController := new(controllers.MediaController)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	admin := router.Group("/admin")
	{
		admin.POST("/register", adminController.Register)
		admin.POST("/login", adminController.Login)
		merchants := admin.Group("/merchants")
		{
			merchants.Use(middleware.AdminAuthMiddleware)
			merchants.POST("/", merchantController.Create)
			merchants.GET("/", merchantController.GetAllMerchant)
			merchants.POST("/:merchantId/items", itemController.CreateItem)
		}
	}
	users := router.Group("/users")
	{
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)
	}
	// image
	image := router.Group("/")
	{
		image.Use(middleware.AdminAuthMiddleware)
		image.POST("/image", mediaController.UploadImage)
	}
	// users auth
	router.Use(middleware.UserAuthMiddleware)
	{
		router.GET("/merchants/nearby/:coordinate", merchantController.GetNearbyMerchant)
		router.POST("/users/estimate", estimateController.EstimatePrice)
		router.POST("/users/orders", userController.Order)
		router.GET("/users/orders", userController.OrderHistory)
	}
	

	return router
}