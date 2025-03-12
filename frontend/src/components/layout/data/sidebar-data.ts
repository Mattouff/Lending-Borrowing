import { Command, HomeIcon, UserIcon } from 'lucide-react';
import { type SidebarData } from '../types';

export const sidebarData: SidebarData = {
  user: {
    name: 'satnaing',
    email: 'satnaingdev@gmail.com',
    avatar: '/avatars/shadcn.jpg',
  },
  teams: [
    {
      name: 'NFT Certs',
      logo: Command,
      plan: 'Avax + Certification',
    },
  ],
  navGroups: [
    {
      title: 'General',
      items: [
        {
          title: 'Explore all',
          url: '/explore/all',
          icon: HomeIcon,
        },
        {
          title: 'User',
          url: '/explore/user',
          icon: UserIcon,
        },
      ],
    },
  ],
};
