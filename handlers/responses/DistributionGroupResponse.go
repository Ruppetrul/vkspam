package responses

import "vkspam/models"

type DistributionGroupResponse struct {
	Id            int                   `json:"id"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	UserId        int                   `json:"user_id"`
	Distributions []models.Distribution `json:"distributions"`
}
