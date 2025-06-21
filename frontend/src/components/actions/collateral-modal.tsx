'use client';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useCollateralActions } from '@/hooks/use-collateral-actions';
import { AlertCircle, Loader2, Shield, TrendingUp } from 'lucide-react';
import { useState } from 'react';

interface CollateralModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  type: 'deposit' | 'withdraw';
  maxAmount?: string;
  currentBalance?: string;
  healthFactor?: string;
}

export function CollateralModal({
  open,
  onOpenChange,
  type,
  maxAmount = '0',
  currentBalance = '0',
  healthFactor = '2.45',
}: CollateralModalProps) {
  const [amount, setAmount] = useState('');
  const { depositCollateral, withdrawCollateral } = useCollateralActions();

  const isDeposit = type === 'deposit';
  const mutation = isDeposit ? depositCollateral : withdrawCollateral;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!amount || parseFloat(amount) <= 0) return;

    mutation.mutate(amount, {
      onSuccess: () => {
        setAmount('');
        onOpenChange(false);
      },
    });
  };

  const setMaxAmount = () => {
    setAmount(maxAmount);
  };

  const calculateNewHealthFactor = () => {
    if (!amount) return healthFactor;

    const currentCol = parseFloat(currentBalance);
    const amountNum = parseFloat(amount);
    const currentHealth = parseFloat(healthFactor);

    // Simplified calculation - in reality this would be more complex
    const newCollateral = isDeposit
      ? currentCol + amountNum
      : currentCol - amountNum;
    const ratio = newCollateral / Math.max(currentCol, 1);
    const newHealth = currentHealth * ratio;

    return Math.max(newHealth, 0).toFixed(2);
  };

  const getHealthColor = (factor: string) => {
    const value = parseFloat(factor);
    if (value >= 2) return 'text-green-400';
    if (value >= 1.5) return 'text-yellow-400';
    return 'text-red-400';
  };

  const newHealthFactor = calculateNewHealthFactor();
  const isRiskyWithdrawal = !isDeposit && parseFloat(newHealthFactor) < 1.5;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="bg-gray-950 border-gray-800 text-white">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Shield className="h-5 w-5 text-blue-400" />
            <span>{isDeposit ? 'Add' : 'Remove'} Collateral</span>
            <Badge
              variant="outline"
              className="bg-blue-500/10 text-blue-400 border-blue-500/20"
            >
              ETH
            </Badge>
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Current Position */}
          <div className="bg-gray-900/50 rounded-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm text-gray-400">Current Collateral</span>
              <span className="text-sm font-mono text-blue-400">
                {currentBalance} ETH
              </span>
            </div>
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm text-gray-400">Health Factor</span>
              <span
                className={`text-sm font-mono ${getHealthColor(healthFactor)}`}
              >
                {healthFactor}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">
                {isDeposit ? 'Available to Deposit' : 'Available to Withdraw'}
              </span>
              <span className="text-sm font-mono text-white">
                {maxAmount} ETH
              </span>
            </div>
          </div>

          {/* Amount Input */}
          <div className="space-y-2">
            <Label htmlFor="amount" className="text-gray-300">
              Amount
            </Label>
            <div className="relative">
              <Input
                id="amount"
                type="number"
                placeholder="0.00"
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
                className="bg-gray-900 border-gray-700 text-white pr-16"
                step="0.001"
                min="0"
                max={maxAmount}
              />
              <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center space-x-2">
                <button
                  type="button"
                  onClick={setMaxAmount}
                  className="text-xs text-blue-400 hover:text-blue-300"
                >
                  MAX
                </button>
                <span className="text-sm text-gray-400">ETH</span>
              </div>
            </div>
          </div>

          {/* Health Factor Impact */}
          <div className="bg-gray-900/50 rounded-lg p-4 space-y-3">
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">Current Health Factor</span>
              <span className={`font-mono ${getHealthColor(healthFactor)}`}>
                {healthFactor}
              </span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">New Health Factor</span>
              <span className={`font-mono ${getHealthColor(newHealthFactor)}`}>
                {newHealthFactor}
              </span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">
                {isDeposit ? 'New' : 'Remaining'} Collateral
              </span>
              <span className="text-blue-400 font-mono">
                {isDeposit
                  ? (
                      parseFloat(currentBalance) + parseFloat(amount || '0')
                    ).toFixed(3)
                  : (
                      parseFloat(currentBalance) - parseFloat(amount || '0')
                    ).toFixed(3)}{' '}
                ETH
              </span>
            </div>
          </div>

          {/* Warnings */}
          {isRiskyWithdrawal && (
            <div className="flex items-start space-x-2 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <AlertCircle className="h-4 w-4 text-red-400 mt-0.5" />
              <div className="text-sm">
                <p className="text-red-400 font-medium">Liquidation Risk</p>
                <p className="text-red-300">
                  This withdrawal will put your position at high risk of
                  liquidation. Health factor will drop to {newHealthFactor}.
                </p>
              </div>
            </div>
          )}

          {isDeposit && (
            <div className="flex items-start space-x-2 p-3 bg-green-500/10 border border-green-500/20 rounded-lg">
              <TrendingUp className="h-4 w-4 text-green-400 mt-0.5" />
              <div className="text-sm">
                <p className="text-green-400 font-medium">Improved Health</p>
                <p className="text-green-300">
                  Adding collateral will improve your health factor and
                  borrowing capacity.
                </p>
              </div>
            </div>
          )}

          {/* Liquidation threshold warning */}
          {parseFloat(newHealthFactor) < 1.2 && !isDeposit && (
            <div className="flex items-start space-x-2 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <AlertCircle className="h-4 w-4 text-red-400 mt-0.5" />
              <div className="text-sm">
                <p className="text-red-400 font-medium">Critical Warning</p>
                <p className="text-red-300">
                  Your position may be liquidated immediately. Consider repaying
                  debt instead.
                </p>
              </div>
            </div>
          )}

          {/* Submit Button */}
          <Button
            type="submit"
            className={`w-full ${
              isDeposit
                ? 'bg-blue-600 hover:bg-blue-700'
                : isRiskyWithdrawal
                  ? 'bg-red-600 hover:bg-red-700'
                  : 'bg-gray-600 hover:bg-gray-700'
            } text-white`}
            disabled={
              !amount ||
              parseFloat(amount) <= 0 ||
              mutation.isPending ||
              (parseFloat(newHealthFactor) < 1.1 && !isDeposit)
            }
          >
            {mutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                Processing...
              </>
            ) : (
              <>
                <Shield className="h-4 w-4 mr-2" />
                {isDeposit ? 'Add' : 'Remove'} {amount || '0.00'} ETH
              </>
            )}
          </Button>

          {/* Transaction fee info */}
          <div className="text-xs text-gray-400 text-center">
            Network fee: ~$2-5 depending on gas price
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
