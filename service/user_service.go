package service

import (
	"Seckill/datamoudles"
	"Seckill/repositories"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	IsPwdSuccess(userName string, pwd string) (user *datamoudles.User, isOk bool)
	AddUser(user *datamoudles.User) (userID int64, err error)
}

func NewService(repository repositories.IUserRepository) IUserService {
	return &UserService{repository}
}

type UserService struct {
	UserRepository repositories.IUserRepository
}

func (u *UserService) IsPwdSuccess(userName string, pwd string) (user *datamoudles.User, isOk bool) {
	//TODO implement me
	user, err := u.UserRepository.Select(userName)
	if err != nil {
		return
	}
	isOk, _ = ValidatePassword(pwd, user.HashPassword)
	if !isOk {
		return &datamoudles.User{}, false
	}
	return

}

func ValidatePassword(userPassword string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
func (u *UserService) AddUser(user *datamoudles.User) (userID int64, err error) {
	//TODO implement me
	pwdByte, errPwd := GeneratePassword(user.HashPassword)
	if errPwd != nil {
		return userID, errPwd
	}
	user.HashPassword = string(pwdByte)
	return u.UserRepository.Insert(user)
}
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
