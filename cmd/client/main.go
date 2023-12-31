package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/ak-karimzai/go-grpc/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	laptop.Id = "389221f6-82f1-4af9-835d-c52899151e02"
	req := pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.CreateLaptop(ctx, &req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("laptop already exist")
		} else {
			log.Fatal("cannot create laptop: ", err)
		}
		return
	}
	log.Printf("created laptop with id: %s", res.Id)
}
