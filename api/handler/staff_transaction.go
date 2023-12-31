package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create staffTransaction godoc
// @ID create_staffTransaction
// @Router /staffTransaction [POST]
// @Summary Create StaffTransaction
// @Description Create StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param StaffTransaction body models.CreateStaffTransaction true "CreateStaffTransactionRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateStaffTransaction(c *gin.Context) {

	var createStaffTransaction models.CreateStaffTransaction
	err := c.ShouldBindJSON(&createStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "error StaffTransaction should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.StaffTransaction().Create(c.Request.Context(), &createStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	saless, err := h.strg.Sales().GetByID(c.Request.Context(), &models.SalesPrimaryKey{Id: resp.SalesID})
	if err != nil {
		h.handlerResponse(c, "storage.saless.getById", http.StatusInternalServerError, err.Error())
		return
	}

	// stafff, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: resp.StaffID})
	// if err != nil {
	// 	h.handlerResponse(c, "storage.Staff.getById", http.StatusInternalServerError, err.Error())
	// 	return
	// }
	cashier, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: saless.CashierID})
	if err != nil {
		h.handlerResponse(c, "storage.chashier.getById", http.StatusInternalServerError, err.Error())
		return
	}
	tarifff, err := h.strg.StaffTarif().GetByID(c.Request.Context(), &models.StaffTarifPrimaryKey{Id: cashier.TarifID})
	if err != nil {
		h.handlerResponse(c, "storage.tarifff.getById", http.StatusInternalServerError, err.Error())
		return
	}

	if resp.Type == "success" {
		if saless.PaymentType == "Cash" {
			if tarifff.Type == "fixed" {
				cashier.Balance += tarifff.AmountForCash
				update, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       cashier.Id,
					BranchID: cashier.BranchID,
					TarifID:  cashier.TarifID,
					Type:     cashier.Type,
					Name:     cashier.Name,
					Balance:  cashier.Balance,
				})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
					return
				}
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, update)
				if len(saless.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: saless.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance += tarifff.AmountForCash
					update, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{

						Id:       shop_assistent.Id,
						BranchID: shop_assistent.BranchID,
						TarifID:  shop_assistent.TarifID,
						Type:     shop_assistent.Type,
						Name:     shop_assistent.Name,
						Balance:  shop_assistent.Balance,
					})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->tarif.getById", http.StatusInternalServerError, err.Error())
						return
					}
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, update)
				}
			} else if tarifff.Type == "percent" {
				cashier.Balance = cashier.Balance + ((saless.Price * tarifff.AmountForCash) / 100)
				h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       cashier.Id,
					BranchID: cashier.BranchID,
					TarifID:  cashier.TarifID,
					Type:     cashier.Type,
					Name:     cashier.Name,
					Balance:  cashier.Balance,
				})
				if len(saless.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: saless.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance = shop_assistent.Balance + ((saless.Price * tarifff.AmountForCash) / 100)
					update, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
						Id:       shop_assistent.Id,
						BranchID: shop_assistent.BranchID,
						TarifID:  shop_assistent.TarifID,
						Type:     shop_assistent.Type,
						Name:     shop_assistent.Name,
						Balance:  shop_assistent.Balance,
					})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->tarif.getById", http.StatusInternalServerError, err.Error())
						return
					}
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, update)
				}
			}
		} else if saless.PaymentType == "Card" {
			if tarifff.Type == "fixed" {
				cashier.Balance += tarifff.AmountForCard
				update, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       cashier.Id,
					BranchID: cashier.BranchID,
					TarifID:  cashier.TarifID,
					Type:     cashier.Type,
					Name:     cashier.Name,
					Balance:  cashier.Balance,
				})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
					return
				}
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, update)
				if len(saless.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: saless.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance += tarifff.AmountForCard
					update, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{

						Id:       shop_assistent.Id,
						BranchID: shop_assistent.BranchID,
						TarifID:  shop_assistent.TarifID,
						Type:     shop_assistent.Type,
						Name:     shop_assistent.Name,
						Balance:  shop_assistent.Balance,
					})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->tarif.getById", http.StatusInternalServerError, err.Error())
						return
					}
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, update)
				}
			}
		} else if tarifff.Type == "percent" {
			cashier.Balance = cashier.Balance + ((saless.Price * tarifff.AmountForCard) / 100)
			h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
				Id:       cashier.Id,
				BranchID: cashier.BranchID,
				TarifID:  cashier.TarifID,
				Type:     cashier.Type,
				Name:     cashier.Name,
				Balance:  cashier.Balance,
			})
			if len(saless.ShopAssistantID) > 0 {
				shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: saless.ShopAssistantID})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
					return
				}
				shop_assistent.Balance = shop_assistent.Balance + ((saless.Price * tarifff.AmountForCash) / 100)
				update, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       shop_assistent.Id,
					BranchID: shop_assistent.BranchID,
					TarifID:  shop_assistent.TarifID,
					Type:     shop_assistent.Type,
					Name:     shop_assistent.Name,
					Balance:  shop_assistent.Balance,
				})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->tarif.getById", http.StatusInternalServerError, err.Error())
					return
				}
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, update)
			}
		}

	}

	h.handlerResponse(c, "create StaffTransaction resposne", http.StatusCreated, resp)

}

// GetByID staffTransaction godoc
// @ID get_by_id_staffTransaction
// @Router /staffTransaction/{id} [GET]
// @Summary Get By ID StaffTransaction
// @Description Get By ID StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdStaffTransaction(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id StaffTransaction resposne", http.StatusOK, resp)
}

// GetList staffTransaction godoc
// @ID get_list_staffTransaction
// @Router /staffTransaction [GET]
// @Summary Get List StaffTransaction
// @Description Get List StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Param salesId query string false "sales_id"
// @Param staffId query string false "staff_id"
// @Param type query string false "type"
// @Param amount query string false "amount"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListStaffTransaction(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list StaffTransaction offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list StaffTransaction limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.StaffTransaction().GetList(c.Request.Context(), &models.StaffTransactionGetListRequest{
		Offset:  offset,
		Limit:   limit,
		Search:  c.Query("search"),
		SalesId: c.Query("sales_id"),
		StaffId: c.Query("staff_id"),
		Type:    c.Query("type"),
		Amount:  c.Query("amount"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list StaffTransaction resposne", http.StatusOK, resp)
}

// Update staffTransaction godoc
// @ID update_staffTransaction
// @Router /staffTransaction/{id} [PUT]
// @Summary Update StaffTransaction
// @Description Update StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param StaffTransaction body models.UpdateStaffTransaction true "UpdateStaffTransactionRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateStaffTransaction(c *gin.Context) {

	var (
		id                     string = c.Param("id")
		updateStaffTransaction models.UpdateStaffTransaction
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "error StaffTransaction should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateStaffTransaction.Id = id
	rowsAffected, err := h.strg.StaffTransaction().Update(c.Request.Context(), &updateStaffTransaction)
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.StaffTransaction.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.StaffTransaction().GetByID(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: updateStaffTransaction.Id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create StaffTransaction resposne", http.StatusAccepted, resp)
}

// Delete staffTransaction godoc
// @ID delete_staffTransaction
// @Router /staffTransaction/{id} [DELETE]
// @Summary Delete StaffTransaction
// @Description Delete StaffTransaction
// @Tags StaffTransaction
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteStaffTransaction(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.StaffTransaction().Delete(c.Request.Context(), &models.StaffTransactionPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.StaffTransaction.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create StaffTransaction resposne", http.StatusNoContent, nil)
}
