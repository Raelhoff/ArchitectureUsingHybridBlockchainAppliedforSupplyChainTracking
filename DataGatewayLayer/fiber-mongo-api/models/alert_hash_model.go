package models

type AssetHash struct {
	ID                 string `json:"ID"`
	IdDevice           string `json:"IdDevice"`
	Type               string `json:"Type"`
	ProductionDate     string `json:"ProductionDate"`
	ProductionLocation string `json:"ProductionLocation"`
	Description        string `json:"Description"`
}
