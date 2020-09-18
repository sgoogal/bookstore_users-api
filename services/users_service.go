package services

import (
	"github.com/sgoogal/bookstore_users-api/domain/users"
	"github.com/sgoogal/bookstore_users-api/utils/crypto_utils"
	"github.com/sgoogal/bookstore_users-api/utils/date_utils"
	"github.com/sgoogal/bookstore_users-api/utils/errors"
)

var UserService userServiceInterface = &userService{}

type userServiceInterface interface {
	CreateUser(users.User) (*users.User, *errors.RestErr)
	GetUser(int64) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(int64) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
	GetLoginUser(string, string) (*users.User, *errors.RestErr)
}

type userService struct{}

func (s *userService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.DateCreated = date_utils.GetNowDBFormat()
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := &users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *userService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current, err := UserService.GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}

	} else {
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil
}

func (s *userService) DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func (s *userService) SearchUser(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *userService) GetLoginUser(userEmail string, userPassword string) (*users.User, *errors.RestErr) {
	result := &users.User{Email: userEmail, Password: crypto_utils.GetMd5(userPassword)}
	if err := result.GetLoginUser(); err != nil {
		return nil, err
	}
	return result, nil
}
