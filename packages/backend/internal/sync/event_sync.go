package sync

import (
	"context"
	"time"

	"linkworld-backend/internal/models"
	"linkworld-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventSync struct {
	ethClient *ethclient.Client
	userRepo  *repository.UserRepository
	contracts map[string]common.Address
}

func NewEventSync(ethClient *ethclient.Client, userRepo *repository.UserRepository, contracts map[string]common.Address) *EventSync {
	return &EventSync{
		ethClient: ethClient,
		userRepo:  userRepo,
		contracts: contracts,
	}
}

func (s *EventSync) Start(ctx context.Context) {
	go s.syncUserRegistered(ctx)
}

func (s *EventSync) syncUserRegistered(ctx context.Context) {
	// Listen for UserRegistered events from UserRegistry contract
	// TODO: Implement actual filter query with ethclient

	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(30 * time.Second)
		}
	}
}

func (s *EventSync) processUserRegistered(log types.Log) error {
	if len(log.Topics) < 2 {
		return nil
	}

	userAddr := common.HexToAddress(log.Topics[1].Hex())
	user := &models.User{
		WalletAddr:   userAddr.Hex(),
		RegisteredAt: time.Unix(0, 0),
		IsActive:     true,
	}
	return s.userRepo.Create(user)
}

func (s *EventSync) processDepositMade(log types.Log) error {
	return nil
}