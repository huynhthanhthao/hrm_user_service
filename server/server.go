package server

import (
	"context"
	"log"
	"net"

	"github.com/google/uuid"
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

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// Sinh UUID cho user
	userID := uuid.New().String()

	// Tạo user mới
	user := db.User{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	// Lưu vào database
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &proto.CreateUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	// Tìm user theo id
	var user db.User
	if err := s.DB.Where("id = ?", req.Id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "user with id %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	// Cập nhật thông tin
	user.Name = req.Name
	user.Email = req.Email

	// Lưu thay đổi
	if err := s.DB.Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &proto.UpdateUserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	// Tìm user theo id
	var user db.User
	if err := s.DB.Where("id = ?", req.Id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Errorf(codes.NotFound, "user with id %s not found", req.Id)
		}
		return nil, status.Errorf(codes.Internal, "failed to find user: %v", err)
	}

	// Xóa user
	if err := s.DB.Delete(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &proto.DeleteUserResponse{
		Message: "User deleted successfully",
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