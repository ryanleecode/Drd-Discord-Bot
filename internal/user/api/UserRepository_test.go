package api_test

import (
	"drdgvhbh/discordbot/internal/user/api"
	"drdgvhbh/discordbot/internal/user/domain"
	"drdgvhbh/discordbot/internal/user/entity"
	"drdgvhbh/discordbot/mocks"
	"testing"

	pq "github.com/lib/pq"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userRepositoryConstructor struct {
	suite.Suite
}

func (constructor *userRepositoryConstructor) TestShouldNotBeASingleton() {
	assert := assert.New(constructor.T())

	userRepository := api.CreateUserRepository(
		&mocks.Connector{},
		&mocks.UserDataTransferMapper{})
	anotherUserRepository := api.CreateUserRepository(
		&mocks.Connector{},
		&mocks.UserDataTransferMapper{})

	assert.True(userRepository != anotherUserRepository)
}

func TestUserRepositoryConstructorSuite(t *testing.T) {
	suite.Run(t, new(userRepositoryConstructor))
}

type userRepositoryProvider struct {
	suite.Suite
}

func (provider *userRepositoryProvider) TestShouldBeASingleton() {
	assert := assert.New(provider.T())

	userRepository := api.ProvideUserRepository(
		&mocks.Connector{},
		&mocks.UserDataTransferMapper{})
	anotherUserRepository := api.ProvideUserRepository(
		&mocks.Connector{},
		&mocks.UserDataTransferMapper{})

	assert.True(userRepository == anotherUserRepository)
}

func TestUserRepositoryProviderSuite(t *testing.T) {
	suite.Run(t, new(userRepositoryProvider))
}

type insertion struct {
	suite.Suite
	userDataTransferMapper *mocks.UserDataTransferMapper
	databaseConnector      *mocks.Connector
	insertedRow            *mocks.Row
	userRepository         *api.UserRepository
}

func (insertion *insertion) SetupTest() {
	insertion.userDataTransferMapper = &mocks.UserDataTransferMapper{}
	insertion.userDataTransferMapper.Mock.On(
		"CreateDTOFrom",
		mock.Anything,
	).Return(&api.User{})

	insertion.insertedRow = &mocks.Row{}
	insertion.databaseConnector = &mocks.Connector{}
	insertion.databaseConnector.Mock.On(
		"QueryRow",
		mock.Anything,
	).Return(insertion.insertedRow)

	insertion.userRepository = api.CreateUserRepository(
		insertion.databaseConnector,
		insertion.userDataTransferMapper)
}

func (
	insertion *insertion,
) TestShouldReturnDuplicateUserInsertionWhenThereIsAUniqueViolation() {
	assert := assert.New(insertion.T())
	insertionError := &pq.Error{
		Code: "23505",
	}

	insertedRow := insertion.insertedRow
	insertedRow.Mock.On("Scan").Return(insertionError)

	userRepository := insertion.userRepository

	user := &entity.User{}
	err := userRepository.InsertUser(user)
	if assert.Error(err) {
		assert.Equal(
			domain.CreateDuplicateUserInsertionError(user),
			err)
	}
}

func (insertion *insertion) TestShouldNotHaveAnErrorUponSuccessfulInsertion() {
	assert := assert.New(insertion.T())

	insertedRow := insertion.insertedRow
	insertedRow.Mock.On("Scan").Return(nil)

	userRepository := insertion.userRepository

	user := &entity.User{}
	err := userRepository.InsertUser(user)

	assert.Nil(err)
}

func (insertion *insertion) TestShouldReturnOriginalPGErrorIfDoesNotHandleIt() {
	assert := assert.New(insertion.T())

	const INTERNAL_ERROR_CODE = "XX000"
	insertionError := &pq.Error{
		Code: INTERNAL_ERROR_CODE,
	}

	insertedRow := insertion.insertedRow
	insertedRow.Mock.On("Scan").Return(insertionError)

	userRepository := insertion.userRepository

	user := &entity.User{}
	err := userRepository.InsertUser(user)
	if assert.Error(err) {
		assert.Equal(insertionError, err)
	}
}

func TestUserRepositoryInsertSuite(t *testing.T) {
	suite.Run(t, new(insertion))
}
