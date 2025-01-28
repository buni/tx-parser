package repository

import (
	"context"
	"testing"

	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/stretchr/testify/suite"
)

type InMemorySubscriptionRepositoryTestSuite struct {
	suite.Suite
	repo *InMemorySubscriptionRepository
	ctx  context.Context
}

func (s *InMemorySubscriptionRepositoryTestSuite) SetupTest() {
	s.repo = NewInMemorySubscriptionRepository()
	s.ctx = context.Background()
}

func (s *InMemorySubscriptionRepositoryTestSuite) TestCreateAndListSubscriptions() {
	tokenType := entity.TokenTypeETH
	address := "0x2527d2ed1dd0e7de193cf121f1630caefc23ac70"

	err := s.repo.Create(s.ctx, tokenType, address)
	s.NoError(err)

	addresses, err := s.repo.List(s.ctx, tokenType)
	s.NoError(err)
	s.Equal([]string{address}, addresses)
}

func (s *InMemorySubscriptionRepositoryTestSuite) TestListSubscriptionsNotFound() {
	tokenType := entity.TokenTypeETH

	addresses, err := s.repo.List(s.ctx, tokenType)
	s.Error(err)
	s.Nil(addresses)
	s.Equal(entity.ErrNotFound, err)
}

func TestInMemorySubscriptionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InMemorySubscriptionRepositoryTestSuite))
}
