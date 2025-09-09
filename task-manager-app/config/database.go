package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	err := godotenv.Load()
	if err!=nil{
		log.Fatal("Gagal Loading .env")
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx,cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client,err := mongo.Connect(clientOptions)
	if err!=nil{
		log.Fatal("Gagal koneksi ke Database : ",err)
	}

	err = client.Ping(ctx,nil)
	if err!=nil{
		log.Fatalf("Gagal Ping ke database : %v",err)
	}
	log.Println("Sukses konek!")
	DB = client.Database(dbName)
}