'use client';

import { LendingModal } from '@/components/actions/lending-modal';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  AlertTriangle,
  Loader2,
  Minus,
  Plus,
  Shield,
  TrendingDown,
  TrendingUp,
} from 'lucide-react';
import { useState } from 'react';
import { BorrowingModal } from '../actions/borrowing-modal';
import { CollateralModal } from '../actions/collateral-modal';

interface UserDashboardProps {
  lendingBalance: any;
  borrowingBalance: any;
  collateralBalance: any;
  isLoading: boolean;
}

export function UserDashboard({
  lendingBalance,
  borrowingBalance,
  collateralBalance,
  isLoading,
}: UserDashboardProps) {
  // Modal states
  const [lendingModal, setLendingModal] = useState<{
    open: boolean;
    type: 'deposit' | 'withdraw';
  }>({
    open: false,
    type: 'deposit',
  });
  const [borrowingModal, setBorrowingModal] = useState<{
    open: boolean;
    type: 'borrow' | 'repay';
  }>({
    open: false,
    type: 'borrow',
  });
  const [collateralModal, setCollateralModal] = useState<{
    open: boolean;
    type: 'deposit' | 'withdraw';
  }>({
    open: false,
    type: 'deposit',
  });

  const positions = [
    {
      title: 'Lending Position',
      value: lendingBalance?.balance || '0',
      apy: lendingBalance?.apy || '5.2',
      icon: <TrendingUp className="h-5 w-5 text-green-400" />,
      type: 'lending',
      color: 'border-green-500/20 bg-green-500/5',
      actions: [
        {
          label: 'Deposit',
          icon: Plus,
          variant: 'default' as const,
          onClick: () => setLendingModal({ open: true, type: 'deposit' }),
        },
        {
          label: 'Withdraw',
          icon: Minus,
          variant: 'outline' as const,
          onClick: () => setLendingModal({ open: true, type: 'withdraw' }),
        },
      ],
    },
    {
      title: 'Borrowed Amount',
      value: borrowingBalance?.balance || '0',
      rate: borrowingBalance?.rate || '7.8',
      icon: <TrendingDown className="h-5 w-5 text-red-400" />,
      type: 'borrowing',
      color: 'border-red-500/20 bg-red-500/5',
      actions: [
        {
          label: 'Borrow',
          icon: Plus,
          variant: 'default' as const,
          onClick: () => setBorrowingModal({ open: true, type: 'borrow' }),
        },
        {
          label: 'Repay',
          icon: Minus,
          variant: 'outline' as const,
          onClick: () => setBorrowingModal({ open: true, type: 'repay' }),
        },
      ],
    },
    {
      title: 'Collateral',
      value: collateralBalance?.balance || '0',
      healthFactor: collateralBalance?.healthFactor || '2.45',
      icon: <Shield className="h-5 w-5 text-blue-400" />,
      type: 'collateral',
      color: 'border-blue-500/20 bg-blue-500/5',
      actions: [
        {
          label: 'Add',
          icon: Plus,
          variant: 'default' as const,
          onClick: () => setCollateralModal({ open: true, type: 'deposit' }),
        },
        {
          label: 'Remove',
          icon: Minus,
          variant: 'outline' as const,
          onClick: () => setCollateralModal({ open: true, type: 'withdraw' }),
        },
      ],
    },
  ];

  const getHealthColor = (factor: string) => {
    const value = parseFloat(factor);
    if (value >= 2) return 'text-green-400';
    if (value >= 1.5) return 'text-yellow-400';
    return 'text-red-400';
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Your Positions</h2>
        <div className="flex items-center space-x-2">
          <Badge
            variant="outline"
            className="bg-blue-500/10 text-blue-400 border-blue-500/20"
          >
            Portfolio
          </Badge>
        </div>
      </div>

      {/* Positions Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {positions.map((position, index) => (
          <Card
            key={index}
            className={`bg-gray-950/50 border-gray-800 ${position.color}`}
          >
            <CardHeader className="pb-3">
              <CardTitle className="text-sm font-medium text-gray-300 flex items-center justify-between">
                {position.title}
                {position.icon}
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {isLoading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <>
                  <div className="space-y-2">
                    <div className="text-2xl font-bold font-mono text-white">
                      ${parseFloat(position.value).toFixed(2)}
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      {position.apy && (
                        <span className="text-green-400">
                          APY: {position.apy}%
                        </span>
                      )}
                      {position.rate && (
                        <span className="text-red-400">
                          APR: {position.rate}%
                        </span>
                      )}
                      {position.healthFactor && (
                        <span
                          className={`font-medium ${getHealthColor(
                            position.healthFactor
                          )}`}
                        >
                          Health: {position.healthFactor}
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="flex space-x-2">
                    {position.actions.map((action, actionIndex) => (
                      <Button
                        key={actionIndex}
                        size="sm"
                        className="flex-1 h-8 text-xs"
                        variant={action.variant}
                        onClick={action.onClick}
                      >
                        <action.icon className="h-3 w-3 mr-1" />
                        {action.label}
                      </Button>
                    ))}
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Risk Alert */}
      {collateralBalance?.healthFactor &&
        parseFloat(collateralBalance.healthFactor) < 1.5 && (
          <Card className="bg-red-500/10 border-red-500/20">
            <CardContent className="p-4">
              <div className="flex items-center space-x-3">
                <AlertTriangle className="h-5 w-5 text-red-400" />
                <div>
                  <p className="text-sm font-medium text-red-400">
                    Low Health Factor
                  </p>
                  <p className="text-xs text-red-300">
                    Your position is at risk of liquidation. Consider adding
                    collateral or repaying debt.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        )}

      {/* Modals */}
      <LendingModal
        open={lendingModal.open}
        onOpenChange={(open) => setLendingModal((prev) => ({ ...prev, open }))}
        type={lendingModal.type}
        maxAmount={lendingBalance?.availableToWithdraw || '1000'}
        currentBalance={lendingBalance?.balance || '0'}
        apy={lendingBalance?.apy || '5.2'}
      />

      <BorrowingModal
        open={borrowingModal.open}
        onOpenChange={(open) =>
          setBorrowingModal((prev) => ({ ...prev, open }))
        }
        type={borrowingModal.type}
        maxAmount={borrowingBalance?.availableToBorrow || '500'}
        currentDebt={borrowingBalance?.balance || '0'}
        rate={borrowingBalance?.rate || '7.8'}
      />

      <CollateralModal
        open={collateralModal.open}
        onOpenChange={(open) =>
          setCollateralModal((prev) => ({ ...prev, open }))
        }
        type={collateralModal.type}
        maxAmount={collateralBalance?.availableToWithdraw || '100'}
        currentBalance={collateralBalance?.balance || '0'}
      />
    </div>
  );
}
