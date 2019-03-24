package play

import (
	"context"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ISQL interface {
	Open(connectionString string)
	Read(tableName string) []string
	Close()
}

type INoSQL interface {
	Open(connectionString string)
	Read(tableName string) []string
	Close()
}

type classify struct {
	gorm.Model
	id   string
	name string `gorm:"column classify_name"`
}

func (classify) TableName() string {
	return "dm_classify"
}

type asset struct {
	gorm.Model
	id         string
	classifyID string
	name       string
	json       string
}

type device struct {
	asset
}

func (d device) DeviceName() string {
	return d.name
}

type part struct {
	asset
}

func (p part) PartName() string {
	return p.name
}

func FromSQL2NoSQL(sql string, nosql string, tables ...string) {
	//result := make(map[string]int)
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://192.168.3.19:27017"))
	if err == nil {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		collection := client.Database("Logging").Collection("DefaultLog")
		cur, err := collection.Find(ctx, bson.M{"Level": "Info"})
		if err != nil {
			log.Fatal(err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			var result bson.M
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			// do something with result....
		}
	}
}
