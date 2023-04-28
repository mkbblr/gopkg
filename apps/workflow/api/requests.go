package api

// Workflow create request
// @Description Workflow create request
type CreateWorkflowRequest struct {
	Description string `json:"desc" example:"configuration change workflow"`
} //@name response.Body(CreateWorkflowRequest)

// Workflow delete request
// @Description Workflow delete request
type DeleteWorkflowRequest struct {
	ID int `json:"id" example:"1"`
} //@name response.Body(DeleteWorkflowRequest)
