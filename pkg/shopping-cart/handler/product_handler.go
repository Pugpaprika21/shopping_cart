package handler

import (
	"fmt"
	"net/http"

	"github.com/Pugpaprika21/internal/utils"
	"github.com/Pugpaprika21/pkg/shopping-cart/dto"
	"github.com/Pugpaprika21/pkg/shopping-cart/models"
	"github.com/Pugpaprika21/pkg/shopping-cart/server"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type productHandler struct {
	server *server.EchoServerEnvironment
	query  *gorm.DB
}

func NewProductHandler(server *server.EchoServerEnvironment) *productHandler {
	return &productHandler{
		server: server,
		query:  server.Connect.ORM,
	}
}

func (p *productHandler) GetAllProduct(c echo.Context) error {
	var products []dto.ProductFetchRow
	p.query.Table("products as p").
		Select("p.id, p.name AS product_name, p.price AS product_price, p.qty AS product_quantity, p.detail AS product_detail, p.category_id, c.name AS category_name, CONCAT(pa.filename, pa.extension) AS product_attachment").
		Joins("INNER JOIN product_attachments AS pa ON p.id = pa.product_id").
		Joins("INNER JOIN categories AS c ON p.category_id = c.id").
		Order("pa.id DESC").
		Scan(&products)

	var productsResp []dto.ProductResp
	for _, product := range products {
		productsResp = append(productsResp, dto.ProductResp{
			ID:                product.ID,
			ProductName:       product.ProductName,
			ProductPrice:      product.ProductPrice,
			ProductQuantity:   product.ProductQuantity,
			ProductDetail:     product.ProductDetail,
			CategoryName:      product.CategoryName,
			ProductAttachment: fmt.Sprintf("%sproducts/attachments/%s", p.server.App.BaseURL, product.ProductAttachment),
		})
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
		Data:       productsResp,
	})
}

func (p *productHandler) GetProductAttachment(c echo.Context) error {
	attachment := c.Param("attachment")
	filePath := fmt.Sprintf("pkg/shopping-cart/uploads/%s", attachment)
	return c.File(filePath)
}

func (p *productHandler) SaveProduct(c echo.Context) error {
	var body dto.ProductReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	prodcut := models.Product{
		Name:       body.Name,
		Price:      body.Price,
		Qty:        body.Qty,
		Detail:     body.Detail,
		CategoryID: body.CategoryID,
	}

	err := p.query.Create(&prodcut).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	form, _ := c.MultipartForm()
	filehandler, _ := utils.NewUploadFile().SetPath("pkg/shopping-cart/uploads").GetMultiparts(form)
	files := filehandler.GetResultUploader()

	productAtt := models.ProductAttachment{
		Filename:  files[0].Filename,
		Path:      files[0].Path,
		Size:      files[0].Size,
		MimeType:  files[0].MimeType,
		Extension: files[0].Extension,
		ProductID: prodcut.ID,
	}

	err = p.query.Create(&productAtt).Error
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

func (p *productHandler) GetProductByID(c echo.Context) error {
	productID := utils.UintFromString(c.Param("productId"))
	var product dto.ProductFetchRow
	p.query.Table("products as p").
		Select("p.id, p.name AS product_name, p.price AS product_price, p.qty AS product_quantity, p.detail AS product_detail, p.category_id, c.name AS category_name, CONCAT(pa.filename, pa.extension) AS product_attachment").
		Joins("INNER JOIN product_attachments AS pa ON p.id = pa.product_id").
		Where("p.id = ?", productID).
		Order("pa.id DESC").
		Scan(&product)

	if product.ID == 0 {
		return c.JSON(http.StatusNotFound, dto.Respone{
			Statusbool: false,
			Data:       nil,
		})
	}

	productsResp := dto.ProductResp{
		ID:                product.ID,
		ProductName:       product.ProductName,
		ProductPrice:      product.ProductPrice,
		ProductQuantity:   product.ProductQuantity,
		ProductDetail:     product.ProductDetail,
		CategoryName:      product.CategoryName,
		ProductAttachment: fmt.Sprintf("%sproducts/attachments/%s", p.server.App.BaseURL, product.ProductAttachment),
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
		Data:       productsResp,
	})
}

func (p *productHandler) UpdateProductByID(c echo.Context) error {
	productID := utils.UintFromString(c.Param("productId"))
	var body dto.ProductReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	prodcut := models.Product{
		Name:       body.Name,
		Price:      body.Price,
		Qty:        body.Qty,
		Detail:     body.Detail,
		CategoryID: body.CategoryID,
	}

	err := p.query.Where("id = ?", productID).Updates(&prodcut).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	form, _ := c.MultipartForm()
	filehandler, _ := utils.NewUploadFile().SetPath("pkg/shopping-cart/uploads").GetMultiparts(form)
	files := filehandler.GetResultUploader()

	productAtt := models.ProductAttachment{
		Filename:  files[0].Filename,
		Path:      files[0].Path,
		Size:      files[0].Size,
		MimeType:  files[0].MimeType,
		Extension: files[0].Extension,
		ProductID: prodcut.ID,
	}

	err = p.query.Where("product_id = ?", productID).Updates(&productAtt).Error
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

func (p *productHandler) DeleteProductByID(c echo.Context) error {
	productID := utils.UintFromString(c.Param("productId"))

	var product models.Product
	p.query.Unscoped().Where("id = ?", productID).Delete(&product)

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}
