package handler

import (
	"net/http"
	"time"

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

	var prodcuts models.Product
	prodcutLength := len(body.Products)
	if prodcutLength > 0 {
		for _, product := range body.Products {
			var currentProduct struct {
				ID              uint
				ProductQuantity uint
			}
			// ตรวจสอบจำนวนสินค้าที่มีอยู่
			err := o.query.Model(&prodcuts).Select("id, qty AS product_quantity").Where("id = ?", product.ProductID).Scan(&currentProduct).Error
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
			err = o.query.Model(&prodcuts).Where("id = ?", product.ProductID).Update("qty", newQuantity).Error
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
				UserID:          body.UserID,
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

func (o *orderItemsHandler) UpdateOrderItemsByUserID(c echo.Context) error {
	userID := utils.UintFromString(c.Param("userId"))

	var body dto.OrderItemsReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	var orderCount int64
	var orderItems models.OrderItems
	o.query.Model(&orderItems).Where("status = ? AND user_id = ?", models.Pending, userID).Count(&orderCount)
	if orderCount == 0 {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    "ไม่พบข้อมูลการสั่งซื้อของคุณ",
			Statusbool: false,
		})
	}

	var products models.Product
	productLength := len(body.Products)
	if productLength > 0 {
		for _, product := range body.Products {
			var currentProduct struct {
				ID              uint
				ProductQuantity uint
			}
			// ตรวจสอบจำนวนสินค้าที่มีอยู่
			err := o.query.Model(&products).Select("id, qty AS product_quantity").Where("id = ?", product.ProductID).Scan(&currentProduct).Error
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
			err = o.query.Model(&products).Where("id = ?", product.ProductID).Update("qty", newQuantity).Error
			if err != nil {
				return c.JSON(http.StatusInternalServerError, dto.Respone{
					Message:    err.Error(),
					Statusbool: false,
				})
			}
			// ตรวจสอบว่ามี orderItem อยู่แล้วหรือไม่
			var existingOrderItem models.OrderItems
			err = o.query.Model(&existingOrderItem).Where("user_id = ? AND product_id = ? AND status = ?", userID, product.ProductID, models.Pending).First(&existingOrderItem).Error
			if err != nil {
				if gorm.ErrRecordNotFound == err {
					// ถ้าไม่พบ orderItem ให้เพิ่มใหม่
					orderItem := models.OrderItems{
						ProductID:       product.ProductID,
						ProductQuantity: product.ProductQuantity,
						UserID:          userID,
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
				} else {
					return c.JSON(http.StatusInternalServerError, dto.Respone{
						Message:    err.Error(),
						Statusbool: false,
					})
				}
			} else {
				// ถ้าพบ orderItem ให้ทำการอัปเดต
				existingOrderItem.ProductQuantity += product.ProductQuantity
				if err := o.query.Save(&existingOrderItem).Error; err != nil {
					return c.JSON(http.StatusInternalServerError, dto.Respone{
						Message:    err.Error(),
						Statusbool: false,
					})
				}
			}
		}
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}

func (o *orderItemsHandler) UpdateOrderItemsByID(c echo.Context) error {
	userID := utils.UintFromString(c.Param("userId"))
	orderID := utils.UintFromString(c.Param("orderId"))

	var orderCount int64
	var orderItems []dto.OrderItemsFetchRow
	var orderItemsModel models.OrderItems
	o.query.Model(&orderItemsModel).Select("id").Where("id = ? AND status = ? AND user_id = ?", orderID, models.Pending, userID).Scan(&orderItems).Count(&orderCount)
	if orderCount == 0 {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    "ไม่พบข้อมูลการสั่งซื้อของคุณ",
			Statusbool: false,
		})
	}

	err := o.query.Model(&orderItemsModel).Where("id = ? AND status = ? AND user_id = ?", orderID, models.Pending, userID).Update("status", models.Completed).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}

func (o *orderItemsHandler) ConfirmBuyOrderItems(c echo.Context) error {
	var body dto.OrderItemsReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	var orderCount int64
	var orderItems []dto.OrderItemsFetchRow
	var orderItemsModel models.OrderItems
	o.query.Model(&orderItemsModel).Select("id").Where("status = ? AND user_id = ?", models.Pending, body.UserID).Scan(&orderItems).Count(&orderCount)
	if orderCount == 0 {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    "ไม่พบข้อมูลการสั่งซื้อของคุณ",
			Statusbool: false,
		})
	}

	for _, order := range orderItems {
		o.query.Model(&orderItemsModel).Where("status = ? AND id = ?", models.Pending, order.ID).Update("status", models.Completed)
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}

func (o *orderItemsHandler) DeleteOrderItemsByID(c echo.Context) error {
	userID := utils.UintFromString(c.Param("userId"))
	orderID := utils.UintFromString(c.Param("orderId"))

	var orderCount int64
	var orderItems models.OrderItems
	o.query.Model(&orderItems).Where("user_id = ? AND id = ?", userID, orderID).Count(&orderCount)
	if orderCount == 0 {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    "ไม่พบข้อมูลการสั่งซื้อของคุณ",
			Statusbool: false,
		})
	}

	fields := map[string]interface{}{
		"status":     models.Cancelled,
		"updated_at": time.Now(),
		"deleted_at": time.Now(),
	}

	err := o.query.Model(&orderItems).Where("user_id = ? AND id = ?", userID, orderID).Updates(fields).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}
