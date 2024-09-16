package auth

import (
	"context"
	"errors"
	ssov1 "github.com/aleksvin8888/sso-protos/gen/go/sso"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sso/internal/services/auth"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth     Auth
	validate *validator.Validate
}

func Register(gRPC *grpc.Server, auth Auth, validate *validator.Validate) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{
		auth:     auth,
		validate: validate,
	})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validateLogin(req, s.validate); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredential) {
			return nil, status.Error(codes.InvalidArgument, "invalid argument")
		}

		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	if err := validateRegister(req, s.validate); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExist) {
			return nil, status.Error(codes.AlreadyExists, "user already exist")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {

	if err := validateIsAdmin(req, s.validate); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest, validate *validator.Validate) error {

	if err := validate.Var(req.GetEmail(), "required"); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if err := validate.Var(req.GetEmail(), "email"); err != nil {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}

	if err := validate.Var(req.GetPassword(), "required"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if err := validate.Var(req.GetAppId(), "required"); err != nil {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest, validate *validator.Validate) error {
	if err := validate.Var(req.GetEmail(), "required"); err != nil {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if err := validate.Var(req.GetEmail(), "email"); err != nil {
		return status.Error(codes.InvalidArgument, "email is invalid")
	}

	if err := validate.Var(req.GetPassword(), "required"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if err := validate.Var(req.GetPassword(), "min=8"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required min 8 symbols")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest, validate *validator.Validate) error {

	if err := validate.Var(req.GetUserId(), "required"); err != nil {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
