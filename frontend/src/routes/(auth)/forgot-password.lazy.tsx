import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute('/(auth)/forgot-password')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/(auth)/forgot-password"!</div>
}
