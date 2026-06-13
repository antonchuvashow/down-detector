package apiinspector

type CreateRequest struct {
	Inspector Inspector `json:"inspector" binding:"required"`
}
