package repository

import (
	"context"
	"testing"

	"github.com/buni/tx-parser/internal/app/entity"
	"github.com/stretchr/testify/suite"
)

type InMemoryTransactionRepositoryTestSuite struct {
	suite.Suite
	repo *InMemoryTransactionRepository
	ctx  context.Context
}

func (s *InMemoryTransactionRepositoryTestSuite) SetupTest() {
	s.repo = NewInMemoryTransactionRepository()
	s.ctx = context.Background()
}

func (s *InMemoryTransactionRepositoryTestSuite) TestBatchCreateAndListTransactions() {
	tokenType := entity.TokenTypeETH
	transactions := []entity.Transaction{
		{
			ID:          "tx1",
			TokenType:   tokenType,
			Hash:        "hash1",
			FromAddress: "from1",
			ToAddress:   "to1",
			Value:       "100",
			Address:     "address1",
		},
		{
			ID:          "tx2",
			TokenType:   tokenType,
			Hash:        "hash2",
			FromAddress: "from2",
			ToAddress:   "to2",
			Value:       "200",
			Address:     "address2",
		},
	}
	for _, v := range transactions {
		err := s.repo.Create(s.ctx, v)
		s.NoError(err)
	}

	got, err := s.repo.List(s.ctx, tokenType, "address1")
	s.NoError(err)
	s.Equal([]entity.Transaction{transactions[0]}, got)

	got, err = s.repo.List(s.ctx, tokenType, "address2")
	s.NoError(err)
	s.Equal([]entity.Transaction{transactions[1]}, got)
}

func (s *InMemoryTransactionRepositoryTestSuite) TestListTransactionsNotFound() {
	tokenType := entity.TokenTypeETH

	got, err := s.repo.List(s.ctx, tokenType, "nonexistent")
	s.NoError(err)
	s.Empty(got)
}

func (s *InMemoryTransactionRepositoryTestSuite) TestListTransactionsAddressNotFound() {
	tokenType := entity.TokenTypeETH

	transactions := []entity.Transaction{
		{
			ID:          "tx1",
			TokenType:   tokenType,
			Hash:        "hash1",
			FromAddress: "from1",
			ToAddress:   "to1",
			Value:       "100",
			Address:     "address1",
		},
		{
			ID:          "tx2",
			TokenType:   tokenType,
			Hash:        "hash2",
			FromAddress: "from2",
			ToAddress:   "to2",
			Value:       "200",
			Address:     "address2",
		},
	}

	for _, v := range transactions {
		err := s.repo.Create(s.ctx, v)
		s.NoError(err)
	}

	got, err := s.repo.List(s.ctx, tokenType, "nonexistent")
	s.NoError(err)
	s.Empty(got)
}

func TestInMemoryTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryTransactionRepositoryTestSuite))
}
