package packge

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var contractSmartHash *gateway.Contract // Variável global para o contrato

func ReadContractSmartHash() *gateway.Contract {
	// Client instance
	var wallet *gateway.Wallet = ReadWallet()

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile("/home/workspace/bitburket/server/connection.json")),
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

// Inicializa o contrato no início do aplicativo
func InitializeContractSmartHash() {
	contractSmartHash = ReadContractSmartHash()
}

// Obtém a instância do contrato
func GetContractSmartHash() *gateway.Contract {
	return contractSmartHash
}
