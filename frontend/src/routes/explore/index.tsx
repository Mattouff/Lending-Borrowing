import { createFileRoute, Navigate } from '@tanstack/react-router';

export const Route = createFileRoute('/explore/')({
  component: () => <Navigate to="/explore/all" />,
});
