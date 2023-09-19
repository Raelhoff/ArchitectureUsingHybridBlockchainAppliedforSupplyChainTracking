package main

import (
//	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"os"
)

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

var contractSmartAlert *gateway.Contract = ReadContractSmartAlert()

type Asset struct {
	ID          string `json:"ID"`
	IdEdge      string `json:"IdEdge"`
	IdNodo      string `json:"IdNodo"`
	Type        string `json:"Type"`
	Hash        string `json:"Hash"`
	Date        string `json:"Date"`
	Temperatura string `json:"Temperatura"`
	Umidade     string `json:"Umidade"`
	Rele        string `json:"Rele"`
	Description string `json:"Description"`
}

func main() {
	/*
		resultInvoke, err := contractSmartAlert.SubmitTransaction("CreateAsset", device.Hash, String(device.IdEdge), String(device.IdNodo), String(device.Type), device.Hash,
					device.Date, device.Temperatura, device.Umidade, device.Rele, device.Description)
				if err != nil {
					fmt.Printf("Failed to commit transaction: %v", err)
					deleteDeviceAlertHash(newDevice)
				} else {
					fmt.Println(resultInvoke)
					fmt.Println("Delete ", device.Hash)
					deleteDeviceAlertHash(newDevice)
				}
	*/

	fmt.Println("\n --- Get0 ---")
	// ./minifab query -n smartalert -p   '"QueryAssets", "{\"selector\":{\"Hash\":\"b2d56713e73fdc0e1cd8874e0ac4fc5eae096fa19c6b5007c46be64c103c6801c72565a7b60ef48a60e568d952239af66e599be8abc8e19592073a74d07a7955\"}}"'

/*
	result, err := contractSmartAlert.SubmitTransaction("QueryAssets", "selector", "\"{\"Hash\":\"b2d56713e73fdc0e1cd8874e0ac4fc5eae096fa19c6b5007c46be64c103c6801c72565a7b60ef48a60e568d952239af66e599be8abc8e19592073a74d07a7955\"}\"")
	if err != nil {
		fmt.Printf("Failed to commit transaction: %v", err)
	} else {
		fmt.Println("Commit is successful - GET0")
//		var smart Asset
//		json.Unmarshal([]byte(string(result)), &smart)
		fmt.Println(result)
		//fmt.Println(smart.Name)
		//fmt.Println(smart.Type)
		//fmt.Println(smart.Idade)
		///fmt.Println(smart.Endereco)
		//fmt.Println(smart.password)
	}
*/


	result, err := contractSmartAlert.SubmitTransaction("AssetExists", "Assert1asdda");
        if err != nil {
                fmt.Printf("Failed to commit transaction: %v", err)
        } else {
                fmt.Println("Commit is successful - GET0")
//              var smart Asset
//              json.Unmarshal([]byte(string(result)), &smart)
                //fmt.Println(result)
		
		if string(result) == "true" {
			fmt.Println("O resultado  verdadeiro!")
		} else {
			fmt.Println("O resultado false.")
		}
              
        }


}
