package models

//import "time"

type QueryAssets struct {
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
