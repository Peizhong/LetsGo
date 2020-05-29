package db

import (
	"context"
	"errors"
	"log"
	"reflect"
	"time"

	"github.com/peizhong/letsgo/pkg/config"
	"github.com/peizhong/letsgo/pkg/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoHandler struct {
	connStr string
}

func structTobsonM(i interface{})bson.M {
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	_, m := data.GetMapAsJson(v.Interface())
	bm := mapTobsonM(m)
	return bm
}

func mapTobsonM(m map[string]interface{}) bson.M {
	// b := (*bson.M)(unsafe.Pointer(&m))
	b := bson.M{}
	for k, v := range m {
		b[k] = v
	}
	return b
}

// query, ordered params
func queryTobsonD(m []Query) bson.D {
	d := bson.D{}
	for _, v := range m {
		if v.Op=="="{
			d = append(d, bson.E{
				Key:   v.Key,
				Value: v.Value,
			})
		}
	}
	return d
}

func (m *MongoHandler) SetContext(string, interface{}) {

}

func (m *MongoHandler) do(f func(client *mongo.Client) (interface{}, error)) (interface{}, error) {
	_, m.connStr = getDBConnString("mongo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.connStr))
	if err != nil {
		return nil, err
	}
	return f(client)
}

func (m *MongoHandler) colletion(i interface{}, f func(*mongo.Collection) (int, error)) (int, error) {
	_, m.connStr = getDBConnString("mongo")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	// todo: 连接池
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.connStr))
	if err != nil {
		return 0, err
	}
	table := data.GetTypeName(i)
	// use db
	// databases hold collections of documents
	collection := client.Database(config.DBName).Collection(table)
	return f(collection)
}

func (m *MongoHandler) Ping() error {
	_, err := m.do(func(client *mongo.Client) (interface{}, error) {
		ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
		err := client.Ping(ctx, readpref.Primary())
		return nil, err
	})
	return err
}

func (m *MongoHandler) Create(i interface{}) error {
	_, err := m.colletion(i, func(collection *mongo.Collection) (int, error) {
		_, m := data.GetMapAsJson(i)
		// M is an unordered representation of a BSON document
		// Elements will be serialized in an undefined, random order
		b := mapTobsonM(m)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		res, err := collection.InsertOne(ctx, b)
		log.Println("insert", res.InsertedID)
		return 1, err
	})
	return err
}

func (m *MongoHandler) Get(i interface{}, q ...Query) error {
	_, err := m.colletion(i, func(collection *mongo.Collection) (int, error) {
		// D is an ordered representation of a BSON document
		filter := queryTobsonD(q)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		err := collection.FindOne(ctx, filter).Decode(i)
		return 0, err
	})
	return err
}

func (m *MongoHandler) Gets(i interface{}, q ...Query) (int, error) {
	cnt, err := m.colletion(i, func(collection *mongo.Collection) (int, error) {
		t := reflect.TypeOf(i)
		if t.Kind() != reflect.Ptr {
			log.Println(t.Kind())
			return 0, errors.New("not ptr")
		}
		t = t.Elem()
		if t.Kind() != reflect.Slice {
			log.Println(t.Kind())
			return 0, errors.New("ptr not to slice")
		}
		t = t.Elem()
		if t.Kind() != reflect.Ptr {
			log.Println(t.Kind())
			return 0, errors.New("slice item not ptr")
		}
		t = t.Elem()
		if t.Kind() != reflect.Struct {
			log.Println(t.Kind())
			return 0, errors.New("prt not to struct")
		}
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		filter := queryTobsonD(q)
		// return a cursor
		cur, err := collection.Find(ctx, filter)
		if err != nil {
			return 0, err
		}
		defer cur.Close(ctx)
		var cnt int
		for cur.Next(ctx) {
			n := reflect.New(t)
			// prt to stuct
			err = cur.Decode(n.Interface())
			if err != nil {
				return cnt, err
			}
			log.Println(n.Interface())
			cnt++
		}
		return cnt, nil
	})
	return cnt, err
}

func (m *MongoHandler) Update(i interface{}) error {
	_, err := m.colletion(i, func(collection *mongo.Collection) (int, error) {
		if keymap := data.GetPrimaryKey(i); keymap != nil {
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			filter := mapTobsonM(keymap)
			// https://docs.mongodb.com/manual/reference/operator/update/
			// use $set
			update := structTobsonM(i)
			collection.ReplaceOne(ctx, filter, update)
		}
		return 0, nil
	})
	return err
}
