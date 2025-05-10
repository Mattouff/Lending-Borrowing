package routes

import (
    "github.com/Mattouff/Lending-Borrowing/internal/domain/repository"
    "github.com/Mattouff/Lending-Borrowing/internal/domain/service"
)

// Services container for passing services to route handlers
type Services struct {
	UserService        service.UserService
	LendingService     service.LendingService
	BorrowingService   service.BorrowingService
	CollateralService  service.CollateralService
	LiquidationService service.LiquidationService
}

// Repositories container for passing repositories to route handlers
type Repositories struct {
    UserRepository     repository.UserRepository
    PositionRepository repository.PositionRepository
    TransactionRepository repository.TransactionRepository
}