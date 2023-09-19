package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"simu/configs"
	"simu/models"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var idEDGE int32 = 1234 // Defina idEDGE
var devicesEdge = returnListDeviceEdge()
var devicesNodo = returnListDeviceNode()
var simulationCount int
var interator  int  = 1000

func generateRandomData(idEdge, idNodo int32) models.Devices {
	sensors := []models.Sensor{
		{Type: 0, Value: rand.Float32() * 100},
		{Type: 1, Value: rand.Float32() * 100},
	}

	return models.Devices{
		Version:   111.0,
		Id:        idNodo,
		Timestamp: time.Now().Unix(),
		Input1:    rand.Int31(),
		Input2:    rand.Int31(),
		Output:    rand.Int31(),
		Sensors:   sensors,
	}
}

func sendData(device models.Devices, ip string, idEdge int32) {
	url := fmt.Sprintf("http://%s/lora", ip) // Substitua "YOUR_LORA_ENDPOINT"

	payloadJSON, err := json.Marshal(device)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadJSON))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response if needed
}

func getDeviceIp(idEdge int32) string {
	for _, device := range devicesEdge {
		if device.IdEdge == idEdge {
			return device.Ip
		}
	}
	return "" // Retorna uma string vazia caso não encontre o ID
}

func simulateAndSend(devicesEdge []models.DeviceEdge, devicesNodo []models.DeviceNodo) {
	ticker := time.NewTicker(5 * time.Second) // Ticker para um minuto
	defer ticker.Stop()

	for range ticker.C {
		if simulationCount >= interator {
			fmt.Println("Simulation completed %d iterations. Stopping simulation.", interator)
			return
		}
		fmt.Printf("Simulation %d iterations", simulationCount)

		// Execute the command at the beginning of the iteration
		cmd := exec.Command("go", "run", "/home/workspace/bitburket/simulator/timeGenerator_1245.go")
		if err := cmd.Run(); err != nil {
			fmt.Println("Error executing command:", err)
			return
		}
		// Increment the simulation count
                simulationCount++
		
		var wg sync.WaitGroup
		for _, edgeDevice := range devicesEdge {
			for _, nodoDevice := range devicesNodo {
				if idEDGE == edgeDevice.IdEdge && edgeDevice.IdEdge == nodoDevice.IdEdge {
					wg.Add(1)
					go func(edgeDevice models.DeviceEdge, nodoDevice models.DeviceNodo) {
						defer wg.Done()

						randomData := generateRandomData(edgeDevice.IdEdge, nodoDevice.IdNodo)
						ip := getDeviceIp(edgeDevice.IdEdge)
						sendData(randomData, ip, edgeDevice.IdEdge)
						timestamp := time.Now().Format("2006-01-02 15:04:05")
						fmt.Printf("[%s] Sent data from data from Edge %d to Node %d\n", timestamp, edgeDevice.IdEdge, nodoDevice.IdNodo)

					}(edgeDevice, nodoDevice)
				}
			}
		}
		wg.Wait()
	}
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

func removeDevicesWithIdEdge(idEdge int32) error {
	collection := configs.GetCollection(configs.DB, "devices_nodo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.DeleteMany(ctx, bson.M{"IdEdge": idEdge})
	if err != nil {
		return err
	}

	return nil
}

func addDummyNodeDevices(idEdge int32, numDevices int) error {
	collection := configs.GetCollection(configs.DB, "devices_nodo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < numDevices; i++ {
		newNodeDevice := models.DeviceNodo{
			IdNodo: rand.Int31(), // You can generate the IdNodo as needed
			IdEdge: idEdge,       // Use the provided idEdge parameter
			Period: 60,           // Example period value
			Date:   time.Now().Format("2006-01-02 15:04:05"),
			Hash:   "dummyhash", // Example hash value
		}

		_, err := collection.InsertOne(ctx, newNodeDevice)
		if err != nil {
			return err
		}
	}

	return nil
}

func removeRandomNodeDevices(idEdge int32, numDevices int) error {
	collection := configs.GetCollection(configs.DB, "devices_nodo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"IdEdge": idEdge}
	options := options.Find().SetLimit(0) // SetLimit(0) for retrieving all matching documents
	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var matchingDevices []models.DeviceNodo
	if err := cursor.All(ctx, &matchingDevices); err != nil {
		return err
	}

	numToRemove := min(numDevices, len(matchingDevices))

	indicesToRemove := generateRandomIndices(len(matchingDevices), numToRemove)

	for _, idx := range indicesToRemove {
		deviceToRemove := matchingDevices[idx]

		_, err := collection.DeleteOne(ctx, bson.M{"_id": deviceToRemove.IdEdge})
		if err != nil {
			return err
		}
	}

	return nil
}

func generateRandomIndices(maxIndex, count int) []int {
	indices := rand.Perm(maxIndex)
	return indices[:count]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func initDB() {
	configs.ConnectDB()
}

func main() {
	fmt.Printf("Inicializando Simulador!!\n")
	initDB()

	fmt.Println("Devices Edge:")
	for _, device := range devicesEdge {
		fmt.Printf("IdEdge: %d, Ip: %s, User: %s, Date: %s, Hash: %s\n", device.IdEdge, device.Ip, device.User, device.Date, device.Hash)
	}

	fmt.Println("Devices Nodo:")
	for _, device := range devicesNodo {
		fmt.Printf("IdNodo: %d, IdEdge: %d, Period: %d, Date: %s, Hash: %s\n", device.IdNodo, device.IdEdge, device.Period, device.Date, device.Hash)
	}

	//idEdgeToAddNodes := int32(1708) // Specify the desired idEdge value
	//numNodesToAdd := 1              // Specify the number of nodes to add
	//err := addDummyNodeDevices(idEdgeToAddNodes, numNodesToAdd)
	//if err != nil {
	//	log.Fatalf("Error adding dummy node devices: %v", err)
	//}

//	addDummyNodeDevices(1234, 5)

	//	addDummyNodeDevices(2111, 1)

	//addDummyNodeDevices(7111, 1)

	//addDummyNodeDevices(9876, 1)

	// Remove devices with IdEdge == 1708
	//idEdgeToRemove := int32(1708)
	//err := removeDevicesWithIdEdge(idEdgeToRemove)
	//if err != nil {
	//	log.Fatalf("Error removing devices with IdEdge %d: %v", idEdgeToRemove, err)
	//}

	simulateAndSend(devicesEdge, devicesNodo)
}
