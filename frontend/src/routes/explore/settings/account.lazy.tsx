import { createLazyFileRoute } from '@tanstack/react-router'
import SettingsAccount from '@/features/settings/account'

export const Route = createLazyFileRoute('/explore/settings/account')({
  component: SettingsAccount,
})
