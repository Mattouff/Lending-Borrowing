'use client';

import { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Loader2, TrendingDown, AlertCircle, DollarSign } from 'lucide-react';
import { useBorrowingActions } from '@/hooks/use-borrowing-actions';

interface BorrowingModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  type: 'borrow' | 'repay';
  maxAmount?: string;
  currentDebt?: string;
  rate?: string;
}

export function BorrowingModal({ 
  open, 
  onOpenChange, 
  type, 
  maxAmount = '0',
  currentDebt = '0',
  rate = '0' 
}: BorrowingModalProps) {
  const [amount, setAmount] = useState('');
  const { borrow, repay } = useBorrowingActions();

  const isBorrow = type === 'borrow';
  const mutation = isBorrow ? borrow : repay;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!amount || parseFloat(amount) <= 0) return;
    
    mutation.mutate(amount, {
      onSuccess: () => {
        setAmount('');
        onOpenChange(false);
      }
    });
  };

  const setMaxAmount = () => {
    setAmount(isBorrow ? maxAmount : currentDebt);
  };

  const calculateInterest = () => {
    if (!amount) return '0.00';
    const principal = parseFloat(amount);
    const rateDecimal = parseFloat(rate) / 100;
    return (principal * rateDecimal / 365).toFixed(6); // Daily interest
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="bg-gray-950 border-gray-800 text-white">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <TrendingDown className="h-5 w-5 text-red-400" />
            <span>{isBorrow ? 'Borrow' : 'Repay'} USDC</span>
            <Badge variant="outline" className="bg-red-500/10 text-red-400 border-red-500/20">
              APR {rate}%
            </Badge>
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Current Position */}
          <div className="bg-gray-900/50 rounded-lg p-4">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm text-gray-400">Current Debt</span>
              <span className="text-sm font-mono text-red-400">${currentDebt}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-sm text-gray-400">
                {isBorrow ? 'Available to Borrow' : 'Max Repayment'}
              </span>
              <span className="text-sm font-mono text-white">
                ${isBorrow ? maxAmount : currentDebt}
              </span>
            </div>
          </div>

          {/* Amount Input */}
          <div className="space-y-2">
            <Label htmlFor="amount" className="text-gray-300">Amount</Label>
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
                max={isBorrow ? maxAmount : currentDebt}
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

          {/* Borrowing Details */}
          {isBorrow && (
            <div className="bg-gray-900/50 rounded-lg p-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Daily Interest</span>
                <span className="text-red-400 font-mono">${calculateInterest()}</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Borrowing APR</span>
                <span className="text-red-400">{rate}%</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">You will receive</span>
                <span className="text-green-400 font-mono">{amount || '0.00'} USDC</span>
              </div>
            </div>
          )}

          {/* Repayment Details */}
          {!isBorrow && (
            <div className="bg-gray-900/50 rounded-lg p-4 space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Current APR</span>
                <span className="text-red-400">{rate}%</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Remaining Debt</span>
                <span className="text-red-400 font-mono">
                  ${(parseFloat(currentDebt) - parseFloat(amount || '0')).toFixed(2)}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Interest Saved (Daily)</span>
                <span className="text-green-400 font-mono">-${calculateInterest()}</span>
              </div>
            </div>
          )}

          {/* Health Factor Warning */}
          {isBorrow && (
            <div className="flex items-start space-x-2 p-3 bg-yellow-500/10 border border-yellow-500/20 rounded-lg">
              <AlertCircle className="h-4 w-4 text-yellow-400 mt-0.5" />
              <div className="text-sm">
                <p className="text-yellow-400 font-medium">Collateral Required</p>
                <p className="text-yellow-300">
                  Ensure you have sufficient collateral to maintain a healthy position ratio.
                </p>
              </div>
            </div>
          )}

          {/* Risk Warning for large borrowing */}
          {isBorrow && amount && parseFloat(amount) > parseFloat(maxAmount) * 0.8 && (
            <div className="flex items-start space-x-2 p-3 bg-red-500/10 border border-red-500/20 rounded-lg">
              <AlertCircle className="h-4 w-4 text-red-400 mt-0.5" />
              <div className="text-sm">
                <p className="text-red-400 font-medium">High Risk Warning</p>
                <p className="text-red-300">
                  Borrowing this amount will significantly increase your liquidation risk.
                </p>
              </div>
            </div>
          )}

          {/* Submit Button */}
          <Button
            type="submit"
            className={`w-full ${isBorrow ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700'} text-white`}
            disabled={!amount || parseFloat(amount) <= 0 || mutation.isPending}
          >
            {mutation.isPending ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin mr-2" />
                Processing...
              </>
            ) : (
              <>
                <DollarSign className="h-4 w-4 mr-2" />
                {isBorrow ? 'Borrow' : 'Repay'} {amount || '0.00'} USDC
              </>
            )}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  );
}
