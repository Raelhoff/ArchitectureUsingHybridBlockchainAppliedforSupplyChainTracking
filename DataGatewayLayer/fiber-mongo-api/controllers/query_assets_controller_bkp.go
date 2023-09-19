package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mongo-api/models"
	"mongo-api/packge"
	"mongo-api/responses"
	"time"
	//	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func QueryAssets(c *fiber.Ctx) error {
	// Extrair o seletor do corpo da solicitação
	selector := c.Body()

	// Obter a instância do contrato
	contractSmartAlert := packge.GetContractSmartAlert()

	// Executar a chamada de método no contrato para consultar ativos pelo seletor
	result, err := contractSmartAlert.SubmitTransaction("QueryAssets", string(selector))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to execute transaction: %v", err),
		})
	}

	// Processar o resultado da transação
	var assets []models.QueryAssets
	err = json.Unmarshal(result, &assets)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to decode JSON: %v", err),
		})
	}

	// Retornar a lista de ativos encontrados
	return c.JSON(assets)
}

func CreateAsset(c *fiber.Ctx) error {
	var newAsset models.QueryAssets

	fmt.Println("CreateAlert hYPER")
	
	// Log do corpo da requisição
	bodyBytes := c.Request().Body()
	fmt.Printf("Requisição recebida: %s\n", bodyBytes)

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05") // O layout define o formato desejado

	fmt.Println("Ini - Data e Hora formatadas:", formattedTime)



	if err := c.BodyParser(&newAsset); err != nil {
		fmt.Println("Invalid data format")
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{
			Message: "Invalid data format",
		})
	}

	// Obter a instância do contrato
	contractSmartAlert := packge.GetContractSmartAlert()

	// Executar a chamada de método no contrato para criar um ativo
	resultInvoke, err := contractSmartAlert.SubmitTransaction(
		"CreateAsset",
		newAsset.ID,
		newAsset.IdEdge,
		newAsset.IdNodo,
		newAsset.Type,
		newAsset.Hash,
		newAsset.Date,
		newAsset.Temperatura,
		newAsset.Umidade,
		newAsset.Rele,
		newAsset.Description,
	)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{
			Message: fmt.Sprintf("Failed to commit transaction: %v", err),
		})
	}

	fmt.Println(resultInvoke)
	fmt.Println("Asset Created:", newAsset.Hash)

	// Retorne uma resposta de sucesso
	return c.JSON(responses.SuccessResponse{
		Data: newAsset,
	})
}
