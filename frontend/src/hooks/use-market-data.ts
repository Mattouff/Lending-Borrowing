'use client';

import { useSimpleKit } from '@/components/simplekit/simplekit';
import { api } from '@/service/api';
import { useQuery } from '@tanstack/react-query';

export function useMarketData() {
    const { isAuthenticated } = useSimpleKit();

    // Market Overview (public)
    const { data: marketOverview, isLoading: marketLoading } = useQuery({
        queryKey: ['market-overview'],
        queryFn: api.getMarketOverview,
        refetchInterval: 30000, // Refresh every 30s
    });

    // Tokens Market Data (public)
    const { data: tokensData, isLoading: tokensLoading } = useQuery({
        queryKey: ['tokens-market'],
        queryFn: api.getTokensMarketData,
        refetchInterval: 30000,
    });

    // Lending Pool Info (public)
    const { data: lendingPool, isLoading: lendingPoolLoading } = useQuery({
        queryKey: ['lending-pool-info'],
        queryFn: api.getLendingPoolInfo,
        refetchInterval: 30000,
    });

    // Borrowing Stats (public)
    const { data: borrowingStats, isLoading: borrowingStatsLoading } = useQuery({
        queryKey: ['borrowing-stats'],
        queryFn: api.getBorrowingStats,
        refetchInterval: 30000,
    });

    // User Data (authenticated only)
    const { data: lendingBalance } = useQuery({
        queryKey: ['lending-balance'],
        queryFn: api.getLendingBalance,
        enabled: isAuthenticated,
        refetchInterval: 15000,
    });

    const { data: borrowingBalance } = useQuery({
        queryKey: ['borrowing-balance'],
        queryFn: api.getBorrowingBalance,
        enabled: isAuthenticated,
        refetchInterval: 15000,
    });

    const { data: collateralBalance } = useQuery({
        queryKey: ['collateral-balance'],
        queryFn: api.getCollateralBalance,
        enabled: isAuthenticated,
        refetchInterval: 15000,
    });

    const isLoadingAPI = marketLoading || tokensLoading || lendingPoolLoading || borrowingStatsLoading;

    return {
        // Market data
        marketOverview,
        tokensData,
        lendingPool,
        borrowingStats,

        // User data
        lendingBalance,
        borrowingBalance,
        collateralBalance,

        // Loading states
        isLoadingAPI,
        isAuthenticated,
    };
}