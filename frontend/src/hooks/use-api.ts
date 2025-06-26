'use client';

import { useAuth } from './use-auth';

const API_BASE_URL = 'http://localhost:8080/api/v1';

export function useApi() {
    const { getAuthToken, logout } = useAuth();

    const makeAuthenticatedRequest = async (
        endpoint: string,
        options: RequestInit = {}
    ) => {
        const token = getAuthToken();

        if (!token) {
            throw new Error('No authentication token available');
        }

        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            ...options,
            headers: {
                ...options.headers,
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        // Si le token est expiré ou invalide, déconnecter l'utilisateur
        if (response.status === 401) {
            logout();
            throw new Error('Authentication expired');
        }

        if (!response.ok) {
            throw new Error(`API request failed: ${response.statusText}`);
        }

        return response.json();
    };

    // Exemples de fonctions API spécifiques
    const getUserProfile = () => makeAuthenticatedRequest('/users/profile');

    const getBorrowingBalance = () => makeAuthenticatedRequest('/borrowing/balance');

    const getLendingBalance = () => makeAuthenticatedRequest('/lending/balance');

    return {
        makeAuthenticatedRequest,
        getUserProfile,
        getBorrowingBalance,
        getLendingBalance,
    };
}