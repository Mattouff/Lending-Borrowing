import { createFileRoute, Navigate } from '@tanstack/react-router';
// index.tsx à la racine du projet

export const Route = createFileRoute('/explore/')({
  component: () => <Navigate to="/explore/all" />,
});

// import { createFileRoute, Navigate } from '@tanstack/react-router';

// export const Route = createFileRoute('/explore/')({
//   component: () => <Navigate to="/explore/all" />,
// });
