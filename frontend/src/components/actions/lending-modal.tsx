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
import { useLendingActions } from '@/hooks/use-lending-actions';
import { AlertCircle, Loader2, TrendingUp } from 'lucide-react';
import { useState } from 'react';

interface LendingModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  type: 'deposit' | 'withdraw';
  maxAmount?: string;
  currentBalance?: string;
  apy?: string;
}

export function LendingModal({
  open,
  onOpenChange,
  type,
  maxAmount = '0',
  currentBalance = '0',
  apy = '0',
}: LendingModalProps) {
  const [amount, setAmount] = useState('');
  const { deposit, withdraw } = useLendingActions();

  const isDeposit = type === 'deposit';
  const mutation = isDeposit ? deposit : withdraw;

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

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="bg-gray-950 border-gray-800 text-white">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <TrendingUp className="h-5 w-5 text-green-400" />
            <span>{isDeposit ? 'Deposit' : 'Withdraw'} USDC</span>
            <Badge
              variant="outline"
              className="bg-green-500/10 text-green-400 border-green-500/20"
            >
              APY {apy}%
            </Badge>
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Current Position */}
          <div className="bg-gray-900/50 rounded-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm text-gray-400">Current Position</span>
              <span className="text-sm font-mono text-white">
                ${currentBalance}
              </span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">Available</span>
              <span className="text-sm font-mono text-green-400">
                ${maxAmount}
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
                className="bg-gray-900 border-gray-700 text-white pr-20"
                step="0.01"
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
                <span className="text-sm text-gray-400">USDC</span>
              </div>
            </div>
          </div>

          {/* Transaction Info */}
          <div className="bg-gray-900/50 rounded-lg p-4 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">Transaction Fee</span>
              <span className="text-white">~$0.05</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-400">
                You will {isDeposit ? 'receive' : 'withdraw'}
              </span>
              <span className="text-green-400 font-mono">
                {amount || '0.00'} {isDeposit ? 'LP tokens' : 'USDC'}
              </span>
            </div>
          </div>

          {/* Warning */}
          {!isDeposit && (
            <div className="flex items-start space-x-2 p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg">
              <AlertCircle className="h-4 w-4 text-yellow-400 mt-0.5" />
              <div className="text-sm">
                <p className="text-yellow-400 font-medium">Withdrawal Notice</p>
                <p className="text-yellow-300">
                  Withdrawals may take up to 24h to process depending on pool
                  utilization.
                </p>
              </div>
            </div>
          )}

          {/* Submit Button */}
          <Button
            type="submit"
            className="w-full bg-blue-600 hover:bg-blue-700 text-white"
            disabled={!amount || parseFloat(amount) <= 0 || mutation.isPending}
          >
            {mutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                Processing...
              </>
            ) : (
              `${isDeposit ? 'Deposit' : 'Withdraw'} ${amount || '0.00'} USDC`
            )}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
