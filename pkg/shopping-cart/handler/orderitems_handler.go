package handler

import (
	"net/http"

	"github.com/Pugpaprika21/internal/utils"
	"github.com/Pugpaprika21/pkg/shopping-cart/dto"
	"github.com/Pugpaprika21/pkg/shopping-cart/models"
	"github.com/Pugpaprika21/pkg/shopping-cart/server"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type orderItemsHandler struct {
	server *server.EchoServerEnvironment
	query  *gorm.DB
}

func NewOrderItemsHandler(server *server.EchoServerEnvironment) *orderItemsHandler {
	return &orderItemsHandler{
		server: server,
		query:  server.Connect.ORM,
	}
}

func (o *orderItemsHandler) SaveOrderItems(c echo.Context) error {
	var body dto.OrderItemsReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	prodcutLength := len(body.Products)
	if prodcutLength > 0 {
		for _, product := range body.Products {
			var currentProduct struct {
				ID              uint
				ProductQuantity uint
			}

			// ตรวจสอบจำนวนสินค้าที่มีอยู่
			err := o.query.Table("products").Select("id, qty AS product_quantity").Where("id = ?", product.ProductID).Scan(&currentProduct).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.Respone{
					Message:    err.Error(),
					Statusbool: false,
				})
			}

			// ตรวจสอบว่าสินค้ามีเพียงพอหรือไม่
			if currentProduct.ProductQuantity < product.ProductQuantity {
				return c.JSON(http.StatusBadRequest, dto.Respone{
					Message:    "สินค้าไม่เพียงพอ",
					Statusbool: false,
				})
			}

			// ลดจำนวนสินค้าในสต็อก
			newQuantity := currentProduct.ProductQuantity - product.ProductQuantity
			err = o.query.Table("products").Where("id = ?", product.ProductID).Update("qty", newQuantity).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.Respone{
					Message:    err.Error(),
					Statusbool: false,
				})
			}

			// สร้างรายการสั่งซื้อใหม่
			orderItem := models.OrderItems{
				ProductID:       product.ProductID,
				ProductQuantity: product.ProductQuantity,
				UserID:          1,
				UserIP:          c.RealIP(),
				Status:          models.Pending,
			}

			// บันทึกข้อมูลการสั่งซื้อ
			if err := o.query.Create(&orderItem).Error; err != nil {
				return c.JSON(http.StatusInternalServerError, dto.Respone{
					Message:    err.Error(),
					Statusbool: false,
				})
			}
		}
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}

func (o *orderItemsHandler) GetOrderItemsByUserID(c echo.Context) error {
	userID := utils.UintFromString(c.Param("userId"))
	orderStatus := c.QueryParam("order_status")

	var orderCount int64
	var orderItems []dto.OrderItemsFetchRow
	o.query.Table("order_items AS oi").
		Select("oi.id, oi.product_id, oi.product_quantity, oi.user_id, p.name AS product_name, p.price AS product_price, c.name AS category_name, oi.status AS order_status").
		Joins("INNER JOIN products AS p ON oi.product_id = p.id").
		Joins("INNER JOIN categories AS c ON c.id = p.category_id").
		Where("oi.status = ? AND oi.user_id = ?", orderStatus, userID).
		Order("oi.id DESC").
		Count(&orderCount).
		Scan(&orderItems)

	if orderCount == 0 {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    "ไม่พบข้อมูลการสั่งซื้อของคุณ",
			Statusbool: false,
		})
	}

	const taxRate = 0.07
	var orderTotal float32
	var orderTotalWithTax float32
	var orderItemResp []dto.OrderItemsResp
	for _, orderItem := range orderItems {
		subtotal := float32(orderItem.ProductQuantity) * orderItem.ProductPrice
		tax := subtotal * taxRate
		total := subtotal + tax
		orderItemResp = append(orderItemResp, dto.OrderItemsResp{
			ID:              orderItem.ID,
			ProductQuantity: orderItem.ProductQuantity,
			ProductName:     orderItem.ProductName,
			ProductPrice:    orderItem.ProductPrice,
			CategoryName:    orderItem.CategoryName,
		})
		orderTotalWithTax += total
		orderTotal += float32(orderItem.ProductQuantity) * orderItem.ProductPrice
	}

	orderItemsRespBody := dto.OrderItemsRespBody{
		OrderItems:  orderItemResp,
		TotalTax:    orderTotalWithTax,
		Total:       orderTotal,
		OrderStatus: orderStatus,
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
		Data:       orderItemsRespBody,
	})
}

func (o *orderItemsHandler) SaveOrderItemsByID(c echo.Context) error {
	return nil
}

func (o *orderItemsHandler) DeleteOrderItemsByID(c echo.Context) error {
	return nil
}

func (o *orderItemsHandler) UpdateOrderItems(c echo.Context) error {
	return nil
}

func (o *orderItemsHandler) ConfirmAsSaveOrderItems(c echo.Context) error {
	return nil
}
