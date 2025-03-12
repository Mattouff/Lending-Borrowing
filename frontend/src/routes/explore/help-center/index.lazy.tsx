import { createLazyFileRoute } from '@tanstack/react-router'
import ComingSoon from '@/components/coming-soon'

export const Route = createLazyFileRoute('/explore/help-center/')({
  component: ComingSoon,
})
