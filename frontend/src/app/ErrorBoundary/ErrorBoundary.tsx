import { Component, ErrorInfo, PropsWithChildren, ReactNode } from 'react'
import { reloadPage } from './reloadPage'

const RELOAD_KEY = 'chunk-reload-at'
const RELOAD_COOLDOWN_MS = 10_000

/** После выката старые файлы сборки исчезают: открытая вкладка просит несуществующий кусок кода */
const isChunkLoadError = (error: Error) =>
  /Failed to fetch dynamically imported module|Importing a module script failed|error loading dynamically imported module|ChunkLoadError/i.test(
    error.message,
  )

type Props = PropsWithChildren<{
  fallback: (error: Error) => ReactNode
}>

type State = { error: Error | null }

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null }

  static getDerivedStateFromError(error: Error): State {
    return { error }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    console.error(error, info.componentStack)

    if (!isChunkLoadError(error)) return

    // Один раз подтягиваем свежую версию сами; пауза не даёт зациклить перезагрузки
    const lastReload = Number(sessionStorage.getItem(RELOAD_KEY) ?? 0)
    if (Date.now() - lastReload > RELOAD_COOLDOWN_MS) {
      sessionStorage.setItem(RELOAD_KEY, String(Date.now()))
      reloadPage()
    }
  }

  render() {
    return this.state.error ? this.props.fallback(this.state.error) : this.props.children
  }
}
