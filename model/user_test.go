package model_test

import (
	"testing"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserServiceSuite struct {
	suite.Suite
	model.UserService
}

func (s *UserServiceSuite) SetupSuite() {
	// create new user service
	const DB_URI = "postgresql://postgres:popTop123@localhost:5432/bebo-gallery_test?sslmode=disable"

	// open db connection to be used in all services
	db, err := gorm.Open(postgres.Open(DB_URI), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Unable to Open db connection. Error:", err)
	}

	// create new UserService
	userService := model.NewUserService(db)

	s.UserService = userService
}

func (s *UserServiceSuite) SetupTest() {
	// Drop existing tables in test database
	s.UserService.ResetUserDB()
}

func (s *UserServiceSuite) TestCreateUser() {
	user := model.User{
		FirstName: "Abanoub",
		LastName:  "Fathy",
		Email:     "aop4ever@gmail.com",
		Password:  "12212154554554asdsa",
	}
	err := s.UserService.CreateUser(&user)
	s.Require().NoError(err, "It should be no error while create user")
	s.Assert().NotEqual(user.ID.String(), "", "The Id of created user should not be empty")
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}
