package models

type AssetHashBack struct {
	ID          string `json:"id"`
	IdEdge      int32  `json:"idedge" validate:"required"`
	IdNodo      int32  `json:"idnodo" validate:"required"`
	Hash        string `json:"hash"`
	Type        int32  `json:"type" validate:"required"`
	Date        string `gorm:"not null;default:'1970-01-01 00:00:01'" json:"date,omitempty"`
	Temperatura string `json:"temperatura"`
	Umidade     string `json:"umidade"`
	Rele        string `json:"rele"`
	Description string `json:"description"`
}
