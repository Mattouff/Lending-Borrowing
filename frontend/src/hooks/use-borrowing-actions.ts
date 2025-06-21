'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/service/api';
import { toast } from 'sonner';

export function useBorrowingActions() {
  const queryClient = useQueryClient();

  const borrow = useMutation({
    mutationFn: (amount: string) => api.borrow(amount),
    onSuccess: () => {
      toast.success('Borrow successful!');
      queryClient.invalidateQueries({ queryKey: ['borrowing-balance'] });
      queryClient.invalidateQueries({ queryKey: ['borrowing-stats'] });
      queryClient.invalidateQueries({ queryKey: ['collateral-balance'] });
      queryClient.invalidateQueries({ queryKey: ['market-overview'] });
    },
    onError: (error: any) => {
      toast.error(`Borrow failed: ${error.message}`);
    },
  });

  const repay = useMutation({
    mutationFn: (amount: string) => api.repay(amount),
    onSuccess: () => {
      toast.success('Repayment successful!');
      queryClient.invalidateQueries({ queryKey: ['borrowing-balance'] });
      queryClient.invalidateQueries({ queryKey: ['borrowing-stats'] });
      queryClient.invalidateQueries({ queryKey: ['collateral-balance'] });
      queryClient.invalidateQueries({ queryKey: ['market-overview'] });
    },
    onError: (error: any) => {
      toast.error(`Repayment failed: ${error.message}`);
    },
  });

  return {
    borrow,
    repay,
  };
}
