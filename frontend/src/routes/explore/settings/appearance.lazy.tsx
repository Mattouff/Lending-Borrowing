import { createLazyFileRoute } from '@tanstack/react-router'
import SettingsAppearance from '@/features/settings/appearance'

export const Route = createLazyFileRoute('/explore/settings/appearance')(
  { component: SettingsAppearance }
)
