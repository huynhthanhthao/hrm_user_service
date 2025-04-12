package server

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"user_service/database"
	"user_service/proto"
)

type UserServer struct {
	proto.UnimplementedUserServiceServer
	DB *gorm.DB
}

func (s *UserServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	var user db.User
	
	if err := s.DB.Where("id = ?", req.Id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "user with id %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &proto.GetUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func StartServer(port string, db *gorm.DB) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, &UserServer{DB: db})
	log.Printf("gRPC server running on port %s", port)
	return grpcServer.Serve(listener)
}