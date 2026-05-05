package handlers

import (
	"pos-backend/internal/models"
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo     *repository.MenuRepository
	categoryRepo *repository.CategoryRepository
}

func NewMenuHandler(menuRepo *repository.MenuRepository, categoryRepo *repository.CategoryRepository) *MenuHandler {
	return &MenuHandler{
		menuRepo:     menuRepo,
		categoryRepo: categoryRepo,
	}
}

func (h *MenuHandler) GetMenu(c *gin.Context) {
	items, err := h.menuRepo.FindAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	// Enrich with category names
	categories, _ := h.categoryRepo.FindAll()
	categoryMap := make(map[string]string)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat.Name
	}

	for i := range items {
		if name, ok := categoryMap[items[i].Category]; ok {
			items[i].CategoryName = name
		}
	}

	utils.SuccessResponse(c, items)
}

func (h *MenuHandler) CreateMenuItem(c *gin.Context) {
	var item models.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.menuRepo.Create(&item); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, item)
}

func (h *MenuHandler) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.menuRepo.Update(id, updates); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Menu item updated successfully"})
}

func (h *MenuHandler) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")
	if err := h.menuRepo.Delete(id); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Menu item deleted successfully"})
}
