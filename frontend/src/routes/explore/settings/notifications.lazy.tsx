import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute(
  '/explore/settings/notifications',
)({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_authenticated/settings/notifications"!</div>
}
