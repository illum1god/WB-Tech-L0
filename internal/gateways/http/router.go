package http

import (
	"WB-L0/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func setupRouter(r *gin.Engine, service usecase.Service, logger *logrus.Logger) {
	handler := NewHandler(service)

	r.Use(ErrorRecoveryMiddleware(logger))

	api := r.Group("/api")
	{
		api.POST("/order", handler.SaveOrder)
		api.GET("/order", handler.GetOrders)
		api.GET("/order/:uid", handler.GetOrderByUID)
	}

	r.LoadHTMLGlob("web/*.html")

	r.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", gin.H{})
	})

	r.GET("/order/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		order, err := service.GetOrderByUID(c, uid)
		if err != nil {
			logger.Errorf("Ошибка получения заказа: %v", err)
			c.HTML(http.StatusInternalServerError, "order.html", gin.H{
				"Error": "Не удалось получить заказ",
			})
			return
		}
		c.HTML(http.StatusOK, "order.html", gin.H{
			"Order": order,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/home")
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
