package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/ak-karimzai/go-grpc/service"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("port", 0, "the server port")
	flag.Parse()

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")

	laptopServer := service.NewLaptopServer(
		laptopStore, imageStore)
	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	// pb.RegisterLaptopServiceServer(grpcServer, )
	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
	fmt.Println("server")
}
