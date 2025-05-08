package response

type ReadUserResponse struct {
    ID       uint           `json:"id"`
    Name     string         `json:"name"`
    Email    string         `json:"email"`
    Projects []ProjectBrief `json:"projects"`
    Tasks    []TaskBrief    `json:"tasks"`
}

type ProjectBrief struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

type TaskBrief struct {
    ID     uint   `json:"id"`
    Name  string `json:"title"`
    Status string `json:"status"`
}