package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create sales godoc
// @ID create_sales
// @Router /sales [POST]
// @Summary Create Sales
// @Description Create Sales
// @Tags Sales
// @Accept json
// @Procedure json
// @Param Sales body models.CreateSales true "CreateSalesRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateSales(c *gin.Context) {

	var createSales models.CreateSales
	err := c.ShouldBindJSON(&createSales)
	if err != nil {
		h.handlerResponse(c, "error Sales should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Sales().Create(c.Request.Context(), &createSales)
	if err != nil {
		h.handlerResponse(c, "storage.Sales.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Sales().GetByID(c.Request.Context(), &models.SalesPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sales.getById", http.StatusInternalServerError, err.Error())
		return
	}

	body := map[string]interface{}{}
	if resp.PaymentType == "Cash" {
		body = map[string]interface{}{
			"sales_id":    id,
			"type":        "TopUp",
			"source_type": "Sales",
			"text":        "Sales finished successfully",
			"staff_id":    resp.CashierID,
			"amount":      resp.Price,
		}
		_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "POST", body)

		if err != nil {
			h.handlerResponse(c, "storage.do request.cash cashier", http.StatusInternalServerError, err.Error())
			return
		}
		if len(resp.ShopAssistantID) > 0 {
			body = map[string]interface{}{
				"sales_id":    resp.Id,
				"type":        "TopUp",
				"source_type": "Sales",
				"text":        "Sales finished successfully",
				"staff_id":    resp.ShopAssistantID,
				"amount":      resp.Price,
			}
			_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "POST", body)
			if err != nil {
				h.handlerResponse(c, "storage.do request.cash assistent", http.StatusInternalServerError, err.Error())
				return
			}
		}
	} else if resp.PaymentType == "Card" {
		body = map[string]interface{}{
			"sales_id":    id,
			"type":        "TopUp",
			"source_type": "Sales",
			"text":        "Sales finished successfully",
			"staff_id":    resp.CashierID,
			"amount":      resp.Price,
		}
		_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "POST", body)
		if err != nil {
			h.handlerResponse(c, "storage.do request.cash cahsier", http.StatusInternalServerError, err.Error())
			return
		}
		if len(resp.ShopAssistantID) > 0 {
			body = map[string]interface{}{
				"sales_id":    resp.Id,
				"type":        "TopUp",
				"source_type": "Sales",
				"text":        "Sales finished successfully",
				"staff_id":    resp.ShopAssistantID,
				"amount":      resp.Price,
			}
			_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "POST", body)
			if err != nil {
				h.handlerResponse(c, "storage.do request.cash assistent", http.StatusInternalServerError, err.Error())
				return
			}
		}
	}

	h.handlerResponse(c, "create Sales resposne", http.StatusCreated, resp)
}

// GetByID sales godoc
// @ID get_by_id_sales
// @Router /sales/{id} [GET]
// @Summary Get By ID Sales
// @Description Get By ID Sales
// @Tags Sales
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdSales(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Sales().GetByID(c.Request.Context(), &models.SalesPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sales.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id Sales resposne", http.StatusOK, resp)
}

// GetList sales godoc
// @ID get_list_sales
// @Router /sales [GET]
// @Summary Get List Sales
// @Description Get List Sales
// @Tags Sales
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Param search query string false "search"
// @Param from query string false "from"
// @Param to query string false "to"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListSales(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list Sales offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list Sales limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Sales().GetList(c.Request.Context(), &models.SalesGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
		From:   c.Query("from"),
		To:     c.Query("to"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.Sales.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list Sales resposne", http.StatusOK, resp)
}

// Update sales godoc
// @ID update_sales
// @Router /sales/{id} [PUT]
// @Summary Update Sales
// @Description Update Sales
// @Tags Sales
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param Sales body models.UpdateSales true "UpdateSalesRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateSales(c *gin.Context) {

	var (
		id          string = c.Param("id")
		updateSales models.UpdateSales
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateSales)
	if err != nil {
		h.handlerResponse(c, "error Sales should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateSales.Id = id
	rowsAffected, err := h.strg.Sales().Update(c.Request.Context(), &updateSales)
	if err != nil {
		h.handlerResponse(c, "storage.Sales.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.Sales.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Sales().GetByID(c.Request.Context(), &models.SalesPrimaryKey{Id: updateSales.Id})
	if err != nil {
		h.handlerResponse(c, "storage.Sales.getById", http.StatusInternalServerError, err.Error())
		return
	}
	cashier, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: resp.CashierID})
	if err != nil {
		h.handlerResponse(c, "storage.sales->staff.getById", http.StatusInternalServerError, err.Error())
		return
	}
	tarif, err := h.strg.StaffTarif().GetByID(c.Request.Context(), &models.StaffTarifPrimaryKey{Id: cashier.TarifID})
	if err != nil {
		h.handlerResponse(c, "storage.sales->tarif.getById", http.StatusInternalServerError, err.Error())
		return
	}

	if resp.Status == "cancel" {
		if tarif.Type == "fixed" {
			if resp.PaymentType == "Cash" {
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom tarif.getById", http.StatusInternalServerError, err.Error())
					return
				}
				cashier.Balance -= tarif.AmountForCash
				updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       cashier.Id,
					BranchID: cashier.BranchID,
					TarifID:  cashier.TarifID,
					Type:     cashier.Type,
					Name:     cashier.Name,
					Balance:  cashier.Balance,
				})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom staff.getById", http.StatusInternalServerError, err.Error())
					return
				}
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				if len(resp.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: resp.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance -= tarif.AmountForCash
					updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{

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
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				}

			} else if resp.PaymentType == "Card" {
				cashier.Balance -= tarif.AmountForCard
				updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       cashier.Id,
					BranchID: cashier.BranchID,
					TarifID:  cashier.TarifID,
					Type:     cashier.Type,
					Name:     cashier.Name,
					Balance:  cashier.Balance,
				})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->staff.update", http.StatusInternalServerError, err.Error())
					return
				}
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				if len(resp.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: resp.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff assistent.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance -= tarif.AmountForCard
					updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
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
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				}

			}
		} else if tarif.Type == "percent" {
			if resp.PaymentType == "Cash" {
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->tarif.getById", http.StatusInternalServerError, err.Error())
					return
				}
				cashier.Balance -= (resp.Price * tarif.AmountForCash) / 100
				updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
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
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				if len(resp.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: resp.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance -= (resp.Price * tarif.AmountForCash) / 100
					updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{

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
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				}

			} else if resp.PaymentType == "Card" {
				cashier.Balance -= (resp.Price * tarif.AmountForCard) / 100
				updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
					Id:       cashier.Id,
					BranchID: cashier.BranchID,
					TarifID:  cashier.TarifID,
					Type:     cashier.Type,
					Name:     cashier.Name,
					Balance:  cashier.Balance,
				})
				if err != nil {
					h.handlerResponse(c, "storage.staffTransactiom->staff.update", http.StatusInternalServerError, err.Error())
					return
				}
				h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				if len(resp.ShopAssistantID) > 0 {
					shop_assistent, err := h.strg.Staff().GetByID(c.Request.Context(), &models.StaffPrimaryKey{Id: resp.ShopAssistantID})
					if err != nil {
						h.handlerResponse(c, "storage.staffTransactiom->staff assistent.getById", http.StatusInternalServerError, err.Error())
						return
					}
					shop_assistent.Balance -= (resp.Price * tarif.AmountForCard) / 100
					updateStaff, err := h.strg.Staff().Update(c.Request.Context(), &models.UpdateStaff{
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
					h.handlerResponse(c, "update staff  resposne", http.StatusOK, updateStaff)
				}

			}
		}
	}
	// if resp.Status == "cancel" {
	// 	body := map[string]interface{}{}
	// 	if resp.PaymentType == "Cash" {
	// 		body = map[string]interface{}{
	// 			"sales_id":    id,
	// 			"type":        "Withdraw",
	// 			"source_type": "Sales",
	// 			"text":        "Sales finished successfully",
	// 			"staff_id":    resp.CashierID,
	// 			"amount":      resp.Price,
	// 		}
	// 		_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "PUT", body)

	// 		if len(resp.ShopAssistantID) > 0 {
	// 			body = map[string]interface{}{
	// 				"sales_id":    resp.Id,
	// 				"type":        "Withdraw",
	// 				"source_type": "Sales",
	// 				"text":        "Sales finished successfully",
	// 				"staff_id":    resp.ShopAssistantID,
	// 				"amount":      resp.Price,
	// 			}

	// 		}
	// 		_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "PUT", body)

	// 	} else if resp.PaymentType == "Card" {
	// 		body = map[string]interface{}{
	// 			"sales_id":    id,
	// 			"type":        "Withdraw",
	// 			"source_type": "Sales",
	// 			"text":        "Sales finished successfully",
	// 			"staff_id":    resp.CashierID,
	// 			"amount":      resp.Price,
	// 		}
	// 		_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "PUT", body)

	// 		if len(resp.ShopAssistantID) > 0 {
	// 			body = map[string]interface{}{
	// 				"sales_id":    resp.Id,
	// 				"type":        "Withdraw",
	// 				"source_type": "Sales",
	// 				"text":        "Sales finished successfully",
	// 				"staff_id":    resp.ShopAssistantID,
	// 				"amount":      resp.Price,
	// 			}
	// 			_, err = helper.DoRequest("http://localhost:8080/staffTransaction", "PUT", body)

	// 		}
	// 	}
	// }

	h.handlerResponse(c, "create Sales resposne", http.StatusAccepted, resp)
}

// Delete sales godoc
// @ID delete_sales
// @Router /sales/{id} [DELETE]
// @Summary Delete Sales
// @Description Delete Sales
// @Tags Sales
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteSales(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Sales().Delete(c.Request.Context(), &models.SalesPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.Sales.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create Sales resposne", http.StatusNoContent, nil)
}
