package main

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/chetanmeniyabacncy/docker_microservice5/helloworld"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath)
	// err := godotenv.Load(os.ExpandEnv("$GOPATH/src/golang-master/.env"))

	// err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":"+os.Getenv("RPCPORT"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	sh := helloworld.Server{}
	sg := grpc.NewServer()
	helloworld.RegisterGreeterServer(sg, &sh)
	log.Printf("s erver listening at %v", lis.Addr())
	if err := sg.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
