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
