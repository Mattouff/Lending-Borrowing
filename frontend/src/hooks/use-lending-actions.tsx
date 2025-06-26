'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/service/api';
import { toast } from 'sonner'; // Install: npm install sonner

export function useLendingActions() {
  const queryClient = useQueryClient();

  const deposit = useMutation({
    mutationFn: (amount: string) => api.deposit(amount),
    onSuccess: () => {
      toast.success('Deposit successful!');
      queryClient.invalidateQueries({ queryKey: ['lending-balance'] });
      queryClient.invalidateQueries({ queryKey: ['lending-pool-info'] });
    },
    onError: (error: any) => {
      toast.error(`Deposit failed: ${error.message}`);
    },
  });

  const withdraw = useMutation({
    mutationFn: (amount: string) => api.withdraw(amount),
    onSuccess: () => {
      toast.success('Withdrawal successful!');
      queryClient.invalidateQueries({ queryKey: ['lending-balance'] });
      queryClient.invalidateQueries({ queryKey: ['lending-pool-info'] });
    },
    onError: (error: any) => {
      toast.error(`Withdrawal failed: ${error.message}`);
    },
  });

  return {
    deposit,
    withdraw,
  };
}