package packge

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"

	"fmt"
	"os"
)

var contractSmartAlert *gateway.Contract // Variável global para o contrato

func ReadWallet() *gateway.Wallet {

	w, err := gateway.NewFileSystemWallet("/home/workspace/bitburket/server/packge/wallets")
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
		gateway.WithConfig(config.FromFile("/home/workspace/bitburket/server/packge/connection.json")),
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

// Inicializa o contrato no início do aplicativo
func InitializeContract() {
	contractSmartAlert = ReadContractSmartAlert()
}

// Obtém a instância do contrato
func GetContractSmartAlert() *gateway.Contract {
	return contractSmartAlert
}
