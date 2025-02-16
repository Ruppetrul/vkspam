package distributions

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Ruppetrul/vkspam_proto/gen/parser"
	"google.golang.org/grpc"
	"net/http"
	"os"
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
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				err.Error(),
			)

			return
		}

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

		if model == nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				true,
				fmt.Sprintf("Ошибка чтения записи из базы'"),
			)
		}

		handlers.ReturnAppBaseResponse(
			w,
			http.StatusOK,
			true,
			responses.DistributionGroupResponse{
				Id:                model.Id,
				Name:              model.Name,
				Description:       model.Description,
				UserId:            model.UserId,
				Sex:               model.Sex,
				Distributions:     *distributions,
				OnlyBirthdayToday: model.OnlyBirthdayToday,
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

		sex := r.FormValue("sex")
		if len(sex) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'sex'",
			)

			return
		}

		sexInt, err := strconv.Atoi(sex)
		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf("Ошибка преобразования: %s", err),
			)

			return
		}

		onlyBirthdayToday := r.FormValue("only_birthday_today")
		if len(onlyBirthdayToday) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'only_birthday_today'",
			)

			return
		}

		onlyBirthdayTodayBool, _ := strconv.ParseBool(onlyBirthdayToday)

		newDistributionGroup := models.DistributionGroup{
			Name:              r.Form.Get("name"),
			Description:       r.Form.Get("description"),
			UserId:            user.Id,
			Sex:               sexInt,
			OnlyBirthdayToday: onlyBirthdayTodayBool,
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

		id, err := h.service.Save(newDistributionGroup)

		if err != nil {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				fmt.Sprintf("Save error: %s", err.Error()),
			)
			return
		}

		if id == 0 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusInternalServerError,
				false,
				"error check id",
			)
			return
		}
		newDistributionGroup.Id = id
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusOK,
			true,
			newDistributionGroup,
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
	if request.Method != http.MethodPost {
		http.Error(writer, "Only Post", http.StatusMethodNotAllowed)
		return
	}

	groupId := request.FormValue("group_id")
	if len(groupId) < 1 {
		http.Error(writer, "Missing required parameter 'group_id'", http.StatusBadRequest)
		return
	}

	distributionHandler := NewDistributionHandler()

	groupIdInt, _ := strconv.Atoi(groupId)

	if 0 == groupIdInt {
		http.Error(writer, "Missing required parameter 'group_id'", http.StatusBadRequest)
		return
	}

	if -1 < GetProgress(groupIdInt) {
		http.Error(writer, "Already running", http.StatusBadRequest)
		return
	}

	distributions, err := distributionHandler.service.GetListByGroup(groupIdInt)
	if err != nil {
		return
	}

	if len(*distributions) < 1 {
		http.Error(writer, "Рассылки не найдены", http.StatusBadRequest)
	}

	conn, err := NewConnection()

	if err != nil {
		http.Error(writer, "Parser not available. Please try later.", http.StatusInternalServerError)
		return
	}
	UpdateProgress(groupIdInt, 0)
	go process(distributions, conn)

	err = json.NewEncoder(writer).Encode("Success")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func process(distributions *[]models.Distribution, conn *grpc.ClientConn) {
	client := pb.NewParserClient(conn)

	for _, distribution := range *distributions {
		processDistribution(&client, &distribution)
	}
}

func processDistribution(client *pb.ParserClient, distribution *models.Distribution) {
	req := &pb.ParsePublicRequest{
		VkToken:   os.Getenv("SYSTEM_VK_TOKEN"),
		PublicUrl: distribution.Url,
		Filters: &pb.MemberFilters{
			Birthday: "15.2", //8.2
			Sex:      2,      //1-woman , 2-man
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := (*client).ParsePublic(ctx, req)
	if err != nil {
		UpdateProgress(distribution.GroupId, -2)
		fmt.Println("error")
		fmt.Println(err.Error())
		return
	}

	for {
		progress, err := stream.Recv()
		if err != nil {
			UpdateProgress(distribution.GroupId, -2)
			fmt.Println("error")
			return
		}

		UpdateProgress(distribution.GroupId, int(progress.Progress))
		fmt.Printf("Прогресс задачи %s: %d%%, Сообщение: %s\n", progress.TaskId, progress.Progress, progress.Message)

		if progress.Progress == 100 {
			break
		}
	}
}
