package response

type ReadProjectResponse struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Tasks       []TaskBrief `json:"tasks"`
}
