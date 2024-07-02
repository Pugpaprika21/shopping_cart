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

type categoryHandler struct {
	server *server.EchoServerEnvironment
	query  *gorm.DB
}

func NewCategoryHandler(server *server.EchoServerEnvironment) *categoryHandler {
	return &categoryHandler{
		server: server,
		query:  server.Connect.ORM,
	}
}

func (ch *categoryHandler) GetAllCategory(c echo.Context) error {
	var categories []dto.CategoryFetchRow
	ch.query.Table("categories").
		Select("id, name AS category_name, description AS category_description").
		Order("id DESC").
		Scan(&categories)

	var categoryResp []dto.CategoryResp
	for _, category := range categories {
		categoryResp = append(categoryResp, dto.CategoryResp{
			ID:                  category.ID,
			CategoryName:        category.CategoryName,
			CategoryDescription: category.CategoryDescription,
		})
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
		Data:       categoryResp,
	})
}

func (ch *categoryHandler) SaveCategory(c echo.Context) error {
	var body dto.CategoryReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	category := models.Category{
		Name:        body.Name,
		Description: body.Description,
	}

	err := ch.query.Create(&category).Error
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})

	// 444444
}

func (ch *categoryHandler) GetCategoryByID(c echo.Context) error {
	categoryID := utils.UintFromString(c.Param("categoryId"))
	var categories dto.CategoryFetchRow
	ch.query.Table("categories").
		Select("id, name AS category_name, description AS category_description").
		Where("id = ?", categoryID).
		Order("id DESC").
		Scan(&categories)

	if categories.ID == 0 {
		return c.JSON(http.StatusNotFound, dto.Respone{
			Statusbool: false,
			Data:       nil,
		})
	}

	categoryResp := dto.CategoryResp{
		ID:                  categories.ID,
		CategoryName:        categories.CategoryName,
		CategoryDescription: categories.CategoryDescription,
	}

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
		Data:       categoryResp,
	})
}

func (ch *categoryHandler) UpdateCategoryByID(c echo.Context) error {
	categoryID := utils.UintFromString(c.Param("categoryId"))
	var body dto.CategoryReqBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Respone{
			Message:    err.Error(),
			Statusbool: false,
		})
	}

	// update Category

	category := models.Category{
		Name:        body.Name,
		Description: body.Description,
	}

	err := ch.query.Where("id = ?", categoryID).Updates(&category).Error
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

func (ch *categoryHandler) DeleteCategoryByID(c echo.Context) error {
	categoryID := utils.UintFromString(c.Param("categoryId"))

	var product models.Category
	ch.query.Unscoped().Where("id = ?", categoryID).Delete(&product)

	return c.JSON(http.StatusOK, dto.Respone{
		Statusbool: true,
	})
}
