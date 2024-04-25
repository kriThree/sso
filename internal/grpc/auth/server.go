package auth

import (
	"context"
	"errors"
	"sso/internal/storage"

	ssov1 "github.com/kriThree/sso_contract/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appId int64) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userId int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

const emptyValue = 0

type ServerApi struct {
	ssov1.UnimplementedAuthServer
	auth Auth
	// Login
	// Register
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &ServerApi{auth: auth})
}

func (s *ServerApi) Login(ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int64(req.GetAppId()))

	if err != nil {

		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerApi) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())

	if err != nil {

		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.RegisterResponse{

		UserId: int64(userId),
	}, nil
}

func (s *ServerApi) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, int64(req.GetUserId()))

	if err != nil {

		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ssov1.IsAdminResponse{

		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is empty")
	}
	return nil
}
func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}
func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is empty")
	}
	return nil
}
