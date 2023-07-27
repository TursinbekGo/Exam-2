package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create branch godoc
// @ID create_branch
// @Router /branch [POST]
// @Summary Create Branch
// @Description Create Branch
// @Tags Branch
// @Accept json
// @Procedure json
// @Param Branch body models.CreateBranch true "CreateBranchRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateBranch(c *gin.Context) {

	var createBranch models.CreateBranch
	err := c.ShouldBindJSON(&createBranch)
	if err != nil {
		h.handlerResponse(c, "error branch should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Branch().Create(c.Request.Context(), &createBranch)
	if err != nil {
		h.handlerResponse(c, "storage.branch.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Branch().GetByID(c.Request.Context(), &models.BranchPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.branch.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create branch resposne", http.StatusCreated, resp)
}

// GetByID branch godoc
// @ID get_by_id_branch
// @Router /branch/{id} [GET]
// @Summary Get By ID Branch
// @Description Get By ID Branch
// @Tags Branch
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdBranch(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Branch().GetByID(c.Request.Context(), &models.BranchPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.branch.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id branch resposne", http.StatusOK, resp)
}

// GetList branch godoc
// @ID get_list_branch
// @Router /branch [GET]
// @Summary Get List Branch
// @Description Get List Branch
// @Tags Branch
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListBranch(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list branch offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list branch limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Branch().GetList(c.Request.Context(), &models.BranchGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.branch.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list branch resposne", http.StatusOK, resp)
}

// Update branch godoc
// @ID update_branch
// @Router /branch/{id} [PUT]
// @Summary Update Branch
// @Description Update Branch
// @Tags Branch
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Branch body models.UpdateBranch true "UpdateBranchRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateBranch(c *gin.Context) {

	var (
		id           string = c.Param("id")
		updateBranch models.UpdateBranch
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateBranch)
	if err != nil {
		h.handlerResponse(c, "error branch should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateBranch.Id = id
	rowsAffected, err := h.strg.Branch().Update(c.Request.Context(), &updateBranch)
	if err != nil {
		h.handlerResponse(c, "storage.branch.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.branch.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Branch().GetByID(c.Request.Context(), &models.BranchPrimaryKey{Id: updateBranch.Id})
	if err != nil {
		h.handlerResponse(c, "storage.branch.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create branch resposne", http.StatusAccepted, resp)
}

// Delete branch godoc
// @ID delete_branch
// @Router /branch/{id} [DELETE]
// @Summary Delete Branch
// @Description Delete Branch
// @Tags Branch
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteBranch(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Branch().Delete(c.Request.Context(), &models.BranchPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.branch.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create branch resposne", http.StatusNoContent, nil)
}

// GetTOP branchs godoc
// @ID get_top_branchs
// @Router /get_top_branchs [GET]
// @Summary Get Top Branchs
// @Description Get Top Branchs
// @Tags Get Top Branchs
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetTopBranchs(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list branch offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list branch limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Branch().GetTopBranch(c.Request.Context(), &models.BranchGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.branch.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list branch resposne", http.StatusOK, resp)
}
