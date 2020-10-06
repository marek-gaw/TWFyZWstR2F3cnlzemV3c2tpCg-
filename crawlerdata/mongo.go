package crawlerdata

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DefaultDatabase = "crawlerDataStore"
const CollectionName = "crawler"

type MongoHandler struct {
	client   *mongo.Client
	database string
}

func NewHandler(address string) *MongoHandler {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cl, _ := mongo.Connect(ctx, options.Client().ApplyURI(address))
	mh := &MongoHandler{
		client:   cl,
		database: DefaultDatabase,
	}
	return mh
}

func (mh *MongoHandler) GetOne(d *UrlToFetch, filter interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection(CollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := collection.FindOne(ctx, filter).Decode(d)
	return err
}

func (mh *MongoHandler) GetOneMax(d *UrlToFetch, filter interface{}, sort interface{}) error {
	//Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection(CollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	queryOptions := options.FindOneOptions{}
	queryOptions.SetSort(sort)
	err := collection.FindOne(ctx, filter, &queryOptions).Decode(d)
	return err
}

func (mh *MongoHandler) GetAll(filter interface{}) []*UrlToFetch {
	collection := mh.client.Database(mh.database).Collection(CollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	cur, err := collection.Find(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*UrlToFetch
	for cur.Next(ctx) {
		data := &UrlToFetch{}
		er := cur.Decode(data)
		if er != nil {
			log.Fatal(er)
		}
		result = append(result, data)
	}
	return result
}

func (mh *MongoHandler) AddOne(d *UrlToFetch) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection(CollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, d)
	return result, err
}

func (mh *MongoHandler) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mh.client.Database(mh.database).Collection(CollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.UpdateOne(ctx, filter, update)

	return result, err
}

func (mh *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection(CollectionName)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	result, err := collection.DeleteOne(ctx, filter)
	return result, err
}
