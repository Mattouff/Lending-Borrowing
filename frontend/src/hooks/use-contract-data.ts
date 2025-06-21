'use client';

import { useReadContracts } from 'wagmi';

// Import des ABIs (vous devez créer ces imports)
import BorrowingABI from '@/contracts/Borrowing.json';
import CollateralABI from '@/contracts/Collateral.json';
import LendingPoolABI from '@/contracts/LendingPool.json';
import TokenABI from '@/contracts/Token.json';


// Adresses des contrats (à mettre à jour avec vos vraies adresses)
const CONTRACT_ADDRESSES = {
    TOKEN: '0x5FbDB2315678afecb367f032d93F642f64180aa3', // Remplacez par votre adresse
    LENDING_POOL: '0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512', // Remplacez par votre adresse
    BORROWING: '0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0', // Remplacez par votre adresse
    COLLATERAL: '0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9', // Remplacez par votre adresse
} as const;

export function useContractData() {
    // Appels multiples pour toutes les données de base
    const { data: contractsData, isLoading, isError, error } = useReadContracts({
        contracts: [
            // Token Contract
            {
                address: CONTRACT_ADDRESSES.TOKEN,
                abi: TokenABI,
                functionName: 'name',
            },
            {
                address: CONTRACT_ADDRESSES.TOKEN,
                abi: TokenABI,
                functionName: 'symbol',
            },
            {
                address: CONTRACT_ADDRESSES.TOKEN,
                abi: TokenABI,
                functionName: 'totalSupply',
            },
            {
                address: CONTRACT_ADDRESSES.TOKEN,
                abi: TokenABI,
                functionName: 'decimals',
            },

            // LendingPool Contract
            {
                address: CONTRACT_ADDRESSES.LENDING_POOL,
                abi: LendingPoolABI,
                functionName: 'name',
            },
            {
                address: CONTRACT_ADDRESSES.LENDING_POOL,
                abi: LendingPoolABI,
                functionName: 'symbol',
            },
            {
                address: CONTRACT_ADDRESSES.LENDING_POOL,
                abi: LendingPoolABI,
                functionName: 'totalSupply',
            },
            {
                address: CONTRACT_ADDRESSES.LENDING_POOL,
                abi: LendingPoolABI,
                functionName: 'totalLending',
            },
            {
                address: CONTRACT_ADDRESSES.LENDING_POOL,
                abi: LendingPoolABI,
                functionName: 'annualInterestRate',
            },

            // Borrowing Contract
            {
                address: CONTRACT_ADDRESSES.BORROWING,
                abi: BorrowingABI,
                functionName: 'totalBorrowed',
            },
            {
                address: CONTRACT_ADDRESSES.BORROWING,
                abi: BorrowingABI,
                functionName: 'getCurrentRate',
            },
            {
                address: CONTRACT_ADDRESSES.BORROWING,
                abi: BorrowingABI,
                functionName: 'rMin',
            },
            {
                address: CONTRACT_ADDRESSES.BORROWING,
                abi: BorrowingABI,
                functionName: 'rMax',
            },

            // Collateral Contract  
            {
                address: CONTRACT_ADDRESSES.COLLATERAL,
                abi: CollateralABI,
                functionName: 'LIQUIDATION_THRESHOLD',
            },
            {
                address: CONTRACT_ADDRESSES.COLLATERAL,
                abi: CollateralABI,
                functionName: 'MIN_COLLATERAL_RATIO',
            },
        ],
    });

    // Parsing des résultats
    const [
        tokenName,
        tokenSymbol,
        tokenTotalSupply,
        tokenDecimals,
        lendingPoolName,
        lendingPoolSymbol,
        lendingPoolTotalSupply,
        totalLending,
        annualInterestRate,
        totalBorrowed,
        currentBorrowRate,
        rMin,
        rMax,
        liquidationThreshold,
        minCollateralRatio,
    ] = contractsData || [];

    return {
        isLoading,
        isError,
        error,

        // Token Data
        token: {
            name: tokenName?.result as string,
            symbol: tokenSymbol?.result as string,
            totalSupply: tokenTotalSupply?.result as bigint,
            decimals: tokenDecimals?.result as number,
        },

        // LendingPool Data
        lendingPool: {
            name: lendingPoolName?.result as string,
            symbol: lendingPoolSymbol?.result as string,
            totalSupply: lendingPoolTotalSupply?.result as bigint,
            totalLending: totalLending?.result as bigint,
            annualInterestRate: annualInterestRate?.result as bigint,
        },

        // Borrowing Data
        borrowing: {
            totalBorrowed: totalBorrowed?.result as bigint,
            currentRate: currentBorrowRate?.result as bigint,
            rMin: rMin?.result as bigint,
            rMax: rMax?.result as bigint,
        },

        // Collateral Data
        collateral: {
            liquidationThreshold: liquidationThreshold?.result as bigint,
            minCollateralRatio: minCollateralRatio?.result as bigint,
        },

        // Adresses des contrats
        addresses: CONTRACT_ADDRESSES,
    };
}