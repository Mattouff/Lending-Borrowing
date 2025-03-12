import { createLazyFileRoute } from '@tanstack/react-router'
import Tasks from '@/features/tasks'

export const Route = createLazyFileRoute('/explore/tasks/')({
  component: Tasks,
})
