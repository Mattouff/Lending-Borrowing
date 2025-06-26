'use client';

import * as React from 'react';

import {
  SimpleKitModal,
  SimpleKitModalBody,
  SimpleKitModalContent,
  SimpleKitModalDescription,
  SimpleKitModalFooter,
  SimpleKitModalHeader,
  SimpleKitModalTitle,
} from '@/components/simplekit/simplekit-modal';
import { Button } from '@/components/ui/button';
import { Check, ChevronLeft, Copy, RotateCcw } from 'lucide-react';
import { formatEther } from 'viem';
import {
  type Connector,
  useAccount,
  useBalance,
  useConnect,
  useDisconnect,
  useEnsAvatar,
  useEnsName,
  useSignMessage,
} from 'wagmi';

const MODAL_CLOSE_DURATION = 320;
const SimpleKitContext = React.createContext<{
  pendingConnector: Connector | null;
  setPendingConnector: React.Dispatch<React.SetStateAction<Connector | null>>;
  isConnectorError: boolean;
  setIsConnectorError: React.Dispatch<React.SetStateAction<boolean>>;
  open: boolean;
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  isAuthenticated: boolean;
  isAuthenticating: boolean;
  authError: string | null;
  showAuthPrompt: boolean;
  login: () => Promise<void>;
  logout: () => void;
  skipAuth: () => void;
  connectAndAuth: () => void; // Nouvelle fonction
  getAuthToken: () => string | null;
}>({
  pendingConnector: null,
  setPendingConnector: () => null,
  isConnectorError: false,
  setIsConnectorError: () => false,
  open: false,
  setOpen: () => false,
  isAuthenticated: false,
  isAuthenticating: false,
  authError: null,
  showAuthPrompt: false,
  login: async () => {},
  logout: () => {},
  skipAuth: () => {},
  connectAndAuth: () => {},
  getAuthToken: () => null,
});

function SimpleKitProvider(props: { children: React.ReactNode }) {
  const { status, address } = useAccount();
  const { signMessageAsync } = useSignMessage();
  const [pendingConnector, setPendingConnector] =
    React.useState<Connector | null>(null);
  const [isConnectorError, setIsConnectorError] = React.useState(false);
  const [open, setOpen] = React.useState(false);

  // États pour l'authentification
  const [isAuthenticating, setIsAuthenticating] = React.useState(false);
  const [authError, setAuthError] = React.useState<string | null>(null);
  const [showAuthPrompt, setShowAuthPrompt] = React.useState(false);
  const [autoAuthAfterConnection, setAutoAuthAfterConnection] =
    React.useState(false); // Nouveau state
  const [isAuthenticated, setIsAuthenticated] = React.useState(() => {
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

  const isConnected = address && !pendingConnector;

  // Fonctions d'authentification (inchangées)
  const getNonce = async (walletAddress: string) => {
    const response = await fetch(
      `http://localhost:8080/api/v1/users/nonce/${walletAddress}`
    );
    if (!response.ok) {
      throw new Error('Failed to get nonce');
    }
    return response.json();
  };

  const authenticate = async (
    walletAddress: string,
    signature: string,
    message: string
  ) => {
    const response = await fetch('http://localhost:8080/api/v1/users/auth', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        address: walletAddress,
        signature: signature,
        message: message,
      }),
    });

    if (!response.ok) {
      throw new Error('Authentication failed');
    }
    return response.json();
  };

  const login = async () => {
    if (!address || !isConnected) {
      setAuthError('Wallet not connected');
      return;
    }

    console.log('Authenticating with address:', address);

    setIsAuthenticating(true);
    setAuthError(null);

    try {
      const { message } = await getNonce(address);
      console.log('Nonce received, requesting signature...');

      const signature = await signMessageAsync({ message });
      console.log('Signature received, authenticating...');

      const authResponse = await authenticate(address, signature, message);
      console.log('Authentication successful');

      localStorage.setItem('auth_token', authResponse.token);
      localStorage.setItem('auth_expires_at', authResponse.expires_at);
      localStorage.setItem('wallet_address', address);

      setIsAuthenticated(true);
      setShowAuthPrompt(false);
      setAutoAuthAfterConnection(false);

      // Fermer le modal après un petit délai pour montrer le succès
      setTimeout(() => {
        setOpen(false);
      }, 500);
    } catch (error) {
      console.error('Authentication error:', error);
      setAuthError(
        error instanceof Error ? error.message : 'Authentication failed'
      );

      // En cas d'erreur en mode autoAuth, on peut proposer le prompt manuel
      if (autoAuthAfterConnection) {
        setAutoAuthAfterConnection(false);
        setShowAuthPrompt(true);
      }
    } finally {
      setIsAuthenticating(false);
    }
  };

  // Fonction pour initier la connexion avec authentification automatique
  const connectAndAuth = () => {
    setAutoAuthAfterConnection(true);
    setOpen(true);
  };

  const skipAuth = () => {
    setShowAuthPrompt(false);
    setAuthError(null);
    setAutoAuthAfterConnection(false);
  };

  const logout = () => {
    localStorage.removeItem('auth_token');
    localStorage.removeItem('auth_expires_at');
    localStorage.removeItem('wallet_address');
    setIsAuthenticated(false);
    setShowAuthPrompt(false);
    setAutoAuthAfterConnection(false);
  };

  const getAuthToken = (): string | null => {
    if (typeof window !== 'undefined') {
      return localStorage.getItem('auth_token');
    }
    return null;
  };

  // Vérifier l'expiration du token
  React.useEffect(() => {
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

    const interval = setInterval(checkTokenExpiry, 60000);
    return () => clearInterval(interval);
  }, []);

  // Se déconnecter automatiquement si le wallet est déconnecté
  React.useEffect(() => {
    if (!isConnected && isAuthenticated) {
      logout();
    }
  }, [isConnected, isAuthenticated]);

  // Authentification automatique après connexion si demandée
  React.useEffect(() => {
    if (
      status === 'connected' &&
      address &&
      !isAuthenticated &&
      autoAuthAfterConnection &&
      !isAuthenticating // Éviter les appels multiples
    ) {
      // Ajouter un délai pour permettre à la connexion de se stabiliser
      const timeout = setTimeout(() => {
        console.log('Starting automatic authentication...');
        login();
      }, 1000); // 1 seconde de délai

      return () => clearTimeout(timeout);
    }
  }, [
    status,
    address,
    isAuthenticated,
    autoAuthAfterConnection,
    isAuthenticating,
  ]);

  // Afficher le prompt d'authentification pour les connexions normales
  React.useEffect(() => {
    if (
      status === 'connected' &&
      address &&
      !isAuthenticated &&
      !autoAuthAfterConnection
    ) {
      const hasSkippedAuth = sessionStorage.getItem(`auth_skipped_${address}`);

      if (!hasSkippedAuth) {
        setShowAuthPrompt(true);
        setOpen(true);
      }
    }
  }, [status, address, isAuthenticated, autoAuthAfterConnection]);

  // Modifier ce useEffect pour mieux gérer la fermeture :
  React.useEffect(() => {
    if (status === 'connected' && pendingConnector) {
      // Ne fermer le modal que si on n'est PAS en mode autoAuth
      if (!autoAuthAfterConnection && !showAuthPrompt) {
        setOpen(false);
      }

      const timeout = setTimeout(() => {
        setPendingConnector(null);
        setIsConnectorError(false);
      }, MODAL_CLOSE_DURATION);

      return () => clearTimeout(timeout);
    }
  }, [status, pendingConnector, autoAuthAfterConnection, showAuthPrompt]);

  return (
    <SimpleKitContext.Provider
      value={{
        pendingConnector,
        setPendingConnector,
        isConnectorError,
        setIsConnectorError,
        open,
        setOpen,
        isAuthenticated,
        isAuthenticating,
        authError,
        showAuthPrompt,
        login,
        logout,
        skipAuth,
        connectAndAuth, // Nouvelle fonction
        getAuthToken,
      }}
    >
      {props.children}
      <SimpleKitModal open={open} onOpenChange={setOpen}>
        <SimpleKitModalContent>
          {isConnected ? <Account /> : <Connectors />}
        </SimpleKitModalContent>
      </SimpleKitModal>
    </SimpleKitContext.Provider>
  );
}

function ConnectWalletButton({ autoAuth = false }: { autoAuth?: boolean }) {
  const simplekit = useSimpleKit();
  const { address } = useAccount();
  const { data: ensName } = useEnsName({ address });
  const { data: ensAvatar } = useEnsAvatar({ name: ensName! });

  const handleClick = () => {
    if (autoAuth && !simplekit.isConnected) {
      simplekit.connectAndAuth(); // Connexion avec authentification automatique
    } else {
      simplekit.toggleModal(); // Comportement normal
    }
  };

  return (
    <Button onClick={handleClick} className="rounded-xl">
      {simplekit.isConnected ? (
        <>
          {ensAvatar && <img src={ensAvatar} alt="ENS Avatar" />}
          {address && (
            <span>
              {ensName ? `${ensName}` : simplekit.formattedAddress}
              {simplekit.isAuthenticated && (
                <span className="ml-1 text-green-500">✓</span>
              )}
            </span>
          )}
        </>
      ) : autoAuth ? (
        'Connect & Authenticate'
      ) : (
        'Connect Wallet'
      )}
    </Button>
  );
}

function Account() {
  const { address } = useAccount();
  const { disconnect } = useDisconnect();
  const { data: ensName } = useEnsName({ address });
  const { data: userBalance } = useBalance({ address });
  const context = React.useContext(SimpleKitContext);

  const formattedAddress = address?.slice(0, 6) + '•••' + address?.slice(-4);
  const formattedUserBalace = userBalance?.value
    ? parseFloat(formatEther(userBalance.value)).toFixed(4)
    : undefined;

  function handleDisconnect() {
    context.logout();
    context.setOpen(false);
    setTimeout(() => {
      disconnect();
    }, MODAL_CLOSE_DURATION);
  }

  function handleSkipAuth() {
    // Marquer que l'utilisateur a refusé l'auth pour cette session
    if (address) {
      sessionStorage.setItem(`auth_skipped_${address}`, 'true');
    }
    context.skipAuth();
  }

  // Si on doit montrer le prompt d'authentification
  if (context.showAuthPrompt && !context.isAuthenticated) {
    return (
      <>
        <SimpleKitModalHeader>
          <SimpleKitModalTitle>Authenticate Your Wallet</SimpleKitModalTitle>
          <SimpleKitModalDescription className="sr-only">
            Sign a message to authenticate with the platform.
          </SimpleKitModalDescription>
        </SimpleKitModalHeader>
        <SimpleKitModalBody className="h-[320px]">
          <div className="flex w-full flex-col items-center justify-center gap-6 md:pt-5">
            <div className="size-24 flex items-center justify-center">
              <img
                className="rounded-full"
                src={`https://avatar.vercel.sh/${address}?size=150`}
                alt="User gradient avatar"
              />
            </div>

            <div className="space-y-1 px-3.5 text-center sm:px-0">
              <h1 className="text-xl font-semibold">Welcome!</h1>
              <p className="text-balance text-sm text-muted-foreground">
                Would you like to sign a message to authenticate and access all
                platform features?
              </p>
            </div>

            <div className="w-full space-y-3">
              <Button
                className="w-full rounded-xl"
                onClick={context.login}
                disabled={context.isAuthenticating}
              >
                {context.isAuthenticating
                  ? 'Authenticating...'
                  : 'Sign Message to Authenticate'}
              </Button>

              <Button
                className="w-full rounded-xl"
                variant="outline"
                onClick={handleSkipAuth}
                disabled={context.isAuthenticating}
              >
                Skip for Now
              </Button>

              {context.authError && (
                <p className="text-xs text-red-500 text-center">
                  {context.authError}
                </p>
              )}
            </div>
          </div>
        </SimpleKitModalBody>
      </>
    );
  }

  // Interface normale après connexion/authentification
  return (
    <>
      <SimpleKitModalHeader>
        <SimpleKitModalTitle>
          {context.isAuthenticated ? 'Authenticated' : 'Connected'}
        </SimpleKitModalTitle>
        <SimpleKitModalDescription className="sr-only">
          Account modal for your connected Web3 wallet.
        </SimpleKitModalDescription>
      </SimpleKitModalHeader>
      <SimpleKitModalBody className="h-[320px]">
        <div className="flex w-full flex-col items-center justify-center gap-6 md:pt-5">
          <div className="size-24 flex items-center justify-center">
            <img
              className="rounded-full"
              src={`https://avatar.vercel.sh/${address}?size=150`}
              alt="User gradient avatar"
            />
          </div>

          <div className="space-y-1 px-3.5 text-center sm:px-0">
            <div className="flex items-center gap-1.5">
              <h1 className="text-xl font-semibold">
                <div>{ensName ? `${ensName}` : formattedAddress}</div>
              </h1>
              <CopyAddressButton />
            </div>
            <p className="text-balance text-sm text-muted-foreground">
              {`${formattedUserBalace ?? '0.00'} ETH`}
            </p>
          </div>

          <div className="w-full space-y-3">
            {!context.isAuthenticated && (
              <Button
                className="w-full rounded-xl"
                onClick={context.login}
                disabled={context.isAuthenticating}
              >
                {context.isAuthenticating
                  ? 'Authenticating...'
                  : 'Sign Message to Login'}
              </Button>
            )}

            {context.isAuthenticated && (
              <div className="text-center">
                <p className="text-sm text-green-600 mb-2">✓ Authenticated</p>
              </div>
            )}

            {context.authError && (
              <p className="text-xs text-red-500 text-center">
                {context.authError}
              </p>
            )}

            <Button
              className="w-full rounded-xl"
              variant="outline"
              onClick={handleDisconnect}
            >
              Disconnect
            </Button>
          </div>
        </div>
      </SimpleKitModalBody>
    </>
  );
}

function Connectors() {
  const context = React.useContext(SimpleKitContext);

  return (
    <>
      <SimpleKitModalHeader>
        <BackChevron />
        <SimpleKitModalTitle>
          {context.pendingConnector?.name ?? 'Connect Wallet'}
        </SimpleKitModalTitle>
        <SimpleKitModalDescription className="sr-only">
          Connect your Web3 wallet or create a new one.
        </SimpleKitModalDescription>
      </SimpleKitModalHeader>
      <SimpleKitModalBody>
        {context.pendingConnector ? <WalletConnecting /> : <WalletOptions />}
      </SimpleKitModalBody>
      <SimpleKitModalFooter>
        <div className="h-0" />
      </SimpleKitModalFooter>
    </>
  );
}

function WalletConnecting() {
  const context = React.useContext(SimpleKitContext);
  const { status } = useAccount();

  // Déterminer le titre et la description en fonction de l'état
  const getTitle = () => {
    if (context.isConnectorError) return 'Request Error';
    if (status === 'connected' && context.isAuthenticating)
      return 'Authenticating...';
    if (status === 'connected') return 'Connected';
    return 'Requesting Connection';
  };

  const getDescription = () => {
    if (context.isConnectorError) {
      return 'There was an error with the request. Click above to try again.';
    }
    if (status === 'connected' && context.isAuthenticating) {
      return 'Please sign the message in your wallet to authenticate.';
    }
    if (status === 'connected') {
      return 'Wallet connected successfully. Preparing authentication...';
    }
    return `Open the ${context.pendingConnector?.name} browser extension to connect your wallet.`;
  };

  return (
    <div className="flex w-full flex-col items-center justify-center gap-9 md:pt-5">
      {context.pendingConnector?.icon && (
        <div className="size-[116px] relative flex items-center justify-center rounded-2xl border p-3">
          <img
            src={context.pendingConnector?.icon}
            alt={context.pendingConnector?.name}
            className="size-full overflow-hidden rounded-2xl"
          />

          {/* Spinner pendant l'authentification */}
          {(context.isAuthenticating ||
            (status === 'connected' && !context.isAuthenticated)) && (
            <div className="absolute inset-0 flex items-center justify-center bg-background/80 rounded-2xl">
              <div className="animate-spin h-6 w-6 border-2 border-primary border-t-transparent rounded-full" />
            </div>
          )}

          {context.isConnectorError ? <RetryConnectorButton /> : null}
        </div>
      )}

      <div className="space-y-3.5 px-3.5 text-center sm:px-0">
        <h1 className="text-xl font-semibold">{getTitle()}</h1>
        <p className="text-balance text-sm text-muted-foreground">
          {getDescription()}
        </p>

        {/* Afficher l'erreur d'authentification si elle existe */}
        {context.authError && (
          <p className="text-sm text-red-500 mt-2">{context.authError}</p>
        )}
      </div>
    </div>
  );
}

function WalletOptions() {
  const context = React.useContext(SimpleKitContext);
  const { connectors, connect } = useConnectors();

  return (
    <div className="flex flex-col gap-3.5">
      {connectors.map((connector) => (
        <WalletOption
          key={connector.uid}
          connector={connector}
          onClick={() => {
            context.setIsConnectorError(false);
            context.setPendingConnector(connector);
            connect({ connector });
          }}
        />
      ))}
    </div>
  );
}

function WalletOption(props: { connector: Connector; onClick: () => void }) {
  const [ready, setReady] = React.useState(false);

  React.useEffect(() => {
    async function checkReady() {
      const provider = await props.connector.getProvider();
      setReady(!!provider);
    }
    checkReady()
      .then(() => null)
      .catch(() => null);
  }, [props.connector]);

  return (
    <Button
      disabled={!ready}
      onClick={props.onClick}
      size="lg"
      variant="secondary"
      className="justify-between rounded-xl px-4 py-7 text-base font-semibold"
    >
      <p>{props.connector.name}</p>
      {props.connector.icon && (
        <img
          src={props.connector.icon}
          alt={props.connector.name}
          className="size-8 overflow-hidden rounded-[6px]"
        />
      )}
    </Button>
  );
}

function CopyAddressButton() {
  const { address } = useAccount();
  const [copied, setCopied] = React.useState(false);

  React.useEffect(() => {
    const timeout = setTimeout(() => {
      if (copied) setCopied(false);
    }, 1000);
    return () => clearTimeout(timeout);
  }, [copied, setCopied]);

  async function handleCopy() {
    setCopied(true);
    await navigator.clipboard.writeText(address!);
  }

  return (
    <button className="text-muted-foreground" onClick={handleCopy}>
      {copied ? (
        <Check className="size-4" strokeWidth={4} />
      ) : (
        <Copy className="size-4" strokeWidth={4} />
      )}
    </button>
  );
}

function BackChevron() {
  const context = React.useContext(SimpleKitContext);

  if (!context.pendingConnector) {
    return null;
  }

  function handleClick() {
    context.setIsConnectorError(false);
    context.setPendingConnector(null);
  }

  return (
    <button
      className="absolute left-[26px] top-[42px] z-50 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none data-[state=open]:bg-accent data-[state=open]:text-muted-foreground md:top-[26px]"
      onClick={handleClick}
    >
      <ChevronLeft className="h-4 w-4" />
      <span className="sr-only">Cancel connection</span>
    </button>
  );
}

function RetryConnectorButton() {
  const context = React.useContext(SimpleKitContext);
  const { connect } = useConnect({
    mutation: {
      onError: () => context.setIsConnectorError(true),
    },
  });

  function handleClick() {
    if (context.pendingConnector) {
      context.setIsConnectorError(false);
      connect({ connector: context.pendingConnector });
    }
  }

  return (
    <Button
      size="icon"
      variant="secondary"
      className="group absolute -bottom-2 -right-2 rounded-full bg-muted p-1.5 shadow"
      onClick={handleClick}
    >
      <RotateCcw className="size-4 transition-transform group-hover:-rotate-45" />
    </Button>
  );
}

function useConnectors() {
  const context = React.useContext(SimpleKitContext);
  const { connect, connectors } = useConnect({
    mutation: {
      onError: () => context.setIsConnectorError(true),
    },
  });

  const sortedConnectors = React.useMemo(() => {
    let metaMaskConnector: Connector | undefined;
    let injectedConnector: Connector | undefined;

    const formattedConnectors = connectors.reduce(
      (acc: Array<Connector>, curr) => {
        console.log(curr.id);
        switch (curr.id) {
          case 'metaMaskSDK':
            metaMaskConnector = {
              ...curr,
              icon: 'https://utfs.io/f/be0bd88f-ce87-4cbc-b2e5-c578fa866173-sq4a0b.png',
            };
            return acc;
          case 'metaMask':
            injectedConnector = {
              ...curr,
              icon: 'https://utfs.io/f/be0bd88f-ce87-4cbc-b2e5-c578fa866173-sq4a0b.png',
            };
            return acc;
          case 'safe':
            acc.push({
              ...curr,
              icon: 'https://utfs.io/f/164ea200-3e15-4a9b-9ce5-a397894c442a-awpd29.png',
            });
            return acc;
          case 'coinbaseWalletSDK':
            acc.push({
              ...curr,
              icon: 'https://utfs.io/f/53e47f86-5f12-404f-a98b-19dc7b760333-chngxw.png',
            });
            return acc;
          case 'walletConnect':
            acc.push({
              ...curr,
              icon: 'https://utfs.io/f/5bfaa4d1-b872-48a7-9d37-c2517d4fc07a-utlf4g.png',
            });
            return acc;
          default:
            acc.unshift(curr);
            return acc;
        }
      },
      []
    );

    if (
      metaMaskConnector &&
      !formattedConnectors.find(
        ({ id }) =>
          id === 'io.metamask' ||
          id === 'io.metamask.mobile' ||
          id === 'injected'
      )
    ) {
      return [metaMaskConnector, ...formattedConnectors];
    }

    if (injectedConnector) {
      const nonMetaMaskConnectors = formattedConnectors.filter(
        ({ id }) => id !== 'io.metamask' && id !== 'io.metamask.mobile'
      );
      return [injectedConnector, ...nonMetaMaskConnectors];
    }
    return formattedConnectors;
  }, [connectors]);

  return { connectors: sortedConnectors, connect };
}

/*
 * This hook can be moved to a separate file
 * if desired (src/hooks/use-simple-kit.tsx).
 */
function useSimpleKit() {
  const { address } = useAccount();
  const context = React.useContext(SimpleKitContext);

  const isModalOpen = context.open;
  const isConnected = address && !context.pendingConnector;
  const formattedAddress = address?.slice(0, 6) + '•••' + address?.slice(-4);

  function open() {
    context.setOpen(true);
  }

  function close() {
    context.setOpen(false);
  }

  function toggleModal() {
    context.setOpen((prevState) => !prevState);
  }

  return {
    isModalOpen,
    isConnected,
    isAuthenticated: context.isAuthenticated,
    isAuthenticating: context.isAuthenticating,
    authError: context.authError,
    showAuthPrompt: context.showAuthPrompt,
    formattedAddress,
    open,
    close,
    toggleModal,
    login: context.login,
    logout: context.logout,
    skipAuth: context.skipAuth,
    getAuthToken: context.getAuthToken,
    connectAndAuth: context.connectAndAuth, // Add this line
  };
}

export {
  ConnectWalletButton,
  SimpleKitContext,
  SimpleKitProvider,
  useSimpleKit,
};
