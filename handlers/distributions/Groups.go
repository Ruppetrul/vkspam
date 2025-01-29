package distributions

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"vkspam/database"
	"vkspam/middleware"
	"vkspam/models"
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
	user := middleware.GetUserFromContext(r.Context())

	if r.Method == http.MethodGet {
		_, err := w.Write([]byte("There will be group list."))
		if err != nil {
			log.Println("Error index page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		name := r.FormValue("name")
		description := r.FormValue("description")

		if len(name) < 1 {
			http.Error(w, "Missing required parameter 'name'", http.StatusBadRequest)
			return
		}
		if len(description) < 1 {
			http.Error(w, "Missing required parameter 'description'", http.StatusBadRequest)
			return
		}

		newDistributionGroup := models.DistributionGroup{
			Name:        r.Form.Get("name"),
			Description: r.Form.Get("description"),
			UserId:      user.Id,
		}

		if r.Method == http.MethodPut {
			id := r.FormValue("id")
			if len(id) < 1 {
				http.Error(w, "Missing required parameter 'id'", http.StatusBadRequest)
				return
			}

			recordId, err := strconv.Atoi(id)
			if err != nil {
				fmt.Println("Ошибка преобразования:", err)
				return
			}

			newDistributionGroup.Id = recordId
		}

		form := h.service.Save(newDistributionGroup)

		if form != nil {
			http.Error(w, form.Error(), http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodDelete {
		id := r.FormValue("id")
		if len(id) < 1 {
			http.Error(w, "Missing required parameter 'id'", http.StatusBadRequest)
			return
		}

		recordId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid 'id' parameter", http.StatusBadRequest)
			return
		}

		err = h.service.Delete(recordId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Group deleted successfully"))
		if err != nil {
			log.Println("Error writing response:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
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
