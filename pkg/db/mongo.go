package db

import (
	"context"
	"github.com/peizhong/letsgo/pkg/config"
	"github.com/peizhong/letsgo/pkg/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
	"unsafe"
)

type MongoHandler struct {
	connStr string
}


func (m *MongoHandler) Do(f func(client *mongo.Client) (interface{}, error)) (interface{}, error) {
	_, m.connStr = GetDBConnString("mongo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.connStr))
	if err != nil {
		return nil, err
	}
	return f(client)
}

func (m *MongoHandler) Ping() error {
	_, err := m.Do(func(client *mongo.Client) (interface{}, error) {
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		err := client.Ping(ctx, readpref.Primary())
		return nil, err
	})
	return err
}

func (m *MongoHandler) Create(i interface{}) error {
	_, err := m.Do(func(client *mongo.Client) (interface{}, error) {
		t, m := data.GetMapAsJson(i)
		collection := client.Database(config.DBName).Collection(t)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		b := (*bson.M)(unsafe.Pointer(&m))
		res, err := collection.InsertOne(ctx, *b)
		return res.InsertedID, err
	})
	return err
}