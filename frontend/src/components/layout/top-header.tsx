'use client';
import { useAccount, useConnect, useDisconnect, useEnsName } from 'wagmi';

import { ThemeSwitch } from '@/components/theme-switch';
import { Card, CardContent } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { metaMask } from 'wagmi/connectors';

export function TopHeader() {
  const { address } = useAccount();
  const { disconnect } = useDisconnect();
  const { connect } = useConnect();
  const { data: ensName } = useEnsName({ address });

  return (
    <Card className="fixed top-0 left-0 z-50 shadow-sm rounded-md bg-sidebar text-sidebar-foreground mt-3 ml-3 w-[calc(100%-1.5rem)]">
      <CardContent className="flex items-center justify-between h-[var(--header-height)]">
        <h1 className="truncate font-semibold text-2xl">Dapp</h1>
        <div className="flex gap-4">
          {/* Boutons d'actions (connect, settings, etc.) */}
          <ThemeSwitch />
          <Separator orientation="vertical" className="h-6" />
          {/* Affiche le bouton Connect si non connecté, sinon l'adresse avec le bouton Disconnect */}
          {address ? (
            <div className="flex items-center space-x-2">
              <span>{ensName || address}</span>
              <button onClick={() => disconnect()}>Disconnect</button>
            </div>
          ) : (
            <button onClick={() => connect({ connector: metaMask() })}>
              Connect
            </button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
