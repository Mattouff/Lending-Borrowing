'use client';

import { useEffect, useState } from 'react';
import { useAccount, useSignMessage } from 'wagmi';

interface AuthResponse {
    token: string;
    expires_at: string;
}

interface NonceResponse {
    message: string;
    nonce: string;
}

const API_BASE_URL = 'http://localhost:8080/api/v1';

export function useAuth() {
    const { address, isConnected } = useAccount();
    const { signMessageAsync } = useSignMessage();
    const [isAuthenticating, setIsAuthenticating] = useState(false);
    const [authError, setAuthError] = useState<string | null>(null);

    // Vérifier si l'utilisateur est déjà authentifié au chargement
    const [isAuthenticated, setIsAuthenticated] = useState(() => {
        if (typeof window !== 'undefined') {
            const token = localStorage.getItem('auth_token');
            const expiresAt = localStorage.getItem('auth_expires_at');

            if (token && expiresAt) {
                const now = new Date();
                const expiry = new Date(expiresAt);
                return now < expiry;
            }
        }
        return false;
    });

    // Fonction pour obtenir un nonce
    const getNonce = async (walletAddress: string): Promise<NonceResponse> => {
        const response = await fetch(`${API_BASE_URL}/users/nonce/${walletAddress}`);
        if (!response.ok) {
            throw new Error('Failed to get nonce');
        }
        return response.json();
    };

    // Fonction pour authentifier avec la signature
    const authenticate = async (walletAddress: string, signature: string): Promise<AuthResponse> => {
        const response = await fetch(`${API_BASE_URL}/users/auth`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                address: walletAddress,
                signature: signature,
            }),
        });

        if (!response.ok) {
            throw new Error('Authentication failed');
        }
        return response.json();
    };

    // Fonction principale de connexion
    const login = async () => {
        if (!address || !isConnected) {
            setAuthError('Wallet not connected');
            return;
        }

        setIsAuthenticating(true);
        setAuthError(null);

        try {
            // 1. Obtenir le nonce
            const { message } = await getNonce(address);

            // 2. Signer le message
            const signature = await signMessageAsync({ message });

            // 3. Authentifier avec la signature
            const authResponse = await authenticate(address, signature);

            // 4. Stocker le token dans localStorage
            localStorage.setItem('auth_token', authResponse.token);
            localStorage.setItem('auth_expires_at', authResponse.expires_at);
            localStorage.setItem('wallet_address', address);

            setIsAuthenticated(true);
        } catch (error) {
            console.error('Authentication error:', error);
            setAuthError(error instanceof Error ? error.message : 'Authentication failed');
        } finally {
            setIsAuthenticating(false);
        }
    };

    // Fonction de déconnexion
    const logout = () => {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_expires_at');
        localStorage.removeItem('wallet_address');
        setIsAuthenticated(false);
    };

    // Obtenir le token pour les requêtes API
    const getAuthToken = (): string | null => {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('auth_token');
        }
        return null;
    };

    // Vérifier l'expiration du token
    useEffect(() => {
        const checkTokenExpiry = () => {
            if (typeof window !== 'undefined') {
                const expiresAt = localStorage.getItem('auth_expires_at');
                if (expiresAt) {
                    const now = new Date();
                    const expiry = new Date(expiresAt);
                    if (now >= expiry) {
                        logout();
                    }
                }
            }
        };

        // Vérifier toutes les minutes
        const interval = setInterval(checkTokenExpiry, 60000);
        return () => clearInterval(interval);
    }, []);

    // Se déconnecter automatiquement si le wallet est déconnecté
    useEffect(() => {
        if (!isConnected && isAuthenticated) {
            logout();
        }
    }, [isConnected, isAuthenticated]);

    return {
        isAuthenticated,
        isAuthenticating,
        authError,
        login,
        logout,
        getAuthToken,
    };
}