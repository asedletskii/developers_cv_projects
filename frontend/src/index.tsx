import '@app/dayjsSetup'
import { App } from '@app'
import { ErrorBoundary, RootFallback } from '@app/ErrorBoundary'
import { createRoot } from 'react-dom/client'

const container = document.getElementById('root')

if (container) {
  const root = createRoot(container)
  root.render(
    <ErrorBoundary fallback={() => <RootFallback />}>
      <App />
    </ErrorBoundary>,
  )
}
