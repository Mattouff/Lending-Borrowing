import { createLazyFileRoute } from '@tanstack/react-router'
import Chats from '@/features/chats'

export const Route = createLazyFileRoute('/explore/chats/')({
  component: Chats,
})
