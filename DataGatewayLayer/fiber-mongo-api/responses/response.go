// responses/response.go
package responses

type SuccessResponse struct {
	Data interface{} `json:"data"`
	// Outros campos de resposta de sucesso
}

type ErrorResponse struct {
	Message string `json:"message"`
	// Outros campos de resposta de erro
}
