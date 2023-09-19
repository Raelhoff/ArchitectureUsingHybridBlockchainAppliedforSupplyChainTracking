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
type SmartAlertHash struct {
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
	HASH string `json:"HASH"`
}

// adding a base set of assets to the ledger
func (s *SmartAlertHash) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		{HASH: "120.47021-15486957423"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.HASH, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to init assets. %v", err)
		}
	}

	return nil
}

// CreateAsset -> create and adds new asset to the network.
func (s *SmartAlertHash) CreateAsset(ctx contractapi.TransactionContextInterface, hash string) error {
	exists, err := s.AssetExists(ctx, hash)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", hash)
	}

	asset := Asset{
		HASH: hash,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hash, assetJSON)
}

// ReadAsset -> returns specific asset stored in the network
func (s *SmartAlertHash) ReadAsset(ctx contractapi.TransactionContextInterface, hash string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", hash)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset -> updates existing asset in the network.
func (s *SmartAlertHash) UpdateAsset(ctx contractapi.TransactionContextInterface, hash string) error {
	exists, err := s.AssetExists(ctx, hash)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", hash)
	}

	// overwriting original asset with new asset
	asset := Asset{
		HASH: hash,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(hash, assetJSON)
}

// DeleteAsset -> deletes specific asset from the network.
func (s *SmartAlertHash) DeleteAsset(ctx contractapi.TransactionContextInterface, hash string) error {
	exists, err := s.AssetExists(ctx, hash)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", hash)
	}

	return ctx.GetStub().DelState(hash)
}

// AssetExists -> returns true when asset with given ID exists in world state
func (s *SmartAlertHash) AssetExists(ctx contractapi.TransactionContextInterface, hash string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(hash)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAsset -> returns asset in the network
func (s *SmartAlertHash) GetAsset(ctx contractapi.TransactionContextInterface, hash string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(hash)

	if err != nil {
		return nil, fmt.Errorf("failed to get asset %s: %v", hash, err)
	}
	if assetBytes == nil {
		return nil, fmt.Errorf("asset %s does not exist", hash)
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// GetAllAssets -> returns all assets in the network
func (s *SmartAlertHash) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
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

func (t *SmartAlertHash) QueryAssets(ctx contractapi.TransactionContextInterface, queryString string) ([]*Asset, error) {
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
func (t *SmartAlertHash) GetAssetHistory(ctx contractapi.TransactionContextInterface, assetID string) ([]HistoryQueryResult, error) {
	log.Printf("GetAssetHistory: hash %v", assetID)

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
				HASH: assetID,
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

	chaincode, err := contractapi.NewChaincode(new(SmartAlertHash))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
