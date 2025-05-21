package response

type ReadColumnResponse struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Tasks []TaskBrief `json:"tasks"`
}
