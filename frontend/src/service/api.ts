'use client';

const BASE_URL = 'http://localhost:8080/api';

class ApiService {
    private getAuthHeaders() {
        const token = localStorage.getItem('auth_token');
        return {
            'Content-Type': 'application/json',
            ...(token && { Authorization: `Bearer ${token}` }),
        };
    }

    private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
        const response = await fetch(`${BASE_URL}${endpoint}`, {
            ...options,
            headers: {
                ...this.getAuthHeaders(),
                ...options.headers,
            },
        });

        if (!response.ok) {
            throw new Error(`API Error: ${response.status} ${response.statusText}`);
        }

        return response.json();
    }

    // Market Data
    async getMarketOverview() {
        return this.request('/v1/market/overview');
    }

    async getTokensMarketData() {
        return this.request('/v1/market/tokens');
    }

    // Lending Operations  
    async getLendingPoolInfo() {
        return this.request('/v1/lending/pool-info');
    }

    async getLendingBalance() {
        return this.request('/v1/lending/balance');
    }

    async getLendingInfo() {
        return this.request('/v1/lending/info');
    }

    async getLendingTransactions() {
        return this.request('/v1/lending/transactions');
    }

    async deposit(amount: string) {
        return this.request('/v1/lending/deposit', {
            method: 'POST',
            body: JSON.stringify({ amount }),
        });
    }

    async withdraw(amount: string) {
        return this.request('/v1/lending/withdraw', {
            method: 'POST',
            body: JSON.stringify({ amount }),
        });
    }

    // Borrowing Operations
    async getBorrowingStats() {
        return this.request('/v1/borrowing/stats');
    }

    async getBorrowingBalance() {
        return this.request('/v1/borrowing/balance');
    }

    async getBorrowingInfo() {
        return this.request('/v1/borrowing/info');
    }

    async getBorrowingTransactions() {
        return this.request('/v1/borrowing/transactions');
    }

    async borrow(amount: string) {
        return this.request('/v1/borrowing/borrow', {
            method: 'POST',
            body: JSON.stringify({ amount }),
        });
    }

    async repay(amount: string) {
        return this.request('/v1/borrowing/repay', {
            method: 'POST',
            body: JSON.stringify({ amount }),
        });
    }

    // Collateral Operations
    async getCollateralBalance() {
        return this.request('/v1/collateral/balance');
    }

    async getCollateralInfo() {
        return this.request('/v1/collateral/info');
    }

    async depositCollateral(amount: string) {
        return this.request('/v1/collateral/deposit', {
            method: 'POST',
            body: JSON.stringify({ amount }),
        });
    }

    async withdrawCollateral(amount: string) {
        return this.request('/v1/collateral/withdraw', {
            method: 'POST',
            body: JSON.stringify({ amount }),
        });
    }

    // Liquidation
    async getLiquidatablePositions() {
        return this.request('/v1/liquidation/positions');
    }

    async getLiquidationHistory() {
        return this.request('/v1/liquidation/history');
    }

    async getLiquidationBonus() {
        return this.request('/v1/liquidation/bonus');
    }
}

export const api = new ApiService();