'use client';

import { Header } from '@/components/layout/header';
import { Main } from '@/components/layout/main';
import { getExploreContent } from '@/lib/links';
import { useLocation } from '@tanstack/react-router';
import React from 'react';
import Global from './pages/all';
import User from './pages/user';

// 📌 Import des composants Table de shadcn/ui

export default function Source() {
  const location = useLocation();
  const path = getExploreContent(location.pathname);

  const Content: React.FC = () => {
    switch (path[0]) {
      case '/all':
        return <Global />;
      case '/user':
        return <User />;
      default:
        return <Global />;
    }
  };

  return (
    <div>
      <Header fixed path={path[0]} />
      <Main>
        <Content />
      </Main>
    </div>
  );
}
