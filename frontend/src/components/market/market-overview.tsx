'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import {
  Activity,
  BarChart3,
  DollarSign,
  Loader2,
  TrendingDown,
  TrendingUp,
  Users,
} from 'lucide-react';

interface MarketOverviewProps {
  marketData: any;
  lendingPool: any;
  borrowingStats: any;
  contractData: any;
  isLoading: boolean;
}

export function MarketOverview({
  marketData,
  lendingPool,
  borrowingStats,
  contractData,
  isLoading,
}: MarketOverviewProps) {
  // Format numbers like HyperLiquid
  const formatValue = (value: any) => {
    if (!value) return '0';
    const num = parseFloat(value.toString());
    if (num >= 1e9) return `$${(num / 1e9).toFixed(2)}B`;
    if (num >= 1e6) return `$${(num / 1e6).toFixed(2)}M`;
    if (num >= 1e3) return `$${(num / 1e3).toFixed(2)}K`;
    return `$${num.toFixed(2)}`;
  };

  const stats = [
    {
      label: 'TVL',
      value:
        marketData?.tvl || contractData?.token?.totalSupply || '12,500,000',
      change: '+5.2%',
      icon: <DollarSign className="h-4 w-4" />,
      color: 'text-green-400',
    },
    {
      label: 'Total Supplied',
      value:
        lendingPool?.totalSupplied ||
        contractData?.lendingPool?.totalLending ||
        '8,750,000',
      change: '+3.1%',
      icon: <TrendingUp className="h-4 w-4" />,
      color: 'text-blue-400',
    },
    {
      label: 'Total Borrowed',
      value:
        borrowingStats?.totalBorrowed ||
        contractData?.borrowing?.totalBorrowed ||
        '6,200,000',
      change: '+2.8%',
      icon: <TrendingDown className="h-4 w-4" />,
      color: 'text-orange-400',
    },
    {
      label: 'Utilization',
      value: borrowingStats?.utilizationRate || '71',
      suffix: '%',
      change: '+1.2%',
      icon: <Activity className="h-4 w-4" />,
      color: 'text-purple-400',
    },
    {
      label: 'Active Users',
      value: marketData?.activeUsers || '1,234',
      change: '+12%',
      icon: <Users className="h-4 w-4" />,
      color: 'text-cyan-400',
    },
    {
      label: 'Lending APY',
      value:
        lendingPool?.apy || contractData?.lendingPool?.interestRate || '5.2',
      suffix: '%',
      change: '+0.1%',
      icon: <BarChart3 className="h-4 w-4" />,
      color: 'text-green-400',
    },
  ];

  return (
    <div className="space-y-4">
      {/* Title */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Market Overview</h2>
        <Badge
          variant="outline"
          className="bg-green-500/10 text-green-400 border-green-500/20"
        >
          Live
        </Badge>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
        {stats.map((stat, index) => (
          <Card
            key={index}
            className=" border-gray-800 hover:border-gray-700 transition-all duration-200"
          >
            <CardContent className="p-4">
              <div className="flex items-center justify-between mb-2">
                <span className="text-xs font-medium text-gray-400 uppercase tracking-wider">
                  {stat.label}
                </span>
                <div className={stat.color}>{stat.icon}</div>
              </div>

              <div className="space-y-1">
                {isLoading ? (
                  <Loader2 className="h-4 w-4 animate-spin text-gray-400" />
                ) : (
                  <>
                    <div className="text-lg font-bold font-mono text-white">
                      {stat.suffix && !stat.value.includes('$')
                        ? `${stat.value}${stat.suffix}`
                        : stat.value.includes('$')
                          ? stat.value
                          : formatValue(stat.value)}
                    </div>
                    <div className="text-xs text-green-400 font-medium">
                      {stat.change}
                    </div>
                  </>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
