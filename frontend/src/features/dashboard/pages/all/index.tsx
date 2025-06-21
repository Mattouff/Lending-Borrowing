'use client';

import { MarketOverview } from '@/components/market/market-overview';
import { ProtocolStats } from '@/components/protocol/protocol-stats';
import { ConnectWalletButton } from '@/components/simplekit/simplekit';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ChartConfig } from '@/components/ui/chart';
import { UserDashboard } from '@/components/user/user-dashboard';
import { useFormattedContractData } from '@/hooks/use-formatted-contract-data';
import { useMarketData } from '@/hooks/use-market-data';
import { Loader2 } from 'lucide-react';

const chartData = [
  { month: 'January', desktop: 186, mobile: 80 },
  { month: 'February', desktop: 305, mobile: 200 },
  { month: 'March', desktop: 237, mobile: 120 },
  { month: 'April', desktop: 73, mobile: 190 },
  { month: 'May', desktop: 209, mobile: 130 },
  { month: 'June', desktop: 214, mobile: 140 },
];

const chartConfig = {
  desktop: {
    label: 'Desktop',
    color: 'hsl(var(--chart-1))',
  },
  mobile: {
    label: 'Mobile',
    color: 'hsl(var(--chart-2))',
  },
} satisfies ChartConfig;

export default function Global() {
  const {
    formatted,
    isLoading: contractLoading,
    isError,
    error,
  } = useFormattedContractData();

  const {
    marketOverview,
    lendingPool,
    borrowingStats,
    lendingBalance,
    borrowingBalance,
    collateralBalance,
    isAuthenticated,
    isLoadingAPI,
  } = useMarketData();

  if (isError) {
    console.error('Contract data error:', error);
  }

  return (
    <div className="space-y-6">
      {/* Cards de statistiques en haut */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Supply Token
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {contractLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                `${formatted.token.totalSupply} ${formatted.token.symbol}`
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Lending</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {contractLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                `${formatted.lendingPool.totalLending} ${formatted.token.symbol}`
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Interest Rate: {formatted.lendingPool.interestRate}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Total Borrowed
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {contractLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                `${formatted.borrowing.totalBorrowed} ${formatted.token.symbol}`
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Current Rate: {formatted.borrowing.currentRate}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Collateral Status
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {contractLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                formatted.collateral.liquidationThreshold
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Liquidation Threshold
            </p>
          </CardContent>
        </Card>
      </div>
      <div className=" p-4 space-y-6">
        {/* Market Overview Section */}
        <MarketOverview
          marketData={marketOverview}
          lendingPool={lendingPool}
          borrowingStats={borrowingStats}
          contractData={formatted}
          isLoading={isLoadingAPI || contractLoading}
        />

        {/* User Dashboard Section */}
        {isAuthenticated ? (
          <UserDashboard
            lendingBalance={lendingBalance}
            borrowingBalance={borrowingBalance}
            collateralBalance={collateralBalance}
            isLoading={isLoadingAPI}
          />
        ) : (
          <div className="text-center py-12">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gray-800 mb-4">
              <svg
                className="w-8 h-8 text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                />
              </svg>
            </div>
            <h3 className="text-lg font-semibold mb-2">Connect Your Wallet</h3>
            <p className="text-gray-400 mb-6">
              Connect your wallet to access lending, borrowing, and trading
              features
            </p>
            <ConnectWalletButton autoAuth={true} />
          </div>
        )}

        {/* Protocol Stats Section */}
        <ProtocolStats contractData={formatted} isLoading={contractLoading} />
      </div>
    </div>
  );
}
