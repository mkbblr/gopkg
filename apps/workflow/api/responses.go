package api

// Error Response
// @Description Error response body
type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode" example:"666"`
	ErrorMessage string `json:"errorMessage" example:"workflow not found"`
} //@name response.Body(ErrorResponse)

// Workflow API Response
// @Description Workflow API response model
type WorkflowResponse struct {
	ID        int      `json:"id" example:"1"`
	Name      string   `json:"name" example:"account name"`
	PhotoUrls []string `json:"photo_urls" example:"http://test/image/1.jpg,http://test/image/2.jpg"`
} //@name response.Body(WorkflowResponse)
