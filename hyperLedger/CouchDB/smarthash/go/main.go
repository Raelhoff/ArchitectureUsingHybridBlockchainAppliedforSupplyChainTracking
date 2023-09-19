/*
SPDX-License-Identifier: Apache-2.0
*/

package main

//package chaincode

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
	"time"
)

// SmartContract provides functions for managing an Asset
type SmartHash struct {
	contractapi.Contract
}

type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// Asset describes
// ID -> reference/serial number of the specific car part (ex: "120.47021-XXXXXXXXXXX")
// Car -> models of car using this specific car part (ex: "Fiat 500, Fiat Panda, Fiat Punto")
// Description -> detailt description of the car part (ex: "Symmetric vane; split-core castings; Black E-Coat anti-corrosive coating protects; Double disc ground friction surface")
// Brand -> Brand of the car part (ex: "Centric")
// ProductionDate ->  (ex: "DD/MM/YYYY")
// ProductionLocation -> (ex: "Saint Jose, US")
type Asset struct {
	ID                 string `json:"ID"`
	IdDevice           string `json:"IdDevice"`
	Type               string `json:"Type"`
	ProductionDate     string `json:"ProductionDate"`
	ProductionLocation string `json:"ProductionLocation"`
	Description        string `json:"Description"`
}

// adding a base set of assets to the ledger
func (s *SmartHash) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{ID: "120.47021-15486957423", IdDevice: "1234", Type: "Volkswagen", ProductionDate: "04/11/2004", ProductionLocation: "Stuttgart, Germany", Description: "Chassis"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to init assets. %v", err)
		}
	}

	return nil
}

// CreateAsset -> create and adds new asset to the network.
func (s *SmartHash) CreateAsset(ctx contractapi.TransactionContextInterface, id string, idDevice string, tipo string, productiondate string, productionlocation string, description string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:                 id,
		IdDevice:           idDevice,
		Type:               tipo,
		ProductionDate:     productiondate,
		ProductionLocation: productionlocation,
		Description:        description,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset -> returns specific asset stored in the network
func (s *SmartHash) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset -> updates existing asset in the network.
func (s *SmartHash) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, idDevice string, tipo string, productiondate string, productionlocation string, description string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:                 id,
		IdDevice:           idDevice,
		Type:               tipo,
		ProductionDate:     productiondate,
		ProductionLocation: productionlocation,
		Description:        description,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset -> deletes specific asset from the network.
func (s *SmartHash) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists -> returns true when asset with given ID exists in world state
func (s *SmartHash) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAsset -> returns asset in the network
func (s *SmartHash) GetAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", id, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAllAssets -> returns all assets in the network
func (s *SmartHash) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) ([]*Asset, error) {
	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResult.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

func (t *SmartHash) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	return getQueryResultForQueryString(ctx, queryString)
}

// getQueryResultForQueryString executes the passed in query string.
// The result set is built and returned as a byte array containing the JSON results.
func getQueryResultForQueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	return constructQueryResponseFromIterator(resultsIterator)
}

// GetAssetHistory returns the chain of custody for an asset since issuance.
func (t *SmartHash) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: ID %v", assetID)

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				ID: assetID,
			}
		}

		timestamp, err := ptypes.Timestamp(response.Timestamp)
		if err != nil {
			return nil, err
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartHash))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
