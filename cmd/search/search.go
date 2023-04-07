package main

import (
	"log"
	"main/internal/composites"
	proto "main/internal/microservices/search/proto"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	postgresDBC, err := composites.NewPostgresDBComposite()
	if err != nil {
		log.Fatal("postgres composite failed", err)
	}

	searchComposite, err := composites.NewSearchComposite(postgresDBC)
	if err != nil {
		log.Fatal("search composite failed")
	}

	listen, err := net.Listen("tcp", ":"+os.Getenv("SEARCH_PORT"))
	if err != nil {
		log.Fatal("CANNOT LISTEN PORT: ", ":"+os.Getenv("SEARCH_PORT"), err.Error())
	}

	server := grpc.NewServer()

	proto.RegisterSearchServer(server, searchComposite.Service)
	log.Printf("STARTED SEARCH MICROSERVICE ON %s", ":"+os.Getenv("SEARCH_PORT"))
	err = server.Serve(listen)
	if err != nil {
		log.Println("CANNOT LISTEN PORT: ", ":"+os.Getenv("SEARCH_PORT"), err.Error())
	}
}
