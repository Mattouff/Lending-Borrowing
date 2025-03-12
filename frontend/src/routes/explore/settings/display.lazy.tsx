import { createLazyFileRoute } from '@tanstack/react-router'
import SettingsDisplay from '@/features/settings/display'

export const Route = createLazyFileRoute('/explore/settings/display')({
  component: SettingsDisplay,
})
