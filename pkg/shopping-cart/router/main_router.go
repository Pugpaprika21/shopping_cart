package router

import (
	"github.com/Pugpaprika21/pkg/shopping-cart/handler"
	"github.com/Pugpaprika21/pkg/shopping-cart/server"
	"github.com/labstack/echo/v4"
)

func EchoRouter(e *echo.Echo, server *server.EchoServerEnvironment) {
	g := e.Group("/api")
	p := handler.NewProductHandler(server)
	g.GET("/products", p.GetAllProduct)
	g.POST("/products", p.SaveProduct)
	g.GET("/products/attachments/:attachment", p.GetProductAttachment)
	g.GET("/products/:productId", p.GetProductByID)
	g.PUT("/products/:productId", p.UpdateProductByID)
	g.DELETE("/products/:productId", p.DeleteProductByID)

	c := handler.NewCategoryHandler(server)
	g.GET("/category", c.GetAllCategory)
	g.POST("/category", c.SaveCategory)
	g.GET("/category/:categoryId", c.GetCategoryByID)
	g.PUT("/category/:categoryId", c.UpdateCategoryByID)
	g.DELETE("/category/:categoryId", c.DeleteCategoryByID)

	o := handler.NewOrderItemsHandler(server)
	g.POST("/orderitems", o.SaveOrderItems)
	g.GET("/orderitems/:userId/items", o.GetOrderItemsByUserID)
}
