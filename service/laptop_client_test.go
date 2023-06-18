package service_test

import (
	"context"
	"net"
	"testing"

	"github.com/ak-karimzai/go-grpc/pb"
	"github.com/ak-karimzai/go-grpc/sample"
	"github.com/ak-karimzai/go-grpc/serializer"
	"github.com/ak-karimzai/go-grpc/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()

	laptopServer, serverAddress :=
		startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	res, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedId, res.Id)

	other, err := laptopServer.Store.Fetch(res.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	requireSameLaptop(t, other, laptop)
}

func startTestLaptopServer(
	t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.
		NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(
	t *testing.T, serverAdd string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAdd, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(
	t *testing.T, lhs, rhs *pb.Laptop) {
	json1, err := serializer.ProtobufToJSON(lhs)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(rhs)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
