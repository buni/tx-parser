package repository

import (
	"context"
	"testing"

	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/stretchr/testify/suite"
)

type InMemoryBlockRepositoryTestSuite struct {
	suite.Suite
	repo *InMemoryBlockRepository
	ctx  context.Context
}

func (s *InMemoryBlockRepositoryTestSuite) SetupTest() {
	s.repo = NewInMemoryBlockRepository()
	s.ctx = context.Background()
}

func (s *InMemoryBlockRepositoryTestSuite) TestSetAndGetHeight() {
	tokenType := entity.TokenTypeETH
	height := "123456"

	err := s.repo.SetHeight(s.ctx, tokenType, height)
	s.NoError(err)

	retrievedHeight, err := s.repo.GetHeight(s.ctx, tokenType)
	s.NoError(err)
	s.Equal(height, retrievedHeight)
}

func (s *InMemoryBlockRepositoryTestSuite) TestGetHeightNotSet() {
	tokenType := entity.TokenTypeETH

	_, err := s.repo.GetHeight(s.ctx, tokenType)
	s.Error(err)
	s.Equal(entity.ErrBlockHightNotSet, err)
}

func TestInMemoryBlockRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryBlockRepositoryTestSuite))
}
