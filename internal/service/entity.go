package service

type PostRequest struct {
	Title  string `json:"title" validate:"required"`
	Data   string `json:"data"`
	Status string `json:"status"`
}

type RequestWithId struct {
	ID string `validate:"required,intString,min=1"`
}
