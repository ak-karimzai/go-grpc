package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxImageSize = 1 << 20

type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	laptopStore LaptopStore
	imageStore  ImageStore
}

func NewLaptopServer(laptopStore LaptopStore, imageStore ImageStore) *LaptopServer {
	return &LaptopServer{
		laptopStore: laptopStore,
		imageStore:  imageStore,
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

	// time.Sleep(time.Second * 6)
	if err := contextError(ctx); err != nil {
		return nil, err
	}
	err := server.laptopStore.Save(laptop)
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

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logErr(status.Errorf(codes.Canceled, "context cancelled: %v", ctx.Err()))
	case context.DeadlineExceeded:
		return logErr(status.Errorf(codes.DeadlineExceeded, "context timeout: %v", ctx.Err()))
	default:
		return nil
	}
}

func (server *LaptopServer) SearchLaptop(
	req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)

	err := server.laptopStore.Search(
		stream.Context(),
		filter,
		func(laptop *pb.Laptop) error {
			res := &pb.SearchLaptopResponse{
				Laptop: laptop,
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}

			log.Printf("sent laptop with id: %s", laptop.GetId())
			return nil
		},
	)
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (server *LaptopServer) UploadImage(
	stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		logErr(err)
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()

	log.Printf(
		"receive an upload-image request for laptop %s with image type %s",
		laptopId,
		imageType)

	laptop, err := server.laptopStore.Fetch(laptopId)
	if err != nil {
		return logErr(
			status.Errorf(codes.Internal, "cannot find laptop: %v", err))
	}

	if laptop == nil {
		return logErr(
			status.Errorf(codes.InvalidArgument, "laptop %s doesn't exist", laptop))
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		if err := contextError(stream.Context()); err != nil {
			return err
		}
		log.Print("waiting to reveive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}

		if err != nil {
			return logErr(status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err))
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		log.Printf("recevied a chunk with size: %d", size)

		imageSize += size

		if imageSize > maxImageSize {
			return logErr(
				status.Errorf(codes.InvalidArgument,
					"image is too large: %d > %d", imageSize, maxImageSize))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return logErr(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
		}
	}

	imageID, err := server.imageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		return logErr(status.Errorf(codes.Internal, "cannot save imager to the store: %v", err))
	}

	res := &pb.UploadImageResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return logErr(
			status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}
	log.Printf("saved image with id: %s, size: %d", imageID, imageSize)
	return nil
}

func logErr(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
