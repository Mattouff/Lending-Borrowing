import { createConfig, http, injected } from 'wagmi';
import { base, mainnet } from 'wagmi/chains';
import { metaMask, safe, walletConnect } from 'wagmi/connectors';

const projectId = 'cb9a3d203b35a979abb8e955f45cdab1';

export const config = createConfig({
  chains: [mainnet, base],
  connectors: [injected(), walletConnect({ projectId }), metaMask(), safe()],
  transports: {
    [mainnet.id]: http(),
    [base.id]: http(),
  },
});
