package handlers

import (
	"pos-backend/internal/models"
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
}

func NewCategoryHandler(repo *repository.CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	categories, err := h.repo.FindAll(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, categories)
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.Create(c.Request.Context(), &category); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, category)
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, updates); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Category updated successfully"})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Category deleted successfully"})
}

func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	var req struct {
		Categories []struct {
			ID        string `json:"id"`
			SortOrder int    `json:"sortOrder"`
		} `json:"categories"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	for _, cat := range req.Categories {
		if err := h.repo.Update(c.Request.Context(), cat.ID, bson.M{"sortOrder": cat.SortOrder}); err != nil {
			utils.InternalServerErrorResponse(c, err.Error())
			return
		}
	}

	utils.SuccessResponse(c, gin.H{"message": "Categories reordered successfully"})
}