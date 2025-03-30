package distributions

import (
	"net/http"
	"strconv"
	"vkspam/handlers"
)

func (h *DistributionGroupHandler) RunProgress(w http.ResponseWriter, r *http.Request) {
	groupId := r.FormValue("group_id")
	if len(groupId) < 1 {
		http.Error(w, "Missing required parameter 'group_id'", http.StatusBadRequest)
		return
	}

	groupIdInt, err := strconv.Atoi(groupId)
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
		true,
		*GetProgress(groupIdInt),
	)
	return
}
