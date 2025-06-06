package service

import (
	"context"
	"kukuh/go-gin-library-project/app/model"
	"kukuh/go-gin-library-project/app/repository"
	"kukuh/go-gin-library-project/app/web"
	"kukuh/go-gin-library-project/helper"
	"kukuh/go-gin-library-project/helper/token"
	"kukuh/go-gin-library-project/response"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UserService interface {
	Register(ctx context.Context, request *web.Register) (*web.UserResponse, *response.CustomError)
	Login(ctx context.Context, request *web.LoginUserRequest) (*web.LoginUserResponse, *response.CustomError)
	UpdateUserOwn(ctx context.Context, request *web.UpdateUserRequest) (*web.UserResponse, *response.CustomError)
	Delete(ctx context.Context, userId int) *response.CustomError
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	DB             *gorm.DB
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, DB *gorm.DB, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
		Validate:       validate,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, request *web.Register) (*web.UserResponse, *response.CustomError) {
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	password, err := helper.HashPassword(request.Password)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}
	user := model.User{
		Name:      request.Name,
		Email:     request.Email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.UserRepository.Save(s.DB, &user)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	userResponse := web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	return &userResponse, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request *web.LoginUserRequest) (*web.LoginUserResponse, *response.CustomError) {
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	var user model.User

	err = s.UserRepository.FindByEmail(s.DB, &user, request.Email)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	err = helper.CheckPasswordHash(user.Password, request.Password)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	token, err := token.GenerateJwtToken(strconv.Itoa(user.Id))
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}
	loginResponse := web.LoginUserResponse{
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	return &loginResponse, nil
}

func (s *UserServiceImpl) UpdateUserOwn(ctx context.Context, request *web.UpdateUserRequest) (*web.UserResponse, *response.CustomError) {
	err := s.Validate.Struct(request)
	if err != nil {
		return nil, response.BadRequestError(err.Error())
	}

	var user model.User
	err = s.UserRepository.FindById(s.DB, &user, request.Id)
	if err != nil {
		return nil, response.NotFoundError(err.Error())
	}

	password, err := helper.HashPassword(request.Password)
	if err != nil {
		return nil, response.GeneralError(err.Error())
	}

	user.Name = request.Name
	user.Password = password

	err = s.UserRepository.Update(s.DB, &user)
	if err != nil {
		return nil, response.RepositoryError(err.Error())
	}

	userResponse := web.UserResponse{
		Email: user.Email,
		Name:  user.Name,
	}

	return &userResponse, nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, userId int) *response.CustomError {
	var user model.User
	err := s.UserRepository.FindById(s.DB, &user, userId)
	if err != nil {
		return response.RepositoryError(err.Error())
	}

	s.UserRepository.Delete(s.DB, userId)

	return nil
}
