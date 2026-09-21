import { PropsWithChildren } from 'react'
import { ColorModeProvider } from '../ColorModeProvider'
import { LocalizationProvider } from '../LocalizationProvider'
import { ReduxProvider } from '../ReduxProvider'
import { RoutingProvider } from '../RoutingProvider'
import { ToastProvider } from '../ToastProvider'

export const Providers = ({ children }: PropsWithChildren) => {
  return (
    <ReduxProvider>
      <RoutingProvider>
        <ColorModeProvider>
          <LocalizationProvider>
            <ToastProvider>{children}</ToastProvider>
          </LocalizationProvider>
        </ColorModeProvider>
      </RoutingProvider>
    </ReduxProvider>
  )
}
