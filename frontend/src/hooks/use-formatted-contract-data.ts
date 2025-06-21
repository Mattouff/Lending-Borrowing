'use client';

import { formatUnits } from 'viem';
import { useContractData } from './use-contract-data';

export function useFormattedContractData() {
    const contractData = useContractData();

    const formatBigInt = (value: bigint | undefined, decimals: number = 18) => {
        if (!value) return '0';
        try {
            return parseFloat(formatUnits(value, decimals)).toFixed(4);
        } catch {
            return '0';
        }
    };

    const formatPercentage = (value: bigint | undefined, decimals: number = 18) => {
        if (!value) return '0%';
        try {
            const formatted = parseFloat(formatUnits(value, decimals));
            return `${(formatted * 100).toFixed(2)}%`;
        } catch {
            return '0%';
        }
    };

    return {
        ...contractData,

        // Données formatées pour l'affichage
        formatted: {
            token: {
                name: contractData.token.name || 'Loading...',
                symbol: contractData.token.symbol || 'TKN',
                totalSupply: formatBigInt(contractData.token.totalSupply, contractData.token.decimals || 18),
                totalSupplyRaw: contractData.token.totalSupply,
            },

            lendingPool: {
                name: contractData.lendingPool.name || 'Lending Pool',
                symbol: contractData.lendingPool.symbol || 'LP',
                totalSupply: formatBigInt(contractData.lendingPool.totalSupply),
                totalLending: formatBigInt(contractData.lendingPool.totalLending),
                interestRate: formatPercentage(contractData.lendingPool.annualInterestRate),
            },

            borrowing: {
                totalBorrowed: formatBigInt(contractData.borrowing.totalBorrowed),
                currentRate: formatPercentage(contractData.borrowing.currentRate),
                rMin: formatPercentage(contractData.borrowing.rMin),
                rMax: formatPercentage(contractData.borrowing.rMax),
            },

            collateral: {
                liquidationThreshold: formatPercentage(contractData.collateral.liquidationThreshold),
                minCollateralRatio: formatPercentage(contractData.collateral.minCollateralRatio),
            },
        },
    };
}