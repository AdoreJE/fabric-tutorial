package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type ShopContract struct {
	contractapi.Contract
}

type Member struct {
	DocType string `json:"docType"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
}

type Order struct {
	DocType  string `json:"docType"`
	ID       string `json:"id"`
	MemberId string `json:"memberId"`
	Member   Member `json:"member"`
	Item     string `json:"item"`
}

func (p *ShopContract) Join(ctx contractapi.TransactionContextInterface, id, name, agestr, gender string) error {
	age, _ := strconv.Atoi(agestr)
	member := Member{
		DocType: "member",
		ID:      id,
		Name:    name,
		Age:     age,
		Gender:  gender,
	}

	memberBytes, _ := json.Marshal(member)
	return ctx.GetStub().PutState(id, memberBytes)
}

func (p *ShopContract) GetMember(ctx contractapi.TransactionContextInterface, id string) (*Member, error) {
	memberBytes, _ := ctx.GetStub().GetState(id)
	if memberBytes == nil {
		return nil, fmt.Errorf("does not exist")
	}

	var member Member
	json.Unmarshal(memberBytes, &member)

	return &member, nil
}

func (p *ShopContract) MakeOrder(ctx contractapi.TransactionContextInterface, id, memberId, item string) error {

	memberBytes, _ := ctx.GetStub().GetState(memberId)
	if memberBytes == nil {
		return fmt.Errorf("member does not exist")
	}

	var member Member
	json.Unmarshal(memberBytes, &member)

	order := Order{
		DocType:  "order",
		ID:       id,
		MemberId: memberId,
		Member:   member,
		Item:     item,
	}

	orderBytes, _ := json.Marshal(order)
	return ctx.GetStub().PutState(id, orderBytes)
}

func (p *ShopContract) GetOrder(ctx contractapi.TransactionContextInterface, id string) (*Order, error) {
	orderBytes, _ := ctx.GetStub().GetState(id)
	if orderBytes == nil {
		return nil, fmt.Errorf("does not exist")
	}

	var order Order
	json.Unmarshal(orderBytes, &order)

	return &order, nil
}

func (p *ShopContract) QueryString(ctx contractapi.TransactionContextInterface, queryString string) ([]*Member, error) {
	iterator, _ := ctx.GetStub().GetQueryResult(queryString)

	var memberList []*Member
	for iterator.HasNext() {
		kv, _ := iterator.Next()

		var member Member
		json.Unmarshal(kv.Value, &member)
		memberList = append(memberList, &member)
	}

	return memberList, nil
}

func (p *ShopContract) QueryStringWithItem(ctx contractapi.TransactionContextInterface, queryString string) ([]*Member, error) {
	iterator, _ := ctx.GetStub().GetQueryResult(queryString)

	var memberList []*Member
	for iterator.HasNext() {
		kv, _ := iterator.Next()

		var order Order
		json.Unmarshal(kv.Value, &order)

		memberList = append(memberList, &order.Member)
	}

	return memberList, nil
}

func main() {
	shopChaincode, err := contractapi.NewChaincode(&ShopContract{})
	if err != nil {
		fmt.Println("err")
	}
	fmt.Println("start shop chaincode")
	shopChaincode.Start()
}
