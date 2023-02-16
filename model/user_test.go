package model_test

import (
	"testing"

	"github.com/abanoub-fathy/bebo-gallery/model"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	suite.Suite
	*model.Service
}

func (s *UserServiceSuite) SetupSuite() {
	// create new user service
	const DB_URI = "postgresql://postgres:popTop123@localhost:5432/bebo-gallery_test?sslmode=disable"

	// create new service
	service, err := model.NewService(DB_URI)
	if err != nil {
		s.T().Fatal("Unable to create service", err)
	}

	// create new UserService
	s.Service = service
}

func (s *UserServiceSuite) SetupTest() {
	// Reset all the data in the DB
	s.Service.ResetDB()
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
