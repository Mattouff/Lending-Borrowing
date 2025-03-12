'use client';

// import { wagmiContractConfig } from '@/abi/contract'
import { Header } from '@/components/layout/header';
import { Main } from '@/components/layout/main';
// import { Search } from '@/components/search'
// import { usePinataMetadata } from '@/hooks/pinata'
import { useAccount, useConnect, useDisconnect, useEnsName } from 'wagmi';

export default function Certifications() {
  const { address } = useAccount();
  const { disconnect } = useDisconnect();
  const { connect } = useConnect();
  const { data: ensName } = useEnsName({ address });

  return (
    <div>
      {/* <UsersProvider> */}
      <Header fixed></Header>
      <Main>
        <div className="mb-2 flex flex-wrap items-center justify-between space-y-2"></div>
      </Main>
      {/* <UsersDialogs /> */}
      {/* </UsersProvider> */}
    </div>
  );
}
