import { Suspense } from 'react'
import { RouteBoundary } from '../ErrorBoundary'
import { PageLoader } from '../PageLoader'
import { PageRouter } from '../PageRouter'
import { Providers } from '../providers'
import { RootLayout } from '../RootLayout'

export const App = () => {
  return (
    <Providers>
      <RootLayout>
        <RouteBoundary>
          <Suspense fallback={<PageLoader />}>
            <PageRouter />
          </Suspense>
        </RouteBoundary>
      </RootLayout>
    </Providers>
  )
}
