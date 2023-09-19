package main

import (
	"context"
	initializers "fiber-mongo-api/configs" //add this
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"

	//add this
	//add this
	"fiber-mongo-api/models" //add this
	"fiber-mongo-api/packge"
	"fiber-mongo-api/routes" //add this
	"strconv"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

var debug bool = true   // Defina como true para gravar tempos em arquivo de debug
var idEDGE int32 = 1234 // Defina idEDGE

// Define response data
type Data struct {
	ID []int `json:"id"`
}

func init() {
	initializers.ConnectDB()
	initializers.ConnectDB2()

	// run gRPC
	initializers.ConnectGRPC()
}

func returnListDevice() []models.Devices {

	rows, err := initializers.DB2.Query("SELECT * FROM devices")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var devices = make([]models.Devices, 0)

	for rows.Next() {
		singleDevice := models.Devices{}
		err = rows.Scan(&singleDevice.InsertedID,
			&singleDevice.IdEdge, &singleDevice.IdNodo, &singleDevice.Input1,
			&singleDevice.Input2, &singleDevice.Output,
			&singleDevice.Alarm_battery, &singleDevice.Alarm_power,
			&singleDevice.Sensor_error, &singleDevice.Temperatura,
			&singleDevice.Umidade, &singleDevice.CreatedAt,
			&singleDevice.UpdatedAt)

		if err != nil {
			fmt.Println(err)
		}

		devices = append(devices, singleDevice)
	}

	return devices
}

func deleteDevice(devices models.Devices) {

	result := initializers.DB.Delete(&models.Devices{}, "inserted_id = ?", devices.InsertedID)

	if result.RowsAffected == 0 {
		log.Fatal("No note with that Id exists\n")
	} else if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Printf("The statement has affected %d rows\n", result.RowsAffected)
}

func writeDebugTimingSendInfo(timingStr string) error {
	fileName := fmt.Sprintf("./logs/debug_timing_sendInfo_%d.txt", idEDGE)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(timingStr)
	if err != nil {
		return err
	}

	return nil
}

func sendInfo(devices []models.Devices) {
	// index and value
	for i, v := range devices {
		//var device = models.Device(v)

		fmt.Println(i)
		fmt.Println(v.InsertedID)
		fmt.Println(v.IdEdge)
		fmt.Println(v.IdNodo)
		fmt.Println(v.Input1)
		fmt.Println(v.Input2)
		fmt.Println(v.Output)
		fmt.Println(v.Alarm_battery)
		fmt.Println(v.Alarm_power)
		fmt.Println(v.Temperatura)
		fmt.Println(v.Umidade)
		fmt.Println(v.CreatedAt)
		fmt.Println(v.UpdatedAt)

		// Salve o tempo de envio
		sendTime := time.Now()

		var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
		client := packge.NewLoraTransactionClient(conn)

		if client != nil {

			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Println("Client ok!!")

			//t, err := time.Parse(time.RFC3339, v.CreatedAt)
			//if err != nil {
			//		panic(err)
			//	}

			//currentlyTime := timestamppb.New(t)
			//fmt.Println(currentlyTime.AsTime())
			// Lets invoke the remote function from client on the server

			// Create a new LoraRequest instance
			loraReq := &packge.LoraRequest{
				IdEdge:       v.IdEdge,
				IdNodo:       v.IdNodo,
				Input1:       v.Input1,
				Input2:       v.Input2,
				Output:       v.Output,
				AlarmBattery: v.Alarm_battery,
				AlarmPower:   v.Alarm_power,
				SensorError:  v.Sensor_error,
				Sensors: []*packge.LoraRequest_Sensor{
					{Type: 1, Value: v.Temperatura},
					{Type: 2, Value: v.Umidade},
				},
				LastUpdated: v.CreatedAt,
			}

			// Call the MakeTransaction function with the correct LoraRequest type
			tx, err := client.MakeTransaction(context.Background(), loraReq)

			// Salve o tempo de resposta
			receiveTime := time.Now()

			if err != nil {
				fmt.Println("Error, %v", err)
			} else {
				fmt.Println("sendInfo")
				fmt.Println("Msg: ", tx.Msg)
				fmt.Println("Confirmation: ", tx.Confirmation)

				// Calcule a diferença de tempo
				timeTaken := receiveTime.Sub(sendTime)

				// Converta a estrutura em uma representação de string
				timingStr := fmt.Sprintf("Enviado em: %s, Recebido em: %s, Tempo de resposta: %s\n",
					sendTime.Format("02/01/2006 15:04:05"), receiveTime.Format("02/01/2006 15:04:05"), timeTaken.String())

				if debug {
					// Escreva a informação de tempo no arquivo de debug
					if err := writeDebugTimingSendInfo(timingStr); err != nil {
						fmt.Println("Erro ao gravar no arquivo de debug:", err)
					}
				}

				if tx.Confirmation {
					deleteDevice(v)
				}
			}
		}
	}

}

func trySendDBDeviceInformation() {
	for range time.Tick(time.Second * 2) {
		//fmt.Println("Foo")

		var devices []models.Devices = returnListDevice()

		sendInfo(devices)
	}
}

func deleteAlert(alert packge.AlertRequest) {

	result := initializers.DB.Delete(&models.AlertRequest{}, "id = ?", alert.Id)

	if result.RowsAffected == 0 {
		log.Fatal("No note with that Id exists\n")
	} else if result.Error != nil {
		log.Fatal(result.Error)
	}

	fmt.Printf("The statement has affected %d rows\n", result.RowsAffected)
}

func returnListAlert() []packge.AlertRequest {

	rows, err := initializers.DB2.Query("SELECT * FROM alert_requests")

	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()

	var listAlerts = make([]packge.AlertRequest, 0)

	for rows.Next() {
		singleAlert := packge.AlertRequest{}
		err = rows.Scan(&singleAlert.Id, &singleAlert.IdEdge,
			&singleAlert.IdNodo, &singleAlert.Hash,
			&singleAlert.Type, &singleAlert.Date,
			&singleAlert.Temperatura, &singleAlert.Umidade,
			&singleAlert.Rele, &singleAlert.Description)

		if err != nil {
			fmt.Println(err)
		}

		listAlerts = append(listAlerts, singleAlert)
	}

	return listAlerts
}

func writeDebugTimingTrySendDBAlertTimer(timingStr string) error {
	fileName := fmt.Sprintf("./logs/debug_timing_trySendDBAlertTimer_%d.txt", idEDGE)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(timingStr)
	if err != nil {
		return err
	}

	return nil
}

func trySendDBAlertTimer() {

	for range time.Tick(time.Second * 2) {
		var listAlert []packge.AlertRequest = returnListAlert()

		for _, alert := range listAlert {

			// Salve o tempo de envio
			sendTime := time.Now()

			fmt.Println("alertTimer")
			var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
			client := packge.NewAlertClient(conn)

			if client != nil {

				tx, err := client.SendAlert(
					context.Background(),
					&packge.AlertRequest{
						Id:          alert.Id,
						IdEdge:      alert.IdEdge,
						IdNodo:      alert.IdNodo,
						Hash:        alert.Hash,
						Type:        alert.Type,
						Date:        alert.Date,
						Temperatura: alert.Temperatura,
						Umidade:     alert.Umidade,
						Rele:        alert.Rele,
						Description: alert.Description,
					},
				)
				// Salve o tempo de resposta
				receiveTime := time.Now()

				if err != nil {
					fmt.Println("Error, %v", err)
				} else {
					// Calcule a diferença de tempo
					timeTaken := receiveTime.Sub(sendTime)

					if debug {
						// Converta a estrutura em uma representação de string
						timingStr := fmt.Sprintf("Enviado em: %s, Recebido em: %s, Tempo de resposta: %s\n",
							sendTime.Format("02/01/2006 15:04:05"), receiveTime.Format("02/01/2006 15:04:05"), timeTaken.String())

						// Escreva a informação de tempo no arquivo de debug
						if err := writeDebugTimingTrySendDBAlertTimer(timingStr); err != nil {
							fmt.Println("Erro ao gravar no arquivo de debug:", err)
						}
					}
					fmt.Println("----------------------------------------------------------------------------------------")
					fmt.Println("OK - Alert timer")
					fmt.Println(tx)
					if tx.Confirmation {
						deleteAlert(alert)
					}

				}
			}
		}
	}
}

func writeDebugTimingSendAlert(timingStr string) error {
	fileName := fmt.Sprintf("./logs/debug_timing_sendAlert_%d.txt", idEDGE)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(timingStr)
	if err != nil {
		return err
	}

	return nil
}

func sendAlert(idEdge int32, idNodo int32, hash string, _type int32, date string, temperatura string, umidade string, rele string, description string) {

	fmt.Println("sendAlert")
	// Salve o tempo de envio
	sendTime := time.Now()

	var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
	client := packge.NewAlertClient(conn)

	newAlert := models.AlertRequest{
		IdEdge:      idEdge,
		IdNodo:      idNodo,
		Hash:        hash,
		Type:        _type,
		Date:        date,
		Temperatura: temperatura,
		Umidade:     umidade,
		Rele:        rele,
		Description: description,
	}

	if client != nil {
		tx, err := client.SendAlert(
			context.Background(),
			&packge.AlertRequest{
				IdEdge:      idEdge,
				IdNodo:      idNodo,
				Hash:        hash,
				Type:        _type,
				Date:        date,
				Temperatura: temperatura,
				Umidade:     umidade,
				Rele:        rele,
				Description: description,
			},
		)

		if err != nil {
			fmt.Println("Error, %v", err)
			result := initializers.DB.Create(&newAlert)
			if result.Error != nil {
				fmt.Println(result.Error.Error())
			}
		} else {
			// Salve o tempo de resposta
			receiveTime := time.Now()

			// Calcule a diferença de tempo
			timeTaken := receiveTime.Sub(sendTime)

			if debug {
				// Converta a estrutura em uma representação de string
				timingStr := fmt.Sprintf("Enviado em: %s, Recebido em: %s, Tempo de resposta: %s\n",
					sendTime.Format("02/01/2006 15:04:05"), receiveTime.Format("02/01/2006 15:04:05"), timeTaken.String())

				// Escreva a informação de tempo no arquivo de debug
				if err := writeDebugTimingSendAlert(timingStr); err != nil {
					fmt.Println("Erro ao gravar no arquivo de debug:", err)
				}
			}

			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Println("OK - sendAlert")
			if tx.Confirmation == false {
				fmt.Println(tx)
				initializers.DB.Create(&newAlert)
			}
		}
	}

}

func keepAlive() {
	for range time.Tick(time.Minute) {
		fmt.Println("keepAlive")
		var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
		client := packge.NewKeepAliveClient(conn)

		if client != nil {
			//t, err := time.Parse(time.RFC3339, v.CreatedAt)
			//if err != nil {
			//		panic(err)
			//	}

			//currentlyTime := timestamppb.New(t)
			//fmt.Println(currentlyTime.AsTime())
			// Lets invoke the remote function from client on the

			fmt.Println("----------------------------------------------------------------------------------------")

			// Convert Dispositivos slice to JSON
			jsonData, err := json.Marshal(packge.GetDispositivos())
			if err != nil {
				fmt.Println("Error converting to JSON:", err)
				return
			}

			// Print JSON data
			fmt.Println("JSON data:")
			fmt.Println(string(jsonData))

			now := time.Now()
			tx, err := client.MakeKeepAlive(
				context.Background(),
				&packge.KeepAliveRequest{
					IdEdge: idEDGE,
					Type:   2,
					Date:   now.String(),
					Msg:    string(jsonData),
				},
			)

			if err != nil {
				fmt.Println("Error, %v", err)
			} else {
				fmt.Println("----------------------------------------------------------------------------------------")
				fmt.Println("OK INFO (KeepAlive)")
				fmt.Println(tx)

			}
		}
	}
}

func reqReturnAllNodos() Data {

	// HTTP endpoint
	posturl := "http://127.0.0.1:8080/"

	// JSON body
	body := []byte(`{
		"version": 1.000000,
		"type": 5
		}`)

	// Create a HTTP post request
	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	fmt.Println("res.StatusCode:", res.StatusCode)

	if res.StatusCode != http.StatusOK {
		fmt.Println("Erro no cadastro (gerar alerta): %v", res.StatusCode)
	}

	jsonData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
	}
	fmt.Printf("client: response body: %s\n", jsonData)

	var data Data
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		fmt.Println("Erro ao fazer a conversão do JSON para a struct:", err)
	}

	// Percorrendo a lista ID usando o loop for e range
	for i, id := range data.ID {
		fmt.Println(i, id)
	}

	//fmt.Println(data.ID)
	return data
}

func reqPostAddNodo(id int) {
	// HTTP endpoint
	postURL := "http://127.0.0.1:8080/"

	// Criar um mapa com a estrutura do corpo JSON
	data := map[string]interface{}{
		"version": 1.000000,
		"type":    3,
		"id":      id, // Adicionar o ID fornecido ao mapa
	}

	// Serializar o mapa para obter o corpo JSON atualizado
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erro ao serializar o corpo JSON:", err)
		return
	}

	// Fazer a solicitação HTTP POST
	resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Erro na solicitação HTTP POST:", err)
		return
	}
	defer resp.Body.Close()

	// Verificar a resposta HTTP
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro na solicitação HTTP. Status: %s\n", resp.Status)
		return
	}

	fmt.Println("Solicitação HTTP POST bem-sucedida!")
}

func reqPostRemoveNodo(id int) {
	// HTTP endpoint
	postURL := "http://127.0.0.1:8080/"

	// Criar um mapa com a estrutura do corpo JSON
	data := map[string]interface{}{
		"version": 1.000000,
		"type":    4,
		"id":      id, // Remover o ID fornecido ao mapa
	}

	// Serializar o mapa para obter o corpo JSON atualizado
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Erro ao serializar o corpo JSON:", err)
		return
	}

	// Fazer a solicitação HTTP POST
	resp, err := http.Post(postURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Erro na solicitação HTTP POST:", err)
		return
	}
	defer resp.Body.Close()

	// Verificar a resposta HTTP
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erro na solicitação HTTP. Status: %s\n", resp.Status)
		return
	}

	fmt.Println("Solicitação HTTP POST bem-sucedida!")
}

func searchIDNodoLocal(data Data, id int32) bool {
	for _, idDevice := range data.ID {
		if idDevice == int(id) {
			return true
		}
	}
	return false
}

func searchIDNodoMongoDB(devices []models.DeviceNodo, id int32) bool {
	for _, value := range devices {
		if value.IdNodo == id {
			return true
		}
	}
	return false
}

func updatesDevicesMonitoringList(deviceNodo []models.DeviceNodo) {
	var _listDeviceActive = make([]*packge.DeviceActive, 0)

	// Retorna todos dispositivosNodo associado a esse EDGE
	if len(deviceNodo) > 0 {

		fmt.Println("---------------- deviceNodo to NewDeviceActive")
		for i, device := range deviceNodo {
			fmt.Println(i, device.IdNodo)
			fmt.Println(i, device.IdEdge)
			fmt.Println(i, device.Period)
			fmt.Println(i, device.Date)
			fmt.Println(i, device.Hash)

			dispositivoAtivo := packge.NewDeviceActive(int(idEDGE), int(device.IdNodo), "Nodo"+strconv.Itoa(int(i+1)))
			dispositivoAtivo.SetActive(true)
			//fmt.Println(dispositivoAtivo)

			// Criar lista para verificar se dispositivos estão ativos
			_listDeviceActive = append(_listDeviceActive, dispositivoAtivo)
		}

	}

	// Atribuir a lista de dispositivos à variável global
	packge.SetDispositivos(_listDeviceActive)
}

func reqDevices() {
	fmt.Println("reqDevices")
	var data Data
	var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
	client := packge.NewKeepAliveClient(conn)

	if client != nil {
		now := time.Now()
		response, err := client.MakeKeepAlive(
			context.Background(),
			&packge.KeepAliveRequest{
				IdEdge: idEDGE,
				Type:   1,
				Date:   now.String(),
				Msg:    "reqDevices",
			},
		)

		if err != nil {
			fmt.Println("Error, %v", err)
		} else {

			///{
			//	"version": 1.000000,
			//	"type": 3,
			//   "id": 771760128
			// }

			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Println("OK INFO (KeepAlive)")
			fmt.Println(response)

			// Processar a resposta
			fmt.Println("ID do Dispositivo (IdEdge):", response.IdEdge)
			fmt.Println("Tipo:", response.Type)
			fmt.Println("Data:", response.Date)
			fmt.Println("Mensagem:", response.Msg)
			fmt.Println("Confirmação:", response.Confirmation)

			var deviceNodo []models.DeviceNodo
			if err := json.Unmarshal([]byte(response.Msg), &deviceNodo); err != nil {
				fmt.Println("Erro ao fazer a conversão do JSON para a struct:", err)
			}

			updatesDevicesMonitoringList(deviceNodo)

			// Retorna todos dispositivos Lora cadastrado no banco local (ID NODO)
			data = reqReturnAllNodos()

			//REMOVER DISPOSITIVOS QUE ESTÃO CADASTRADOS INCORRETAMENTE
			// Percorrendo a lista de todos Nodos cadastrados SBC

			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Printf("Removendo devicesNodo\n")
			for i, id := range data.ID {
				fmt.Println(i, id)
				// Remove dispositivos que não estão associados ao EDGE

				found := searchIDNodoMongoDB(deviceNodo, int32(id))
				if found {
					fmt.Printf("O IdNodo %d foi encontrado na lista.\n", id)
					// Não faz nada
				} else {
					fmt.Printf("O IdNodo %d não foi encontrado na lista.\n", id)
					// Remove Device
					reqPostRemoveNodo(id)
				}
			}

			// ADICIONA DISPOSITIVOS QUE ESTÃO CADASTRADOS INCORRETAMENTE
			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Printf("Adicionando devicesNodo\n")
			for i, device := range deviceNodo {
				fmt.Println(i, device.IdNodo)
				fmt.Println(i, device.IdEdge)
				fmt.Println(i, device.Period)
				fmt.Println(i, device.Date)
				fmt.Println(i, device.Hash)

				found := searchIDNodoLocal(data, device.IdNodo)
				if found {
					fmt.Printf("O IdNodo %d foi encontrado na lista.\n", device.IdNodo)
					// Não faz nada
				} else {
					fmt.Printf("O IdNodo %d não foi encontrado na lista.\n", device.IdNodo)
					// Adiciona Device
					reqPostAddNodo(int(device.IdNodo))
				}

			}

		}
	}
}

func writeDebugTimingSendAlertDeviceNotResponde(timingStr string) error {
	fileName := fmt.Sprintf("./logs/debug_timing_sendAlertDeviceNotResponde_%d.txt", idEDGE)
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(timingStr)
	if err != nil {
		return err
	}

	return nil
}

func sendAlertDeviceNotResponde(idEdge int32, idNodo int32, _type int32, date string, description string) {
	// Salve o tempo de envio
	sendTime := time.Now()

	fmt.Println("sendAlertDeviceNotResponde")
	var conn *grpc.ClientConn = initializers.ReturnClientGRPC()
	client := packge.NewAlertClient(conn)

	newAlert := models.AlertRequest{
		IdEdge: idEdge,
		IdNodo: idNodo,
		Type:   _type,
		Date:   date,
	}

	if client != nil {
		tx, err := client.SendAlert(
			context.Background(),
			&packge.AlertRequest{
				IdEdge:      idEdge,
				IdNodo:      idNodo,
				Type:        _type,
				Date:        date,
				Description: description,
			},
		)

		if err != nil {
			fmt.Println("Error, %v", err)
			result := initializers.DB.Create(&newAlert)
			if result.Error != nil {
				fmt.Println(result.Error.Error())
			}
		} else {
			// Salve o tempo de resposta
			receiveTime := time.Now()

			if debug {
				// Calcule a diferença de tempo
				timeTaken := receiveTime.Sub(sendTime)

				// Converta a estrutura em uma representação de string
				timingStr := fmt.Sprintf("Enviado em: %s, Recebido em: %s, Tempo de resposta: %s\n",
					sendTime.Format("02/01/2006 15:04:05"), receiveTime.Format("02/01/2006 15:04:05"), timeTaken.String())

				// Escreva a informação de tempo no arquivo de debug
				if err := writeDebugTimingSendAlertDeviceNotResponde(timingStr); err != nil {
					fmt.Println("Erro ao gravar no arquivo de debug:", err)
				}
			}

			fmt.Println("----------------------------------------------------------------------------------------")
			fmt.Printf("Dispositivo %d (Date: %s) está inativo.", idNodo, date)
			fmt.Println("OK - sendAlertDeviceNotResponde")
			if tx.Confirmation == false {
				fmt.Println(tx)
				initializers.DB.Create(&newAlert)
			}
		}
	}

}

func checkDeviceIsActive() {

	for range time.Tick(time.Second * 30) {

		for _, d := range packge.GetDispositivos() {
			if d.IsActive == false && d.Notification == false {
				msg := fmt.Sprintf("Dispositivo %s (ID: %d) está inativo.", d.Name, d.IDNodo)
				fmt.Println(msg) // This will print the formatted message stored in the 'test' variable.

				sendAlertDeviceNotResponde(idEDGE, int32(d.IDNodo), 5, time.Now().Local().String(), msg)
				packge.SetDeviceNotif(int(d.IDNodo), true)
			}
		}
	}
}

func main() {
	go trySendDBDeviceInformation()
	go trySendDBAlertTimer()

	//sendAlert(1234, 123, "adadad", 1, "", "", "", "", "ok")
	reqDevices()
	time.Sleep(time.Second * 5)

	// Access and display the list of global devices
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println("Lista de dispositivos globais:")
	for _, dispositivo := range packge.GetDispositivos() {
		fmt.Printf("ID: %d, IDEdge: %d, IDNodo: %d, Nome: %s\n", dispositivo.ID, dispositivo.IDEdge, dispositivo.IDNodo, dispositivo.Name)
	}
	fmt.Println("----------------------------------------------------------------------------------------")

	// Podemos modificar a lista de dispositivos globais a partir de qualquer lugar
	//packge.GetDispositivos()[0].IsActive = true
	//fmt.Println("Status do Dispositivo", packge.GetDispositivos()[0].IDNodo, packge.GetDispositivos()[0].IsActive)
	//fmt.Println("----------------------------------------------------------------------------------------")

	go checkDeviceIsActive()
	go keepAlive()

	app := fiber.New()

	//run database
	//	configs.ConnectDB()

	//routes
	routes.DeviceRoute(app) //add this

	app.Listen(":3033")

}
