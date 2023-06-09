package main

import (
	"log"
	"main/internal/composites"
	proto "main/internal/microservices/profile/proto"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	postgresDBC, err := composites.NewPostgresDBComposite()
	if err != nil {
		log.Fatal("postgres db composite failed")
	}

	minioComposite, err := composites.NewMinioComposite()
	if err != nil {
		log.Fatal("minio composite failed", err.Error())
	}

	redisComposite, err := composites.NewRedisComposite()
	if err != nil {
		log.Fatal("redis composite failed")
	}

	profileComposite, err := composites.NewProfileComposite(postgresDBC, minioComposite, redisComposite)
	if err != nil {
		log.Fatal("profile composite failed")
	}

	listen, err := net.Listen("tcp", ":"+os.Getenv("PROFILE_PORT"))
	if err != nil {
		log.Fatal("CANNOT LISTEN PORT: ", ":"+os.Getenv("PROFILE_PORT"), err.Error())
	}

	server := grpc.NewServer()

	proto.RegisterProfileServer(server, profileComposite.Service)
	log.Printf("STARTED PROFILE MICROSERVICE ON %s", ":"+os.Getenv("PROFILE_PORT"))
	err = server.Serve(listen)
	if err != nil {
		log.Println("CANNOT LISTEN PORT: ", ":"+os.Getenv("PROFILE_PORT"), err.Error())
	}
}
