package DBManager

import (
	"context"
	"log"

	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var configErr = godotenv.Load()
var dbURL string = os.Getenv("DB_SOURCE_URL")

var SystemCollections VAICollections

type VAICollections struct {
	Inventory          *mongo.Collection
	Material           *mongo.Collection
	UnitsOfMeasurement *mongo.Collection
	Product            *mongo.Collection
	Contact            *mongo.Collection
	PriceList          *mongo.Collection
	Setting            *mongo.Collection
	ProdStages         *mongo.Collection
	Sales              *mongo.Collection
	Purchasing         *mongo.Collection
	InventoryDemand    *mongo.Collection
}

func getMongoDbConnection() (*mongo.Client, error) {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func GetMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
	client, err := getMongoDbConnection()
	if err != nil {
		return nil, err
	}
	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

func InitCollections() bool {
	if configErr != nil {
		return false
	}
	var err error
	SystemCollections.Inventory, err = GetMongoDbCollection("vai_db", "inventories")
	if err != nil {
		return false
	}

	SystemCollections.Setting, err = GetMongoDbCollection("vai_db", "setting")
	if err != nil {
		return false
	}

	SystemCollections.Contact, err = GetMongoDbCollection("vai_db", "contact")
	if err != nil {
		return false
	}

	SystemCollections.Product, err = GetMongoDbCollection("vai_db", "product")
	if err != nil {
		return false
	}

	SystemCollections.Material, err = GetMongoDbCollection("vai_db", "material")
	if err != nil {
		return false
	}

	SystemCollections.ProdStages, err = GetMongoDbCollection("vai_db", "prodstages")
	if err != nil {
		return false
	}

	SystemCollections.PriceList, err = GetMongoDbCollection("vai_db", "price_list")
	if err != nil {
		return false
	}

	SystemCollections.Sales, err = GetMongoDbCollection("vai_db", "sales")
	if err != nil {
		return false
	}

	SystemCollections.Purchasing, err = GetMongoDbCollection("vai_db", "purchasing")
	if err != nil {
		return false
	}
	SystemCollections.UnitsOfMeasurement, err = GetMongoDbCollection("vai_db", "unitsofmeasurement")
	if err != nil {
		return false
	}

	SystemCollections.InventoryDemand, err = GetMongoDbCollection("vai_db", "inventorydemand")
	if err != nil {
		return false
	}

	return true
}
