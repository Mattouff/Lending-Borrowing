import { AppSidebar } from '@/components/layout/app-sidebar';
import { TopHeader } from '@/components/layout/top-header';
import SkipToMain from '@/components/skip-to-main';
import { MARGIN_TOP, SidebarProvider } from '@/components/ui/sidebar';
import { cn } from '@/lib/utils';
import { createFileRoute, Outlet } from '@tanstack/react-router';
import Cookies from 'js-cookie';

export const Route = createFileRoute('/explore')({
  component: RouteComponent,
});

function RouteComponent() {
  const defaultOpen = Cookies.get('sidebar:state') !== 'false';

  return (
    <>
      <TopHeader />
      <SidebarProvider defaultOpen={defaultOpen}>
        <SkipToMain />
        <AppSidebar />
        <div
          id="content"
          className={cn(
            'ml-auto w-full max-w-full mt-3',
            'peer-data-[state=collapsed]:w-[calc(100%-var(--sidebar-width-icon)-1rem)]',
            'peer-data-[state=expanded]:w-[calc(100%-var(--sidebar-width))]',
            'transition-[width] duration-200 ease-linear',
            'flex h-svh flex-col',
            'group-data-[scroll-locked=1]/body:h-full',
            'group-data-[scroll-locked=1]/body:has-[main.fixed-main]:h-svh'
          )}
          style={{ paddingTop: MARGIN_TOP }} // Décale le contenu sous le header
        >
          <Outlet />
        </div>
      </SidebarProvider>
    </>
  );
}
