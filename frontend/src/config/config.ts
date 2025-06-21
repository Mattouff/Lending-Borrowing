import { createConfig, http, injected } from 'wagmi';
import { localhost, mainnet, polygon, sepolia } from 'wagmi/chains';
import { metaMask, walletConnect } from 'wagmi/connectors';

const projectId = 'cb9a3d203b35a979abb8e955f45cdab1';

export const config = createConfig({
  chains: [mainnet, sepolia, polygon, localhost],
  connectors: [
    metaMask(),
    injected({ target: 'metaMask' }),
    walletConnect({ projectId }),
  ],
  transports: {
    [sepolia.id]: http(),
    [mainnet.id]: http(),
    [polygon.id]: http(),
    [localhost.id]: http('http://127.0.0.1:8545'),
  },
  ssr: true,
});
