package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"api/api"
	"api/api/handler"
	pb "api/genproto/auth"

	_ "api/docs"

	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	usersConn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Error while creating new client of USERS service: ", err.Error())
	}
	defer usersConn.Close()

	usc := pb.NewAuthServiceClient(usersConn)
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", "localhost", 6379),
	})
	
	h := handler.NewHandler(usc, *rdb)
	r := api.NewGin(h)

	err = r.Run("localhost:8080")
	if err != nil {
		log.Fatalln("Error while running server: ", err.Error())
	}

}
