package distributions

import (
	"context"
	"fmt"
	pb "github.com/Ruppetrul/vkspam_proto/gen/parser"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
	"vkspam/database"
	"vkspam/handlers"
	"vkspam/handlers/responses"
	"vkspam/middleware"
	"vkspam/models"
	"vkspam/repositories"
	"vkspam/services"
)

type DistributionGroupHandler struct {
	service             services.DistributionGroupService
	distributionService services.DistributionService
}

func NewDistributionGroupHandler() *DistributionGroupHandler {
	db, _ := database.GetDBInstance()

	distributionGroupRepository := repositories.NewDistributionGroupRepository(db.Db)
	distributionGroupService := services.NewDistributionGroupService(distributionGroupRepository)

	distributionRepository := repositories.NewDistributionRepository(db.Db)
	distributionService := services.NewDistributionService(distributionRepository)

	return &DistributionGroupHandler{
		service:             distributionGroupService,
		distributionService: distributionService,
	}
}

func (h *DistributionGroupHandler) Group(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

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
		distributions, err := h.distributionService.GetListByGroup(recordId)

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
			responses.DistributionGroupResponse{
				Id:            model.Id,
				Name:          model.Name,
				Description:   model.Description,
				UserId:        model.UserId,
				Distributions: *distributions,
			},
		)
	}

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		name := r.FormValue("name")
		description := r.FormValue("description")

		if len(name) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'name'",
			)

			return
		}
		if len(description) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'description'",
			)

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
				handlers.ReturnAppBaseResponse(
					w,
					http.StatusBadRequest,
					false,
					"Missing required parameter 'id'",
				)

				return
			}

			recordId, err := strconv.Atoi(id)
			if err != nil {
				handlers.ReturnAppBaseResponse(
					w,
					http.StatusInternalServerError,
					false,
					fmt.Sprintf("Ошибка преобразования: %s", err),
				)

				return
			}

			newDistributionGroup.Id = recordId
		}

		err := h.service.Save(newDistributionGroup)

		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				err.Error(),
			)
		}

		handlers.ReturnAppBaseResponse(
			w,
			http.StatusOK,
			true,
			newDistributionGroup,
		)
	}

	if r.Method == http.MethodDelete {
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

		recordId, err := strconv.Atoi(id)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Invalid 'id' parameter",
			)

			return
		}

		err = h.service.Delete(recordId)
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
			"Group deleted successfully",
		)
	}
}

func (h *DistributionGroupHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusMethodNotAllowed,
			false,
			fmt.Sprintf("Only GET method allowed"),
		)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	distributionGroups, err := h.service.GetList(user.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	handlers.ReturnAppBaseResponse(
		w,
		http.StatusOK,
		true,
		distributionGroups,
	)
}

func (h *DistributionGroupHandler) Run(writer http.ResponseWriter, request *http.Request) {
	conn, err := grpc.Dial("vkspam_parser:10001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewParserClient(conn)

	req := &pb.Request{
		Message: "TEst m",
	}

	ctx, cancel := context.WithTimeout(request.Context(), time.Second*10)
	defer cancel()

	res, err := client.Do(ctx, req)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err.Error())
	} else {

		fmt.Println(res)
	}
}
