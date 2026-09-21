import { ReactNode } from 'react'
import { BrowserRouter } from 'react-router'

interface Props {
  children: ReactNode
}

export const RoutingProvider = ({ children }: Props) => {
  return <BrowserRouter>{children}</BrowserRouter>
}
