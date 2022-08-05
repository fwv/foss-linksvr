package bucket

import (
	"context"
	"flag"
	"linksvr/pkg/zlog"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type BucketDataSource interface {
	InsertBucket(ctx context.Context, b *Bucket) error
	FindBucket(ctx context.Context, openID string, bucketName string) (*Bucket, error)
}

var (
	MONGODB_URL    = flag.String("MONGODB_URL", "mongodb://localhost:27017/foss", "mongodb url")
	databaseName   = "foss"
	collectionName = "buckets"
)

type MongoMetaDataSource struct {
	client    *mongo.Client
	colletion *mongo.Collection
}

func NewMogoMetaDataSource(mongUri string) (*MongoMetaDataSource, error) {
	clientOpts := options.Client().ApplyURI(mongUri)
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		return nil, err
	}
	colletion := client.Database(databaseName).Collection(collectionName)
	ds := &MongoMetaDataSource{
		client:    client,
		colletion: colletion,
	}
	err = ds.createIndexes()
	if err != nil {
		zlog.Error("failed to create metadata mongodb indexes", zap.Error(err))
		return nil, err
	}
	return ds, nil
}

func (d *MongoMetaDataSource) InsertBucket(ctx context.Context, b *Bucket) error {
	_, err := d.colletion.InsertOne(ctx, b)
	if err != nil {
		return err
	}
	return nil
}

func (d *MongoMetaDataSource) FindBucket(ctx context.Context, openID string, bucketName string) (*Bucket, error) {
	// opts := options.FindOne().SetSort(bson.D{{Key: "uploadtime", Value: -1}})
	var bucket *Bucket = &Bucket{}
	if err := d.colletion.FindOne(ctx, bson.D{{Key: "openid", Value: openID}, {Key: "bucketname", Value: bucketName}}).Decode(bucket); err != nil {
		return nil, err
	}
	return bucket, nil
}

func (d *MongoMetaDataSource) createIndexes() error {
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "openid", Value: 1}, {Key: "bucketname", Value: 1}},
	}
	_, err := d.colletion.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "openid", Value: 1}},
	}
	_, err = d.colletion.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return err
	}
	return nil
}

func GetMongoBucketDataSource() BucketDataSource {
	var once sync.Once
	var ds BucketDataSource
	var err error
	once.Do(func() {
		if ds == nil {
			ds, err = NewMogoMetaDataSource(*MONGODB_URL)
			if err != nil {
				zlog.Error("failed to init mongo metadata datasource", zap.Error(err))
			}
		}
	})
	return ds
}
