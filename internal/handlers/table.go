package handlers

import (
	"strconv"

	"pos-backend/internal/models"
	"pos-backend/internal/repository"
	"pos-backend/internal/utils"

	"github.com/gin-gonic/gin"
)

type TableHandler struct {
	repo *repository.TableRepository
}

func NewTableHandler(repo *repository.TableRepository) *TableHandler {
	return &TableHandler{repo: repo}
}

func (h *TableHandler) GetTables(c *gin.Context) {
	tables, err := h.repo.FindAll()
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	utils.SuccessResponse(c, tables)
}

func (h *TableHandler) GetTableByNumber(c *gin.Context) {
	tableNumber, err := strconv.Atoi(c.Param("tableNumber"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid table number")
		return
	}

	table, err := h.repo.FindByNumber(tableNumber)
	if err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}
	if table == nil {
		utils.NotFoundResponse(c, "Table not found")
		return
	}

	utils.SuccessResponse(c, table)
}

func (h *TableHandler) CreateTable(c *gin.Context) {
	var table models.Table
	if err := c.ShouldBindJSON(&table); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.Create(&table); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.CreatedResponse(c, table)
}

func (h *TableHandler) UpdateTable(c *gin.Context) {
	tableNumber, err := strconv.Atoi(c.Param("tableNumber"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid table number")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.BadRequestResponse(c, "Invalid request body")
		return
	}

	if err := h.repo.UpdateByNumber(tableNumber, updates); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Table updated successfully"})
}

func (h *TableHandler) DeleteTable(c *gin.Context) {
	tableNumber, err := strconv.Atoi(c.Param("tableNumber"))
	if err != nil {
		utils.BadRequestResponse(c, "Invalid table number")
		return
	}

	if err := h.repo.DeleteByNumber(tableNumber); err != nil {
		utils.InternalServerErrorResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Table deleted successfully"})
}
