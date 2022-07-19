package response

type Message struct {
	Message string `json:"message" binding:"required"`
}
