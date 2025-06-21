import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';
import React from 'react';

import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb';
import { LinksDefinition } from '@/lib/links';
import { IconCube } from '@tabler/icons-react';
import { Link } from '@tanstack/react-router';

interface HeaderProps extends React.HTMLAttributes<HTMLElement> {
  fixed?: boolean;
  path: string;
  ref?: React.Ref<HTMLElement>;
}

export const Header = ({
  className,
  fixed,
  path,
  children,
  ...props
}: HeaderProps) => {
  const [offset, setOffset] = React.useState(0);

  React.useEffect(() => {
    const onScroll = () => {
      setOffset(document.body.scrollTop || document.documentElement.scrollTop);
    };

    // Add scroll listener to the body
    document.addEventListener('scroll', onScroll, { passive: true });

    // Clean up the event listener on unmount
    return () => document.removeEventListener('scroll', onScroll);
  }, []);

  console.log(
    'Links',
    path as keyof typeof LinksDefinition,
    LinksDefinition[path as keyof typeof LinksDefinition]
  );

  return (
    <div
      className={cn(
        'flex h-16 items-center gap-3 p-4 sm:gap-4 px-6',
        offset > 10 && fixed ? '' : 'shadow-none',
        className
      )}
      {...props}
    >
      <IconCube className="scale-125 sm:scale-100 animate-pulse transition-800" />
      <Separator orientation="vertical" className="h-6" />
      <Breadcrumb>
        <BreadcrumbList>
          <BreadcrumbItem>
            <BreadcrumbLink>Magnum</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink>Dashboard</BreadcrumbLink>
          </BreadcrumbItem>
          <BreadcrumbSeparator />
          <BreadcrumbItem>
            <BreadcrumbLink>
              <Link to={'/explore/' + path}>
                {
                  LinksDefinition[
                    path.replace(/^\//, '') as keyof typeof LinksDefinition
                  ]
                }
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
        </BreadcrumbList>
      </Breadcrumb>
      {children}
    </div>
  );
};

Header.displayName = 'Header';
