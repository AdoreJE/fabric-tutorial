package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type PingPong struct {
	contractapi.Contract
}

type Asset struct {
	Key   string
	Value string
}

func (p *PingPong) ClientInfo(ctx contractapi.TransactionContextInterface) error {
	id, _ := ctx.GetClientIdentity().GetID()
	fmt.Println("GetID : ", id)

	mspid, _ := ctx.GetClientIdentity().GetMSPID()
	fmt.Println("GetMSPID : ", mspid)

	certificate, _ := ctx.GetClientIdentity().GetX509Certificate()
	fmt.Println("GetX509Certificate", certificate)

	value, found, _ := ctx.GetClientIdentity().GetAttributeValue("test1")
	if found {
		fmt.Println("Attribute test1 value : ", value)
	} else {
		fmt.Println("Attribute test1 Not Found")
	}

	err := ctx.GetClientIdentity().AssertAttributeValue("test1", "FABRIC")
	if err != nil {
		return err
	}

	return nil
}

func (p *PingPong) Ping(ctx contractapi.TransactionContextInterface, key, value string) error {
	asset := Asset{
		Key:   key,
		Value: value,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(key, assetBytes)
}

func (p *PingPong) Pong(ctx contractapi.TransactionContextInterface, key string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(key)

	if err != nil {
		return nil, err
	}

	if assetBytes == nil {
		return nil, fmt.Errorf("not found")
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

const objectType = "Car"

func (p *PingPong) CreateCar(ctx contractapi.TransactionContextInterface, brand, model, price string) error {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey(objectType, []string{brand, model})

	return ctx.GetStub().PutState(compositeKey, []byte(price))
}
func (p *PingPong) ReadCar(ctx contractapi.TransactionContextInterface, brand, model string) error {
	compositeKey, _ := ctx.GetStub().CreateCompositeKey(objectType, []string{brand, model})
	assetBytes, _ := ctx.GetStub().GetState(compositeKey)

	fmt.Println(string(assetBytes))
	return nil
}

func (p *PingPong) ReadCarByPartial(ctx contractapi.TransactionContextInterface, brand string) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey(objectType, []string{brand})

	if err != nil {
		return err
	}

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return err
		}

		fmt.Println("namespace : " + kv.GetNamespace())
		fmt.Println("key : " + kv.GetKey())
		fmt.Println("value : " + string(kv.GetValue()))
		fmt.Println()
	}

	iterator.Close()

	return nil
}

func (p *PingPong) Dummy(ctx contractapi.TransactionContextInterface) error {
	key := "Key"

	for i := 1; i <= 10; i++ {
		asset := Asset{
			Key:   key + strconv.Itoa(i),
			Value: strconv.Itoa(i),
		}

		assetBytes, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		ctx.GetStub().PutState(key+strconv.Itoa(i), assetBytes)
	}

	return nil
}

func (p *PingPong) DummyCompositeKey(ctx contractapi.TransactionContextInterface) error {
	key := "Asset"

	for i := 1; i < 30; i++ {
		compositeKey, _ := ctx.GetStub().CreateCompositeKey(key, []string{"A", "aa", strconv.Itoa(i)})
		fmt.Println(compositeKey)
		ctx.GetStub().PutState(compositeKey, []byte(strconv.Itoa(i)))
	}

	for i := 31; i < 60; i++ {
		compositeKey, _ := ctx.GetStub().CreateCompositeKey(key, []string{"A", "bb", strconv.Itoa(i)})
		ctx.GetStub().PutState(compositeKey, []byte(strconv.Itoa(i)))
	}

	for i := 1; i < 30; i++ {
		compositeKey, _ := ctx.GetStub().CreateCompositeKey(key, []string{"B", "bb", strconv.Itoa(i)})
		ctx.GetStub().PutState(compositeKey, []byte(strconv.Itoa(i)))
	}

	return nil
}

func (p *PingPong) Query(ctx contractapi.TransactionContextInterface, startKey, endKey string) error {
	iterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	if err != nil {
		return err
	}

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return err
		}

		fmt.Println("namespace : " + kv.GetNamespace())
		fmt.Println("key : " + kv.GetKey())
		fmt.Println("value : " + string(kv.GetValue()))
		fmt.Println()
	}

	return nil
}

func (p *PingPong) QueryCompositeKey(ctx contractapi.TransactionContextInterface, key string) error {
	iterator, err := ctx.GetStub().GetStateByPartialCompositeKey("Asset", []string{key})

	if err != nil {
		return err
	}

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return err
		}

		fmt.Println("namespace : " + kv.GetNamespace())
		fmt.Println("key : " + kv.GetKey())
		fmt.Println("value : " + string(kv.GetValue()))
		fmt.Println()
	}

	return nil
}

func main() {
	pingpongChaincode, err := contractapi.NewChaincode(&PingPong{})
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println("start pingpong chaincode")
	pingpongChaincode.Start()
}
