import { PropsWithChildren } from 'react'
import { useLocation } from 'react-router'
import { ErrorBoundary } from './ErrorBoundary'
import { ErrorFallback } from './ErrorFallback'

/** Ошибка на одной странице не роняет шапку и подвал; при переходе на другую страницу
 *  граница пересоздаётся (`key`), и приложение оживает без перезагрузки */
export const RouteBoundary = ({ children }: PropsWithChildren) => {
  const { pathname } = useLocation()

  return (
    <ErrorBoundary key={pathname} fallback={(error) => <ErrorFallback error={error} />}>
      {children}
    </ErrorBoundary>
  )
}
