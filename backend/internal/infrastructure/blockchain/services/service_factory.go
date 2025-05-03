package services

import (
	"sync"
)

// ServiceFactory provides a centralized way to access all contract services
type ServiceFactory struct {
	tokenService      *TokenService
	lendingService    *LendingPoolService
	borrowingService  *BorrowingService
	collateralService *CollateralService

	tokenOnce      sync.Once
	lendingOnce    sync.Once
	borrowingOnce  sync.Once
	collateralOnce sync.Once
}

var (
	// Single instance of ServiceFactory (singleton)
	factoryInstance *ServiceFactory
	factoryOnce     sync.Once
)

// GetInstance returns the singleton instance of ServiceFactory
func GetInstance() *ServiceFactory {
	factoryOnce.Do(func() {
		factoryInstance = &ServiceFactory{}
	})
	return factoryInstance
}

// GetTokenService returns a singleton instance of TokenService
func (f *ServiceFactory) GetTokenService() (*TokenService, error) {
	var err error

	f.tokenOnce.Do(func() {
		f.tokenService, err = NewTokenService()
	})

	if err != nil {
		return nil, err
	}

	return f.tokenService, nil
}

// GetLendingPoolService returns a singleton instance of LendingPoolService
func (f *ServiceFactory) GetLendingPoolService() (*LendingPoolService, error) {
	var err error

	f.lendingOnce.Do(func() {
		f.lendingService, err = NewLendingPoolService()
	})

	if err != nil {
		return nil, err
	}

	return f.lendingService, nil
}

// GetBorrowingService returns a singleton instance of BorrowingService
func (f *ServiceFactory) GetBorrowingService() (*BorrowingService, error) {
	var err error

	f.borrowingOnce.Do(func() {
		f.borrowingService, err = NewBorrowingService()
	})

	if err != nil {
		return nil, err
	}

	return f.borrowingService, nil
}

// GetCollateralService returns a singleton instance of CollateralService
func (f *ServiceFactory) GetCollateralService() (*CollateralService, error) {
	var err error

	f.collateralOnce.Do(func() {
		f.collateralService, err = NewCollateralService()
	})

	if err != nil {
		return nil, err
	}

	return f.collateralService, nil
}
