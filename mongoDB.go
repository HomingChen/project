package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initialMongoDBClient() *mongo.Client {
	uri := getEnvValue("mongoDBURI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	return client
}

func closeClient(client *mongo.Client) {
	if err := client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

type Member struct {
	Id       string `bson:"_id,omitempty"`
	Name     string
	Email    string
	Password string `bson:"password,omitempty"`
	Source   string
	Photo    string `bson:"photo,omitempty"`
}

func findMemberById(id string) (Member, error) {
	client := initialMongoDBClient()
	defer closeClient(client)
	coll := client.Database("project").Collection("member")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	var result Member
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Member{}, fmt.Errorf("[ERROR] can't find data: %v", err)
		} else {
			return Member{}, fmt.Errorf("[ERROR] others: %v", err)
		}
	}
	fmt.Printf("[NOTE] member: \nid: %s\nname: %s\nemail: %s\nPWD: %s\nfrom: %s\n",
		result.Id, result.Name, result.Email, result.Password, result.Source)
	return result, nil
}

func findMemberByEmail(email string) (Member, error) {
	client := initialMongoDBClient()
	defer closeClient(client)
	coll := client.Database("project").Collection("member")
	filter := bson.M{"email": email}

	var result Member
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Member{}, fmt.Errorf("[ERROR] can't find data: %v", err)
		} else {
			return Member{}, fmt.Errorf("[ERROR] others: %v", err)
		}
	}

	fmt.Printf("[NOTE] member: \nid: %s\nname: %s\nemail: %s\nPWD: %s\nfrom: %s\n",
		result.Id, result.Name, result.Email, result.Password, result.Source)
	return result, nil
}

type insertMemberResult struct {
	Success bool
	Id      string
}

func insertMember(memberData Member) (insertMemberResult, error) {
	client := initialMongoDBClient()
	defer closeClient(client)
	coll := client.Database("project").Collection("member")
	doc := memberData
	result, err := coll.InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("[ERROR] failed to insert 1 document: %v", err)
		return insertMemberResult{}, fmt.Errorf("failed to insert 1 document: %v", err)
	} else {
		fmt.Printf("[NOTE] inserted 1 document with _ID: %v\n", result.InsertedID)
		return insertMemberResult{Success: true, Id: result.InsertedID.(string)}, nil
	}
}

func isIdDuplicated(id string) bool {
	client := initialMongoDBClient()
	defer closeClient(client)
	coll := client.Database("project").Collection("member")
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}

	var result Member
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("[ERROR] can't find data: %v", err)
			return false
		} else {
			fmt.Printf("[ERROR] others: %v", err)
			return true
		}
	} else {
		return true
	}
}

func isEmailDuplicated(email string) bool {
	client := initialMongoDBClient()
	defer closeClient(client)
	coll := client.Database("project").Collection("member")
	filter := bson.M{"email": email}

	var result Member
	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("[ERROR] can't find data: %v", err)
			return false
		} else {
			fmt.Printf("[ERROR] others: %v", err)
			return true
		}
	} else {
		return true
	}
}
