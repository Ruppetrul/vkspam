package distributions

import (
	"net/http"
	"vkspam/middleware"
)

func List(w http.ResponseWriter, r *http.Request) {
	middleware.GetUserFromContext(r.Context())

	//TODO
}
