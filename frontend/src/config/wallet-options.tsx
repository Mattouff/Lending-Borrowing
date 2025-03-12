import { useConnect } from 'wagmi';

export function WalletOptions() {
  const { connectors, connect } = useConnect();

  return (
    <div>
      <button onClick={() => connect({ connector: connectors[0] })}>
        Connect MetaMask
      </button>
    </div>
  );
}
