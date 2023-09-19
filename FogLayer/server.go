package main

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"go.mongodb.org/mongo-driver/bson" //add this
	"go.mongodb.org/mongo-driver/mongo/options"

	//	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net"
	"os"
	"protofiles/server/configs" //add this
	"protofiles/server/models"
	"protofiles/server/packge" //add this
	"time"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/reflection"
)

type server struct {
	packge.UnimplementedLoraTransactionServer
	packge.UnimplementedListLoraTransactionServer
	packge.UnimplementedAlertServer
	packge.UnimplementedKeepAliveServer
}

const (
	RequestRegisteredDevices = 1
	Ping                     = 2
	Test                     = 3
)

var alertHashBackCollection *mongo.Collection = configs.GetCollection(configs.DB, "alerts")
var alertHashFifoBackCollection *mongo.Collection = configs.GetCollection(configs.DB, "alertsHashFiFo")
var devicesCollection *mongo.Collection = configs.GetCollection(configs.DB, "devices")

var contractSmart *gateway.Contract = ReadContractSmart()
var contractSmartAlert *gateway.Contract = ReadContractSmartAlert()

func InteractWithSmartContractBack(hash string) error {
	// Define an error variable to track any errors during the process
	var interactionError error

	// Use a channel to wait for the interaction to complete
	done := make(chan bool)

	go func() {
		defer func() {
			// Signal that the interaction is done
			done <- true
		}()

		contractSmartAlertHash := packge.ReadContractSmartAlertHash()

		// Execute the contract method using the provided hash
		// Execute the contract method (adjust parameters as needed)
		_, err := contractSmartAlertHash.SubmitTransaction("CreateAsset", hash)
		if err != nil {
			interactionError = err
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Update the alert in the database to mark as processed
		result, err := alertHashBackCollection.UpdateOne(
			ctx,
			bson.M{"hash": hash},
			bson.M{"$set": bson.M{"processed": true}},
		)
		if err != nil {
			fmt.Printf("Failed to update alert in database: %v", err)
			interactionError = err
			return
		}
		fmt.Println(result.ModifiedCount)
	}()

	// Wait for the interaction to complete
	<-done

	return interactionError
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
							// Remove the processed hash from the FIFO collection
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

					// Signal that alert processing is complete
					alertProcessingComplete <- true
				}()

				// Wait for the alert processing to complete before moving to the next hash
				<-alertProcessingComplete
			}
		} else {
			time.Sleep(time.Second)
		}
	}
}

func StartAlertProcessing() {
	go ProcessFifoAndInsertToBlockchain()
}

func (s *server) SendAlert(ctx context.Context, in *packge.AlertRequest) (*packge.AlertResponse, error) {
	// Create a new alert request object
	newAlert := models.AlertRequest{
		Id:          in.Id,
		IdEdge:      in.IdEdge,
		IdNodo:      in.IdNodo,
		Hash:        in.Hash,
		Type:        in.Type,
		Date:        in.Date,
		Temperatura: in.Temperatura,
		Umidade:     in.Umidade,
		Rele:        in.Rele,
		Description: in.Description,
	}

	// Encode and hash the alert data
	dataOut := EncodeToBytes(newAlert)
	sha512Hash := sha512.Sum512([]byte(dataOut))
	newAlert.Hash = hex.EncodeToString(sha512Hash[:])

	// Create a new device alert object
	newDeviceAlert := models.DeviceAlert{
		IdEdge:      in.IdEdge,
		IdNodo:      in.IdNodo,
		Hash:        newAlert.Hash,
		Type:        in.Type,
		Date:        in.Date,
		Temperatura: in.Temperatura,
		Umidade:     in.Umidade,
		Rele:        in.Rele,
		Description: in.Description,
	}

	// Record the start time
	startTime := time.Now()
	formattedStartTime := startTime.Format("2006-01-02 15:04:05")
	fmt.Println("Start - Formatted Date and Time:", formattedStartTime)

	// Create a context with a timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Insert the new device alert into alertHashBackCollection
	if _, err := alertHashBackCollection.InsertOne(dbCtx, newDeviceAlert); err != nil {
		fmt.Printf("Error inserting device alert: %v\n", err)
		return &packge.AlertResponse{Confirmation: false}, err
	}

	// Insert the alert hash into alertHashFifoBackCollection
	hashDocument := bson.M{"hash": newDeviceAlert.Hash}
	if _, err := alertHashFifoBackCollection.InsertOne(dbCtx, hashDocument); err != nil {
		fmt.Printf("Error inserting alert hash: %v\n", err)
		return &packge.AlertResponse{Confirmation: false}, err
	}

	// Record the end time
	endTime := time.Now()
	formattedEndTime := endTime.Format("2006-01-02 15:04:05")
	fmt.Println("End - Formatted Date and Time:", formattedEndTime)

	// Return a successful confirmation
	return &packge.AlertResponse{Confirmation: true}, nil
}

func returnListDeviceNode() []models.DeviceNodo {

	var deviceEdgeCollection *mongo.Collection = configs.GetCollection(configs.DB, "devices_nodo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.DeviceNodo
	defer cancel()

	results, err := deviceEdgeCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return devices
	}

	if err = results.All(ctx, &devices); err != nil {
		log.Fatal(err)
		return devices
	}
	fmt.Println(devices)
	return devices
}

func returnListDeviceEdge() []models.DeviceEdge {

	var deviceEdgeCollection *mongo.Collection = configs.GetCollection(configs.DB, "devices_edge")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var devices []models.DeviceEdge
	defer cancel()

	results, err := deviceEdgeCollection.Find(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		return devices
	}

	if err = results.All(ctx, &devices); err != nil {
		log.Fatal(err)
		return devices
	}
	fmt.Println(devices)
	return devices
}

func convertJSONStringToList(jsonString string) ([]*packge.DeviceActive, error) {
	// Cria uma slice para armazenar os dados decodificados
	var list []*packge.DeviceActive

	// Decodifica a string JSON na slice list
	err := json.Unmarshal([]byte(jsonString), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func convertJSONStringToListEdge(jsonString string) ([]*packge.EdgeActive, error) {
	// Cria uma slice para armazenar os dados decodificados
	var list []*packge.EdgeActive

	// Decodifica a string JSON na slice list
	err := json.Unmarshal([]byte(jsonString), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func saveOrUpdateList(collection *mongo.Collection, deviceList []*packge.DeviceActive) error {
	// Cria um slice de modelos WriteModel para inserção/atualização em massa
	var writes []mongo.WriteModel
	for _, device := range deviceList {
		// Cria um filtro para encontrar o documento pelo campo "ID"
		filter := bson.M{"idnodo": device.IDNodo}

		// Remove o campo "_id" do documento para evitar o erro de imutabilidade
		updateDoc := bson.M{
			"$set": bson.M{
				"idedge":       device.IDEdge,
				"name":         device.Name,
				"isactive":     device.IsActive,
				"notification": device.Notification,
				"updateat":     device.UpdatedAt.Local().String(),
			},
		}

		// Cria um modelo WriteModel para a operação de upsert (inserção ou atualização)
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateDoc).SetUpsert(true)

		// Adiciona o modelo ao slice de writes
		writes = append(writes, model)
	}

	// Realiza a operação de inserção/atualização em massa no MongoDB
	_, err := collection.BulkWrite(context.Background(), writes)
	if err != nil {
		return err
	}

	return nil
}

func saveOrUpdateListEDGE(collection *mongo.Collection, deviceList []*packge.EdgeActive) error {
	// Cria um slice de modelos WriteModel para inserção/atualização em massa
	var writes []mongo.WriteModel
	for _, device := range deviceList {
		// Cria um filtro para encontrar o documento pelo campo "ID"
		filter := bson.M{"idedge": device.IDEdge}

		// Remove o campo "_id" do documento para evitar o erro de imutabilidade
		updateDoc := bson.M{
			"$set": bson.M{
				"name":         device.Name,
				"isactive":     device.IsActive,
				"notification": device.Notification,
				"updateat":     device.UpdatedAt.Local().String(),
			},
		}

		// Cria um modelo WriteModel para a operação de upsert (inserção ou atualização)
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(updateDoc).SetUpsert(true)

		// Adiciona o modelo ao slice de writes
		writes = append(writes, model)
	}

	// Realiza a operação de inserção/atualização em massa no MongoDB
	_, err := collection.BulkWrite(context.Background(), writes)
	if err != nil {
		return err
	}

	return nil
}

func processKeepalive(jsonString string, date string) string {
	fmt.Println(jsonString)

	// Converte a string JSON em uma lista de DeviceActive
	listDeviceActive, err := convertJSONStringToList(jsonString)
	if err != nil {
		fmt.Println("Erro ao converter JSON para lista:", err)
		return "Erro ao converter JSON para lista"
	}

	// Imprime a lista de DeviceActive
	for _, device := range listDeviceActive {
		fmt.Printf("IDEdge: %d, IDNodo: %d, Name: %s, IsActive: %v, Notification: %v, UpdatedAt: %s\n",
			device.IDEdge, device.IDNodo, device.Name, device.IsActive, device.Notification, device.UpdatedAt)
	}

	// Obtem a coleção de dispositivos no MongoDB
	var deviceCollection *mongo.Collection = configs.GetCollection(configs.DB, "nodo_active")

	// Salva a lista de dispositivos no MongoDB
	err = saveOrUpdateList(deviceCollection, listDeviceActive)
	if err != nil {
		fmt.Println("Erro ao salvar lista no MongoDB:", err)
		return "Erro ao salvar lista no MongoDB"
	}

	fmt.Println("Lista de dispositivos salva no MongoDB com sucesso!")
	return "Lista de dispositivos salva no MongoDB com sucesso!"
}

func updatesEdgeMonitoringList(devicesEdge []models.DeviceEdge) {
	// Verifica se tem comunicacao com os Edges
	var _listEdgeActive = make([]*packge.EdgeActive, 0)

	// Retorna todos dispositivosNodo associado a esse EDGE
	if len(devicesEdge) > 0 {

		fmt.Println("---------------- deviceNodo to NewDeviceActive")
		for i, device := range devicesEdge {
			fmt.Println(i, device.IdEdge)
			fmt.Println(i, device.Ip)
			fmt.Println(i, device.User)
			fmt.Println(i, device.Date)
			fmt.Println(i, device.Hash)

			edgeAtivo := packge.NewEdgeActive(int(device.IdEdge), "Edge"+strconv.Itoa(int(i+1)))
			edgeAtivo.SetActive(true)
			//fmt.Println(dispositivoAtivo)

			// Criar lista para verificar se dispositivos estão ativos
			_listEdgeActive = append(_listEdgeActive, edgeAtivo)
		}

		// Access and display the list of global devices
		fmt.Println("----------------------------------------------------------------------------------------")
		fmt.Println("Lista de dispositivos globais (Edge):")
		for _, edge := range _listEdgeActive {
			fmt.Printf("ID: %d, IDEdge: %d, Nome: %s\n", edge.ID, edge.IDEdge, edge.Name)
		}
		fmt.Println("----------------------------------------------------------------------------------------")

	}

	// Atribuir a lista de dispositivos à variável global
	packge.SetEdges(_listEdgeActive)
}

func updateStatusEdge(IDEdge int) {
	// Verifica se tem comunicacao com os Edges
	var listEdgeActive = packge.GetEdges()
	var isFound bool = true
	fmt.Println("---------------- deviceNodo to NewDeviceActive")
	for i, device := range listEdgeActive {

		if device.IDEdge == IDEdge {
			fmt.Println(i, device.IDEdge)
			fmt.Println(i, device.Name)
			fmt.Println(i, device.IsActive)
			fmt.Println(i, device.Notification)
			fmt.Println(i, device.UpdatedAt.Local().String())
			isFound = false

			packge.SetEdgeStatus(IDEdge, true)
			packge.SetEdgeNotif(IDEdge, false)

		}
	}

	if isFound {
		// Não encontrou dispositivo, atualiza lista
		devicesEdge := returnListDeviceEdge()
		// Verifica se tem comunicacao com os Edges
		updatesEdgeMonitoringList(devicesEdge)
	}

	// Access and display the list of global devices
	listEdgeActive = packge.GetEdges()
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println("Lista de dispositivos globais (Edge):")
	for _, edge := range listEdgeActive {
		fmt.Printf("ID: %d, IDEdge: %d, Nome: %s\n", edge.ID, edge.IDEdge, edge.Name)
	}
	fmt.Println("----------------------------------------------------------------------------------------")

}

func (s *server) MakeKeepAlive(ctx context.Context, in *packge.KeepAliveRequest) (*packge.KeepAliveResponse, error) {
	var devicesNodo []models.DeviceNodo
	var devicesNodoAux []models.DeviceNodo

	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("Processa KeepAlive")
	switch in.Type {
	case RequestRegisteredDevices:
		fmt.Println("RequestRegisteredDevices")

		// Atualiza Stutus EDGE
		updateStatusEdge(int(in.IdEdge))

		var devicesEdge = packge.GetEdges()
		for index, element := range devicesEdge {
			if int32(element.IDEdge) == in.IdEdge {
				devicesNodo = returnListDeviceNode()
				for index2, element2 := range devicesNodo {
					if int32(element.IDEdge) == element2.IdEdge {
						devicesNodoAux = append(devicesNodoAux, element2)
					}
					fmt.Println(index2)
				}

				break
			}
			fmt.Println(index)
		}

		b, err := json.Marshal(devicesNodoAux)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return &packge.KeepAliveResponse{}, nil
		}
		fmt.Println(string(b))
		now := time.Now()
		return &packge.KeepAliveResponse{IdEdge: in.IdEdge, Type: RequestRegisteredDevices, Date: now.String(), Msg: string(b)}, nil

	case Ping:
		fmt.Println("Ping")
		now := time.Now()
		return &packge.KeepAliveResponse{IdEdge: in.IdEdge, Type: RequestRegisteredDevices, Date: now.String(), Msg: processKeepalive(in.Msg, in.Date)}, nil

	case Test:
		fmt.Println("Test.")
	default:
		fmt.Println("Too far away.")
	}

	return &packge.KeepAliveResponse{}, nil
}

func EncodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}

func ReadWallet() *gateway.Wallet {

	w, err := gateway.NewFileSystemWallet("./wallets")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	if !w.Exists("Admin") {
		fmt.Println("Failed to get Admin from wallet")
		os.Exit(1)
	}
	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("ReadWallet to wallets")
	return w
}

func ReadContractSmart() *gateway.Contract {
	// Client instance
	var wallet *gateway.Wallet = ReadWallet()

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("./connection.json")),
		gateway.WithIdentity(wallet, "Admin"),
	)

	if err != nil {
		fmt.Println("-------------------------------------------------------------------------------------")
		fmt.Printf("Failed to connect: %v", err)
	} else {
		fmt.Println("-------------------------------------------------------------------------------------")
		fmt.Println("Connect fabric")
	}

	if gw == nil {
		fmt.Println("-------------------------------------------------------------------------------------")
		fmt.Println("Failed to create gateway")
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %v", err)
	}

	//var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	c := network.GetContract("smarthash")

	return c
}

func ReadContractSmartAlert() *gateway.Contract {
	// Client instance
	var wallet *gateway.Wallet = ReadWallet()

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("./connection.json")),
		gateway.WithIdentity(wallet, "Admin"),
	)

	if err != nil {
		fmt.Println("-------------------------------------------------------------------------------------")
		fmt.Printf("Failed to connect: %v", err)
	} else {
		fmt.Println("-------------------------------------------------------------------------------------")
		fmt.Println("Connect fabric")
	}

	if gw == nil {
		fmt.Println("-------------------------------------------------------------------------------------")
		fmt.Println("Failed to create gateway")
	}

	network, err := gw.GetNetwork("mychannel")
	if err != nil {
		fmt.Printf("Failed to get network: %v", err)
	}

	//var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	c := network.GetContract("smartalert")

	return c
}

// [ctx] is used by the goroutines to interact with GRPC
// [in] is the type of TransactionRequest
/*
	This function signature matches the service that we mentioned in the protobuf
*/
func (s *server) MakeTransaction(ctx context.Context, in *packge.LoraRequest) (*packge.LoraResponse, error) {
	/*
		IdDevice:     1234,
			Id:           456,
			Input1:       0,
			Input2:       0,
			Output:       0,
			AlarmBattery: true,
			AlarmPower:   true,
			SensorError:  true,
			Sensors: []*pb.Device_Sensor{
				{Type: 1, Value: 30},
				{Type: 2, Value: 15},
			},
			LastUpdated: timestamppb.Now(),
		}
	*/
	// Business logic will come here
	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("IdEdge: ", in.IdEdge)
	fmt.Println("IdNodo: ", in.IdNodo)
	fmt.Println("Input1: ", in.Input1)
	fmt.Println("Input2: ", in.Input2)
	fmt.Println("Output: ", in.Output)
	fmt.Println("AlarmBattery: ", in.AlarmBattery)
	fmt.Println("AlarmPower: ", in.AlarmPower)
	fmt.Println("SensorError: ", in.SensorError)
	fmt.Println("Sensors: ", in.Sensors)
	fmt.Println("LastUpdated: ", in.LastUpdated)

	uuid.SetRand(nil)

	packge.SetEdgeStatus(int(in.IdEdge), true)

	// invoke
	fmt.Println("\n\n --- invoke ---")
	//var temp float32 = float32(in.Sensors[0].Value)
	//var humid float32 = float32(in.Sensors[1].Value)
	var stemp string = fmt.Sprintf("%f", in.Sensors[0].Value)
	var sUmidade string = fmt.Sprintf("%f", in.Sensors[1].Value)

	//var stemp string = String(in.Sensors[0].Value)
	fmt.Println(in.LastUpdated)
	fmt.Println(stemp)
	fmt.Println(sUmidade)
	fmt.Println(strconv.FormatBool(in.AlarmPower))
	//	var stemp string = strconv.FormatFloat(fmt.Sprintf(temp), 'E', -1, 32)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loc, _ := time.LoadLocation("America/Sao_Paulo")
	t := time.Now().In(loc)

	newDevice := models.Devices{
		Hash:          "Never",
		IdEdge:        in.IdEdge,
		IdNodo:        in.IdNodo,
		Input1:        in.Input1,
		Input2:        in.Input2,
		Output:        in.Output,
		Alarm_battery: in.AlarmBattery,
		Alarm_power:   in.AlarmPower,
		Sensor_error:  in.SensorError,
		Temperatura:   in.Sensors[0].Value,
		Umidade:       in.Sensors[1].Value,
		CreatedAt:     in.LastUpdated,
		UpdatedAt:     t.Format(time.RFC3339),
	}

	dataOut := EncodeToBytes(newDevice)
	//fmt.Println([]byte(dataOut))
	sha_512 := sha512.New()
	// sha from a byte array
	sha_512.Write([]byte(dataOut))
	fmt.Printf("sha512: %x\n", sha_512.Sum(nil))
	newDevice.Hash = hex.EncodeToString(sha_512.Sum(nil))

	_, err := devicesCollection.InsertOne(ctx, newDevice)
	if err != nil {
		response := map[string]string{
			"error": fmt.Sprintf("Failed to create alert: %v", err),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return &packge.LoraResponse{Confirmation: false}, err
	}

	// Returning a response of type Transaction Response
	return &packge.LoraResponse{Msg: "Ok", Confirmation: true}, nil
}

func (s *server) MakeListTransaction(ctx context.Context, in *packge.ListLoraRequest) (*packge.ListLoraResponse, error) {
	/*
		IdDevice:     1234,
			Id:           456,
			Input1:       0,
			Input2:       0,
			Output:       0,
			AlarmBattery: true,
			AlarmPower:   true,
			SensorError:  true,
			Sensors: []*pb.Device_Sensor{
				{Type: 1, Value: 30},
				{Type: 2, Value: 15},
			},
			LastUpdated: timestamppb.Now(),
		}
	*/
	// Business logic will come here
	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("\n\n ------ List: ", len(in.ListDevice))
	for i := 0; i < len(in.ListDevice); i++ {
		fmt.Println("IdEdge: ", in.ListDevice[i].IdEdge)
		fmt.Println("v: ", in.ListDevice[i].IdNodo)
		fmt.Println("Input1: ", in.ListDevice[i].Input1)
		fmt.Println("Input2: ", in.ListDevice[i].Input2)
		fmt.Println("Output: ", in.ListDevice[i].Output)
		fmt.Println("AlarmBattery: ", in.ListDevice[i].AlarmBattery)
		fmt.Println("AlarmPower: ", in.ListDevice[i].AlarmPower)
		fmt.Println("SensorError: ", in.ListDevice[i].SensorError)
		fmt.Println("Sensors: ", in.ListDevice[i].Sensors)
		fmt.Println("LastUpdated: ", in.ListDevice[i].LastUpdated)

	}

	// Returning a response of type Transaction Response
	return &packge.ListLoraResponse{Confirmation: true}, nil
}

func SalveAlert(idEdge int, tipo int32, date string, description string) (string, error) {
	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("Salva Alerta:")
	fmt.Println(idEdge)
	fmt.Println(tipo)
	fmt.Println(date)
	fmt.Println(description)

	newAlert := models.AlertRequest{
		Id:          int32(idEdge), // Assuming Id in the AlertRequest struct corresponds to idEdge
		IdEdge:      int32(idEdge),
		Hash:        "",
		Type:        tipo,
		Date:        date,
		Description: description,
	}

	dataOut := EncodeToBytes(newAlert)
	sha512Hash := sha512.Sum512(dataOut)
	newAlert.Hash = hex.EncodeToString(sha512Hash[:])

	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04:05")
	fmt.Println("Ini - Data e Hora formatadas:", formattedTime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := alertHashBackCollection.InsertOne(ctx, newAlert)
	if err != nil {
		response := map[string]string{
			"error": fmt.Sprintf("Failed to create alert: %v", err),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return string(responseJSON), err
	}

	hashDocument := bson.M{"hash": newAlert.Hash}
	_, err = alertHashFifoBackCollection.InsertOne(ctx, hashDocument)
	if err != nil {
		response := map[string]string{
			"error": fmt.Sprintf("Failed to create alert hash in FIFO collection: %v", err),
		}
		responseJSON, _ := json.Marshal(response)
		fmt.Println(string(responseJSON))
		return string(responseJSON), err
	}

	currentTime2 := time.Now()
	formattedTime2 := currentTime2.Format("2006-01-02 15:04:05")
	fmt.Println("Fim - Data e Hora formatadas:", formattedTime2)

	return "Alerta cadastro com Sucesso", nil
}

func checkEdgeIsActive() {
	for range time.Tick(time.Second * 60) {
		for _, d := range packge.GetEdges() {
			if d.IsActive == false && d.Notification == false {
				fmt.Printf("Dispositivo (Edge) %s (ID: %d) está inativo.\n", d.Name, d.IDEdge)

				result, err := SalveAlert(d.IDEdge, 6, time.Now().Local().String(), "Dispositivo (Edge)está inativo")

				if err != nil {
					fmt.Println("Erro ao salvar o alerta:", err)
				}
				packge.SetEdgeNotif(d.IDEdge, true)
				log.Println(result)

			}
		}
		updateStateEdges()
	}
}

func updateStateEdges() string {

	listEdgeActive := packge.GetEdges()

	// Imprime a lista de DeviceActive
	for _, device := range listEdgeActive {
		fmt.Printf("IDEdge: %d, Name: %s, IsActive: %v, Notification: %v, UpdatedAt: %s\n",
			device.IDEdge, device.Name, device.IsActive, device.Notification, device.UpdatedAt)
	}

	// Obtem a coleção de dispositivos no MongoDB
	var deviceCollection *mongo.Collection = configs.GetCollection(configs.DB, "edge_active")

	// Salva a lista de dispositivos no MongoDB
	err := saveOrUpdateListEDGE(deviceCollection, listEdgeActive)
	if err != nil {
		fmt.Println("Erro ao salvar lista no MongoDB:", err)
		return "Erro ao salvar lista no MongoDB"
	}

	fmt.Println("Lista de dispositivos salva no MongoDB com sucesso!")
	return "Lista de dispositivos salva no MongoDB com sucesso!"
}

func main() {

	StartAlertProcessing()

	devicesEdge := returnListDeviceEdge()
	// Verifica se tem comunicacao com os Edges
	updatesEdgeMonitoringList(devicesEdge)
	go checkEdgeIsActive()

	// NewServer creates a gRPC server which has no service registered and has not started
	// to accept requests yet.
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", ":8017")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// We are making use of the function that compiled proto made for us to register
	// our GRPC server so that the clients can make use of the functions tide to our
	// server remotely via the GRPC server (like MakeTransaction function)

	// The first argument is the grpc server instance
	// The second argument is the service who's methods we want to expose (in our case)
	// we have put it in this program only
	packge.RegisterLoraTransactionServer(s, &server{})
	packge.RegisterListLoraTransactionServer(s, &server{})
	packge.RegisterAlertServer(s, &server{})
	packge.RegisterKeepAliveServer(s, &server{})
	//pb.RegisterAlertServer(s, &server{})

	// Serve accepts incoming connections on the listener lis, creating a new ServerTransport
	// and service goroutine for each. The service goroutines read gRPC requests and then
	// call the registered handlers to reply to them.
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	//	app.Listen(":3033")

}
