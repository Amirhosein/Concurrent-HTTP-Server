package response

type Message struct {
	Message string `json:"message" binding:"required"`
}

type Error struct {
	Error string `json:"error" binding:"required"`
}

type SuccessfulUpload struct {
	FileId string `json:"file_id"`
}
