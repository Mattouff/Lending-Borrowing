import Dashboard from '@/features/dashboard';
import { createFileRoute } from '@tanstack/react-router';

export const Route = createFileRoute('/explore/all/')({
  component: Dashboard,
});
