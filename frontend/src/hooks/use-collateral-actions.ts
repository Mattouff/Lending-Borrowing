'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/service/api';
import { toast } from 'sonner';

export function useCollateralActions() {
  const queryClient = useQueryClient();

  const depositCollateral = useMutation({
    mutationFn: (amount: string) => api.depositCollateral(amount),
    onSuccess: () => {
      toast.success('Collateral deposit successful!');
      queryClient.invalidateQueries({ queryKey: ['collateral-balance'] });
      queryClient.invalidateQueries({ queryKey: ['collateral-info'] });
      queryClient.invalidateQueries({ queryKey: ['borrowing-balance'] });
      queryClient.invalidateQueries({ queryKey: ['market-overview'] });
    },
    onError: (error: any) => {
      toast.error(`Collateral deposit failed: ${error.message}`);
    },
  });

  const withdrawCollateral = useMutation({
    mutationFn: (amount: string) => api.withdrawCollateral(amount),
    onSuccess: () => {
      toast.success('Collateral withdrawal successful!');
      queryClient.invalidateQueries({ queryKey: ['collateral-balance'] });
      queryClient.invalidateQueries({ queryKey: ['collateral-info'] });
      queryClient.invalidateQueries({ queryKey: ['borrowing-balance'] });
      queryClient.invalidateQueries({ queryKey: ['market-overview'] });
    },
    onError: (error: any) => {
      toast.error(`Collateral withdrawal failed: ${error.message}`);
    },
  });

  return {
    depositCollateral,
    withdrawCollateral,
  };
}
