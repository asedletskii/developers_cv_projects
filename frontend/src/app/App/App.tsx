import { Suspense } from 'react'
import { PageLoader } from '../PageLoader'
import { PageRouter } from '../PageRouter'
import { Providers } from '../providers'
import { RootLayout } from '../RootLayout'

export const App = () => {
  return (
    <Providers>
      <RootLayout>
        <Suspense fallback={<PageLoader />}>
          <PageRouter />
        </Suspense>
      </RootLayout>
    </Providers>
  )
}
