package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	Store LaptopStore
}

func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

// CreateLaptop(context.Context, *CreateLaptopRequest) (*CreateLaptopResponse, error)
// CreateLaptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	laptop := req.GetLaptop()

	log.Printf("revice a create-laptop request with id: %s", laptop.Id)
	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument, "laptop id is not a valid uuid: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(
				codes.InvalidArgument, "laptop id is not a valid uuid: %v", err)
		}
		laptop.Id = id.String()
	}

	time.Sleep(time.Second * 6)
	if ctx.Err() == context.Canceled {
		log.Print("request is cancelled")
		return nil, status.Errorf(codes.Canceled, "context cancelled: %v", ctx.Err())
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Print("deadline is exceeded")
		return nil, status.Errorf(codes.DeadlineExceeded, "context timeout: %v", ctx.Err())
	}
	err := server.Store.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(
			code, "cannot save laptop to ther store: %v", err)
	}
	log.Printf("saved laptop with id: %s", laptop.Id)

	res := &pb.CreateLaptopResponse{
		Id: laptop.Id,
	}
	return res, nil
}
