package distributions

import (
	"net/http"
	"vkspam/database"
	"vkspam/repositories"
	"vkspam/services"
)

type DistributionHandler struct {
	service services.DistributionService
}

func NewDistributionHandler() *DistributionHandler {
	db, _ := database.GetDBInstance()

	distributionRepository := repositories.NewDistributionRepository(db.Db)
	distributionService := services.NewDistributionService(distributionRepository)

	return &DistributionHandler{
		service: distributionService,
	}
}

func (*DistributionHandler) Distribution(w http.ResponseWriter, r *http.Request) {

}
