package main

import (
	"context"
	"fmt"
	"log"
	"mongo-api/configs" //add this
	"mongo-api/packge"
	"mongo-api/routes" //add this
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var alertHashBackCollection *mongo.Collection = configs.GetCollection(configs.DB, "alerts")
var alertHashFifoBackCollection *mongo.Collection = configs.GetCollection(configs.DB, "alertsHashFiFo")
var contractSmartAlertHash = packge.ReadContractSmartAlertHash()
// Variável global para rastrear se a fila está vazia
var isQueueEmpty bool

// Mutex para garantir acesso exclusivo à variável isQueueEmpty
var queueEmptyMutex sync.Mutex

func InteractWithSmartContractBack(hash string) error {
	// Define um arquivo de log para registrar o tempo da função SubmitTransaction
	// Abra o arquivo de log em modo de escrita apendicular para concatenar os dados
	logFile, err := os.OpenFile("transaction_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Use um logger para gravar no arquivo de log
	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	// Define um erro variável para rastrear quaisquer erros durante o processo
	var interactionError error

	// Use um canal para aguardar a conclusão da interação
	done := make(chan bool)

	go func() {
		defer func() {
			// Sinaliza que a interação está concluída
			done <- true
		}()

		//contractSmartAlertHash := packge.ReadContractSmartAlertHash()

		startTime := time.Now() // Hora de início da execução da função

		// Execute o método do contrato usando o hash fornecido
		// Execute o método do contrato (ajuste os parâmetros conforme necessário)
		_, err := contractSmartAlertHash.SubmitTransaction("CreateAsset", hash)
		if err != nil {
			interactionError = err
			return
		}

		endTime := time.Now() // Hora de término da execução da função

		// Calcule o tempo de resposta
		responseTime := endTime.Sub(startTime)

		// Registre o tempo de resposta no arquivo de log
		logger.Printf("Data: %s - Tempo de resposta: %s\n", startTime.Format("2006-01-02 15:04:05"), responseTime)

		// ... (resto do código)

	}()

	// Aguarde a conclusão da interação
	<-done

	return interactionError
}

func executePythonScript() {
	cmd := exec.Command("python3", "geraResult.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Erro ao executar o script Python: %v\n", err)
	}
}

var alertProcessingComplete = make(chan bool)

func ProcessFifoAndInsertToBlockchain() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		findOptions := options.FindOne()
		findOptions.Sort = bson.M{"_id": 1}
		var hashDocument bson.M
		err := alertHashFifoBackCollection.FindOne(ctx, bson.M{}, findOptions).Decode(&hashDocument)

		queueEmptyMutex.Lock()
		isQueueEmpty = err != nil
		queueEmptyMutex.Unlock()

		if err == nil {
			hash, ok := hashDocument["hash"].(string)
			if ok {
				go func() {
					err := InteractWithSmartContractBack(hash)
					if err != nil {
						fmt.Printf("Error processing alert: %v\n", err)
						errMessage := err.Error()
						if strings.Contains(errMessage, "already exists") {
							fmt.Println("Asset already exists error detected.")
							// Remove o hash processado da coleção FIFO
							_, deleteErr := alertHashFifoBackCollection.DeleteOne(ctx, bson.M{"hash": hash})
							if deleteErr != nil {
								fmt.Printf("Error deleting hash from FIFO collection: %v\n", deleteErr)
							}
						} else {
							fmt.Println("Other error occurred.")
						}
					} else {
						_, deleteErr := alertHashFifoBackCollection.DeleteOne(ctx, bson.M{"hash": hash})
						if deleteErr != nil {
							fmt.Printf("Error deleting hash from FIFO collection: %v\n", deleteErr)
						}
					}

					// Sinalize que o processamento da fila de alertas está completo
					alertProcessingComplete <- true
				}()
				// Aguarde o processamento do alerta ser concluído antes de passar para o próximo hash
				<-alertProcessingComplete
			}
		} else {
			time.Sleep(time.Second)
			// Verifique se a fila está vazia e execute o script Python se estiver
			queueEmptyMutex.Lock()
			if isQueueEmpty {
				queueEmptyMutex.Unlock()
				// Execute o script Python quando a fila estiver vazia
				executePythonScript()
			} else {
				queueEmptyMutex.Unlock()
			}
		}
	}
}

func StartAlertProcessing() {
	go ProcessFifoAndInsertToBlockchain()
}

func main() {
	app := fiber.New()

	// Configurar as opções de CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:9090", // Origens permitidas
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Content-Type",
	}))

	//run database
	configs.ConnectDB()

	//routes
	routes.DeviceRoute(app) //add this
	routes.AlertaRoute(app)
	routes.DeviceEdgeRoute(app)
	routes.DeviceNodeRoute(app)
	routes.EdgeActiveRoute(app)
	routes.NodeActiveRoute(app)

	// Inicialize o contrato
	packge.InitializeContract()
	routes.QueryRoute(app)

	packge.InitializeContractSmartHash()
	routes.AlertaHashRoute(app)

	routes.AlertaHashBackRoute(app)

	StartAlertProcessing()

	app.Listen(":80")
}
