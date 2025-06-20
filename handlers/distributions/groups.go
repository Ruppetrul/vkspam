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
	Service             services.DistributionGroupService
	distributionService services.DistributionService
}

func NewDistributionGroupHandler() *DistributionGroupHandler {
	db, _ := database.GetDBInstance()

	distributionGroupRepository := repositories.NewDistributionGroupRepository(db.Db)
	distributionGroupService := services.NewDistributionGroupService(distributionGroupRepository)

	distributionRepository := repositories.NewDistributionRepository(db.Db)
	distributionService := services.NewDistributionService(distributionRepository)

	return &DistributionGroupHandler{
		Service:             distributionGroupService,
		distributionService: distributionService,
	}
}

func (h *DistributionGroupHandler) Group(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())

	if 0 == user.Id {
		handlers.ReturnAppBaseResponse(
			w,
			http.StatusBadRequest,
			false,
			fmt.Sprintf("User id undefined"),
		)

		return
	}

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

		model, err := h.Service.Get(recordId)
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
				Id:                  model.Id,
				Name:                model.Name,
				Description:         model.Description,
				UserId:              model.UserId,
				Sex:                 model.Sex,
				Distributions:       *distributions,
				OnlyBirthdayToday:   model.OnlyBirthdayToday,
				OnlyBirthdayFriends: model.OnlyBirthdayFriends,
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

		onlyBirthdayFriends := r.FormValue("only_birthday_friends")
		if len(onlyBirthdayFriends) < 1 {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"Missing required parameter 'only_birthday_friends'",
			)

			return
		}

		onlyBirthdayTodayBool, _ := strconv.ParseBool(onlyBirthdayToday)
		onlyBirthdayFriendsBool, _ := strconv.ParseBool(onlyBirthdayFriends)

		if onlyBirthdayTodayBool && onlyBirthdayFriendsBool {
			handlers.ReturnAppBaseResponse(
				w,
				http.StatusBadRequest,
				false,
				"only_birthday_today and only_birthday_friends cannot be set at the same time",
			)

			return
		}

		newDistributionGroup := models.DistributionGroup{
			Name:                r.Form.Get("name"),
			Description:         r.Form.Get("description"),
			UserId:              user.Id,
			Sex:                 sexInt,
			OnlyBirthdayToday:   onlyBirthdayTodayBool,
			OnlyBirthdayFriends: onlyBirthdayFriendsBool,
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

		id, err := h.Service.Save(newDistributionGroup)

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

		err = h.Service.Delete(recordId)
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
	distributionGroups, err := h.Service.GetList(user.Id)

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

	if -1 < (*GetProgress(groupIdInt)).Progress {
		http.Error(writer, "Already running", http.StatusBadRequest)
		return
	}

	distributions, err := distributionHandler.service.GetListByGroup(groupIdInt)
	if err != nil {
		return
	}

	groupService := NewDistributionGroupHandler()
	distributionGroup, err := groupService.Service.Get(groupIdInt)
	if err != nil {
		handlers.ReturnAppBaseResponse(
			writer,
			http.StatusInternalServerError,
			false,
			err.Error(),
		)

		return
	}

	if len(*distributions) < 1 {
		http.Error(writer, "Рассылки не найдены", http.StatusBadRequest)
		return
	}

	today := time.Now()
	if distributionGroup.LastProcessing.Year() == today.Year() &&
		distributionGroup.LastProcessing.Month() == today.Month() &&
		distributionGroup.LastProcessing.Day() == today.Day() {
		http.Error(writer, "Today already was processing.", http.StatusInternalServerError)
		return
	}

	conn, err := NewConnection()

	if err != nil {
		http.Error(writer, "Parser not available. Please try later.", http.StatusInternalServerError)
		return
	}

	(*distributionGroup).LastProcessing = time.Now()
	_, err = h.Service.Save(*distributionGroup)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	UpdateProgress(groupIdInt, 0, "")
	go process(distributionGroup, distributions, conn)

	err = json.NewEncoder(writer).Encode("Success")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func process(distributionGroup *models.DistributionGroup, distributions *[]models.Distribution, conn *grpc.ClientConn) {
	client := pb.NewParserClient(conn)

	for position, distribution := range *distributions {
		processDistribution(&client, &distribution, distributionGroup, position, len(*distributions))
		time.Sleep(3 * time.Second)
		fmt.Printf("Старутет обработка %d", position)
	}

	DeleteProgress(distributionGroup.Id)
}

func processDistribution(
	client *pb.ParserClient,
	distribution *models.Distribution,
	distributionGroup *models.DistributionGroup,
	position int,
	distributionsCount int,
) {
	var birthdayFilter = ""
	if true == distributionGroup.OnlyBirthdayToday || true == distributionGroup.OnlyBirthdayFriends {
		currentDate := time.Now()
		birthdayFilter = currentDate.Format("2.1")
	}

	req := &pb.ParsePublicRequest{
		VkToken:   os.Getenv("SYSTEM_VK_TOKEN"),
		PublicUrl: distribution.Url,
		Filters: &pb.MemberFilters{
			Birthday:        birthdayFilter,               //8.2
			Sex:             int32(distributionGroup.Sex), //1-woman , 2-man
			BirthdayFriends: distributionGroup.OnlyBirthdayFriends,
		},
		Message: distributionGroup.Description,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	stream, err := (*client).ParsePublic(ctx, req)

	defer stream.CloseSend() //TODO handle err

	if err != nil {
		DeleteProgress(distributionGroup.Id)
		fmt.Println("error")
		fmt.Println(err.Error())
		return
	}

	for {
		progress, err := stream.Recv()
		if err != nil {
			DeleteProgress(distributionGroup.Id)
			fmt.Println("error")
			return
		}

		message := fmt.Sprintf("%d/%d Рассылка %s", position+1, distributionsCount, distribution.Name)
		UpdateProgress(distribution.GroupId, int(progress.Progress), message)
		fmt.Printf("Прогресс задачи %s: %d%%, Сообщение: %s\n", progress.TaskId, progress.Progress, progress.Message)

		if progress.Progress == 100 {
			break
		}
	}
}
