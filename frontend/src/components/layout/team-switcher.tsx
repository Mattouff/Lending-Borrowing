import {
  SidebarMenu,
  SidebarMenuItem,
  useSidebar,
} from '@/components/ui/sidebar';
import * as React from 'react';

export function TeamSwitcher({
  teams,
}: {
  teams: {
    name: string;
    logo: React.ElementType;
    plan: string;
  }[];
}) {
  const { isMobile } = useSidebar();
  const [activeTeam, setActiveTeam] = React.useState(teams[0]);

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <div className="flex items-center gap-3">
          <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
            <activeTeam.logo className="size-4" />
          </div>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-semibold pt-1">Dashboard</span>
            <span className="truncate text-xs">Interactive plan</span>
          </div>
        </div>
        {/* <ChevronsUpDown className="ml-auto" /> */}
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
