'use client';

import {
  ConnectWalletButton,
  useSimpleKit,
} from '@/components/simplekit/simplekit';
import { Button } from '@/components/ui/button';
import { useAuth } from '@/hooks/use-auth';
import { useAccount } from 'wagmi';

export function AuthButton() {
  const { isConnected } = useAccount();
  const { isAuthenticated, isAuthenticating, authError, login, logout } =
    useAuth();
  const simplekit = useSimpleKit();

  if (!isConnected) {
    return <ConnectWalletButton />;
  }

  if (!isAuthenticated) {
    return (
      <div className="flex flex-col gap-2">
        <Button
          onClick={login}
          disabled={isAuthenticating}
          className="rounded-xl"
        >
          {isAuthenticating ? 'Authenticating...' : 'Sign Message to Login'}
        </Button>
        {authError && <p className="text-sm text-red-500">{authError}</p>}
      </div>
    );
  }

  return (
    <div className="flex gap-2">
      <Button
        onClick={simplekit.toggleModal}
        variant="outline"
        className="rounded-xl"
      >
        {simplekit.formattedAddress}
      </Button>
      <Button onClick={logout} variant="destructive" className="rounded-xl">
        Logout
      </Button>
    </div>
  );
}
