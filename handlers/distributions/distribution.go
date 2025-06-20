package distributions

import (
	"fmt"
	"net/http"
	"strconv"
	"vkspam/database"
	"vkspam/handlers"
	"vkspam/models"
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

func (h *DistributionHandler) Distribution(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		id := r.FormValue("id")

		if len(id) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				fmt.Sprintf("Missing required parameter 'id'"),
			)

			return
		}

		recordId, err := strconv.Atoi(id)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
			)

			return
		}

		model, err := h.service.Get(recordId)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				err.Error(),
			)

			return
		}

		handlers.ReturnAppBaseResponse(
			w,
			http.StatusOK,
			true,
			model,
		)
		return
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		groupId := r.FormValue("group_id")
		if len(groupId) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'group_id'",
			)
			return
		}

		name := r.FormValue("name")
		if len(name) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'name'",
			)
			return
		}

		groupIdInt, err := strconv.Atoi(groupId)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf("Error conversation: %s", err),
			)
			return
		}

		url := r.FormValue("url")
		if len(url) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'url'",
			)
			return
		}

		typeF := r.FormValue("type")
		if len(typeF) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'type'",
			)
			return
		}

		if typeF != "Any public" {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Not supported variable option 'type'. Only 'Any public'",
			)
			return
		}

		distribution := models.Distribution{
			GroupId: groupIdInt,
			Type:    typeF,
			Url:     url,
			Name:    name,
		}

		if r.Method == http.MethodPut {
			id := r.FormValue("id")

			if len(id) < 1 {
				handlers.ReturnAppBaseResponse(
					w,
					http.StatusBadRequest,
					false,
					"Missing required parameter 'id'",
				)
				return
			}

			idInt, err := strconv.Atoi(id)
			if err != nil {
				handlers.ReturnAppBaseResponse(
					w,
					http.StatusInternalServerError,
					false,
					fmt.Sprintf("Error conversation: %s", err),
				)
				return
			}

			distribution.Id = idInt
		}

		id, err := h.service.Save(distribution)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf("Save error: %s", err),
			)
			return
		}

		if id == 0 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf("Save error parse id: %s", err),
			)
			return
		}
		distribution.Id = id
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusCreated,
			true,
			distribution,
		)
		return
	}

	if r.Method == http.MethodDelete {
		id := r.FormValue("id")

		if len(id) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				fmt.Sprintf("Missing required parameter 'id'"),
			)

			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
			)

			return
		}

		err = h.service.DeleteById(idInt)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				err.Error(),
			)

			return
		}
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusOK,
			false,
			nil,
		)

		return
	}
}
