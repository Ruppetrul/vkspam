package distributions

import (
	"encoding/json"
	"log"
	"net/http"
	"vkspam/database"
	"vkspam/middleware"
	"vkspam/repositories"
	"vkspam/services"
)

type DistributionGroupHandler struct {
	service services.DistributionGroupService
}

func NewDistributionGroupHandler() *DistributionGroupHandler {
	db, _ := database.GetDBInstance()

	distributionRepository := repositories.NewDistributionRepository(db.Db)
	distributionService := services.NewDistributionGroupService(distributionRepository)

	return &DistributionGroupHandler{
		service: distributionService,
	}
}

func (h *DistributionGroupHandler) Group(w http.ResponseWriter, r *http.Request) {
	middleware.GetUserFromContext(r.Context())

	if r.Method == http.MethodGet {
		_, err := w.Write([]byte("There will be group list."))
		if err != nil {
			log.Println("Error index page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost {
		//TODO create
	}

	if r.Method == http.MethodPut {
		//TODO update
	}

	if r.Method == http.MethodDelete {
		//TODO delete
	}
}

func (h *DistributionGroupHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	distributionGroups, err := h.service.GetList(user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(distributionGroups)
	if err != nil {
		http.Error(w, "Error return list groups", http.StatusInternalServerError)
		log.Println("Error return list groups:", err)
		return
	}
}
