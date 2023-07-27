package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create staff godoc
// @ID create_staff
// @Router /staff [POST]
// @Summary Create Staff
// @Description Create Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param Staff body models.CreateStaff true "CreateStaffRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateStaff(c *gin.Context) {

	var createStaff models.CreateStaff
	err := c.ShouldBindJSON(&createStaff)
	if err != nil {
		h.handlerResponse(c, "error staff should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Staff().Create(c.Request.Context(), &createStaff)
	if err != nil {
		h.handlerResponse(c, "storage.staff.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create staff resposne", http.StatusCreated, resp)
}

// GetByID staff godoc
// @ID get_by_id_staff
// @Router /staff/{id} [GET]
// @Summary Get By ID Staff
// @Description Get By ID Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdStaff(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Staff resposne", http.StatusOK, resp)
}

// GetList staff godoc
// @ID get_list_staff
// @Router /staff [GET]
// @Summary Get List Staff
// @Description Get List Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Param from query string false "from"
// @Param to query string false "to"
// @Param branch_id query string false "branch_id"
// @Param tarif_id query string false "tarif_id"
// @Param type query string false "type"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListStaff(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Staff offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Staff limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Staff().GetList(c.Request.Context(), &models.StaffGetListRequest{
		Offset:   offset,
		Limit:    limit,
		Search:   c.Query("search"),
		From:     c.Query("from"),
		To:       c.Query("to"),
		BranchId: c.Query("branch_id"),
		TarifID:  c.Query("tarif_id"),
		Type:     c.Query("type"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Staff resposne", http.StatusOK, resp)
}

// Update staff godoc
// @ID update_staff
// @Router /staff/{id} [PUT]
// @Summary Update Staff
// @Description Update Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Staff body models.UpdateStaff true "UpdateStaffRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateStaff(c *gin.Context) {

	var (
		id          string = c.Param("id")
		updateStaff models.UpdateStaff
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateStaff)
	if err != nil {
		h.handlerResponse(c, "error Staff should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateStaff.Id = id
	rowsAffected, err := h.strg.Staff().Update(c.Request.Context(), &updateStaff)
	if err != nil {
		h.handlerResponse(c, "storage.Staff.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Staff.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: updateStaff.Id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff resposne", http.StatusAccepted, resp)
}

// Delete staff godoc
// @ID delete_staff
// @Router /staff/{id} [DELETE]
// @Summary Delete Staff
// @Description Delete Staff
// @Tags Staff
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteStaff(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Staff().Delete(c.Request.Context(), &models.StaffPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Staff resposne", http.StatusNoContent, nil)
}

// GetTOP staffs godoc
// @ID get_top_staffs
// @Router /get_top_staffs [GET]
// @Summary Get Top Staffs
// @Description Get Top Staffs
// @Tags Get Top Staffs
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetTOPStaffs(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Staff offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Staff limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Staff().GetTopStaff(c.Request.Context(), &models.StaffGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Staff.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Staff resposne", http.StatusOK, resp)
}
