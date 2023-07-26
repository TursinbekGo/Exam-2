package api

import (
	_ "app/api/docs"
	"app/api/handler"
	"app/config"
	"app/pkg/logger"
	"app/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewApi(r *gin.Engine, cfg *config.Config, storage storage.StorageI, logger logger.LoggerI) {

	handler := handler.NewHandler(cfg, storage, logger)

	r.POST("/branch", handler.CreateBranch)
	r.GET("/branch/:id", handler.GetByIdBranch)
	r.GET("/branch", handler.GetListBranch)
	r.PUT("/branch/:id", handler.UpdateBranch)
	r.DELETE("/branch/:id", handler.DeleteBranch)

	r.POST("/staffTarif", handler.CreateStaffTarif)
	r.GET("/staffTarif/:id", handler.GetByIdStaffTarif)
	r.GET("/staffTarif", handler.GetListStaffTarif)
	r.PUT("/staffTarif/:id", handler.UpdateStaffTarif)
	r.DELETE("/staffTarif/:id", handler.DeleteStaffTarif)

	r.POST("/staff", handler.CreateStaff)
	r.GET("/staff/:id", handler.GetByIdStaff)
	r.GET("/staff", handler.GetListStaff)
	r.PUT("/staff/:id", handler.UpdateStaff)
	r.DELETE("/staff/:id", handler.DeleteStaff)

	r.POST("/sales", handler.CreateSales)
	r.GET("/sales/:id", handler.GetByIdSales)
	r.GET("/sales", handler.GetListSales)
	r.PUT("/sales/:id", handler.UpdateSales)
	r.DELETE("/sales/:id", handler.DeleteSales)

	r.POST("/staffTransaction", handler.CreateStaffTransaction)
	r.GET("/staffTransaction/:id", handler.GetByIdStaffTransaction)
	r.GET("/staffTransaction", handler.GetListStaffTransaction)
	r.PUT("/staffTransaction/:id", handler.UpdateStaffTransaction)
	r.DELETE("/staffTransaction/:id", handler.DeleteStaffTransaction)

	url := ginSwagger.URL("swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}
