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
)

var idEDGE1 int32 = 1234 // Defina idEDGE
var idEDGE2 int32 = 9876 // Defina idEDGE
var devicesEdge = returnListDeviceEdge()
var devicesNodo = returnListDeviceNode()
var simulationCount int
var interator int = 1000

var runCommand9876 bool = true // Variável global booleana
var runCommand1708 bool = true // Variável global booleana

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
	ticker := time.NewTicker(800 * time.Millisecond) // Ticker para um minuto
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

		if runCommand9876 == true {
			// Execute the command at the beginning of the iteration
	                cmd := exec.Command("go", "run", "/home/workspace/bitburket/simulator/timeGenerator_9876.go")
        	        if err := cmd.Run(); err != nil {
                	        fmt.Println("Error executing command:", err)
                        	return
                	}
		
		}

                if runCommand1708 == true {
                        // Execute the command at the beginning of the iteration
                        cmd := exec.Command("go", "run", "/home/workspace/bitburket/simulator/timeGenerator_1708.go")
                        if err := cmd.Run(); err != nil {
                                fmt.Println("Error executing command:", err)
                                return
                        }

                }


		// Increment the simulation count
		simulationCount++

		var wg sync.WaitGroup
		for _, edgeDevice := range devicesEdge {
			for _, nodoDevice := range devicesNodo {
				if edgeDevice.IdEdge == nodoDevice.IdEdge {
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

	simulateAndSend(devicesEdge, devicesNodo)
}
