package task

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/buni/tx-parser/internal/app/contract"
	"github.com/buni/tx-parser/internal/pkg/ethclient"
	"github.com/buni/tx-parser/internal/pkg/scheduler"
)

var _ scheduler.Task = (*EthBlockParser)(nil)

type EthBlockParser struct {
	svc contract.ParserService
}

func NewEthBlockParser(svc contract.ParserService) *EthBlockParser {
	return &EthBlockParser{
		svc: svc,
	}
}

func (p *EthBlockParser) Name() string {
	return "EthBlockParserTask"
}

func (p *EthBlockParser) Interval() time.Duration {
	return time.Second * 1
}

func (p *EthBlockParser) Handle(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := p.svc.ParseNextBlock(ctx); err != nil {
				clientError := &ethclient.ErrorResponse{}
				if errors.As(err, &clientError) {
					if clientError.StatusCode == http.StatusNotFound { // this means we are fully synced, and we can finish the task and wait until the next interval
						return nil
					}
				}
				return fmt.Errorf("process next block: %w", err)
			}
		}
	}
}
