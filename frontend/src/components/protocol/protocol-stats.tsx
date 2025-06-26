'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Loader2 } from 'lucide-react';

interface ProtocolStatsProps {
  contractData: any;
  isLoading: boolean;
}

export function ProtocolStats({ contractData, isLoading }: ProtocolStatsProps) {
  const parameters = [
    {
      label: 'Token Name',
      value: contractData?.token?.name || 'Loading...',
      description: 'Protocol token name',
    },
    {
      label: 'Token Symbol',
      value: contractData?.token?.symbol || 'TKN',
      description: 'Protocol token symbol',
    },
    {
      label: 'Total Supply',
      value: contractData?.token?.totalSupply || '0',
      description: 'Total token supply',
    },
    {
      label: 'Lending Pool Supply',
      value: contractData?.lendingPool?.totalSupply || '0',
      description: 'LP token supply',
    },
    {
      label: 'Interest Rate',
      value: contractData?.lendingPool?.interestRate || '0%',
      description: 'Annual lending rate',
    },
    {
      label: 'Current Borrow Rate',
      value: contractData?.borrowing?.currentRate || '0%',
      description: 'Dynamic borrowing rate',
    },
    {
      label: 'Rate Range',
      value: `${contractData?.borrowing?.rMin || '0%'} - ${contractData?.borrowing?.rMax || '0%'}`,
      description: 'Min/Max borrow rates',
    },
    {
      label: 'Liquidation Threshold',
      value: contractData?.collateral?.liquidationThreshold || '80%',
      description: 'Liquidation trigger point',
    },
  ];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Protocol Parameters</h2>
        <Badge variant="outline" className=" text-gray-300 border-gray-700">
          Smart Contract Data
        </Badge>
      </div>

      <Card className=" border-gray-800">
        <CardHeader>
          <CardTitle className="text-lg">Live Contract Data</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {parameters.map((param, index) => (
              <div key={index} className="space-y-2">
                <div className="text-xs font-medium text-gray-400 uppercase tracking-wider">
                  {param.label}
                </div>
                <div className="text-lg font-mono text-white">
                  {isLoading ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    param.value
                  )}
                </div>
                <div className="text-xs text-gray-500">{param.description}</div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
