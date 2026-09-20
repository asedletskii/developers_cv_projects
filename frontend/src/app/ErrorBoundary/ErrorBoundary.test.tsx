import { act } from 'react'
import { createRoot } from 'react-dom/client'
import { ErrorBoundary } from './ErrorBoundary'
import { reloadPage } from './reloadPage'

jest.mock('./reloadPage', () => ({ reloadPage: jest.fn() }))

const reload = jest.mocked(reloadPage)

const Boom = ({ message }: { message: string }): never => {
  throw new Error(message)
}

const renderBoundary = (message: string) => {
  const container = document.createElement('div')
  const root = createRoot(container)

  act(() => {
    root.render(
      <ErrorBoundary fallback={(error) => <p>fallback: {error.message}</p>}>
        <Boom message={message} />
      </ErrorBoundary>,
    )
  })

  return container
}

describe('ErrorBoundary', () => {
  beforeEach(() => {
    ;(globalThis as { IS_REACT_ACT_ENVIRONMENT?: boolean }).IS_REACT_ACT_ENVIRONMENT = true
    jest.spyOn(console, 'error').mockImplementation(() => undefined)
    sessionStorage.clear()
    reload.mockClear()
  })

  afterEach(() => jest.restoreAllMocks())

  it('показывает запасной экран вместо белого', () => {
    const container = renderBoundary('сломалось')

    expect(container.textContent).toBe('fallback: сломалось')
    expect(reload).not.toHaveBeenCalled()
  })

  it('если не нашёлся кусок кода после выката — сам обновляет страницу, но не чаще раза в 10 секунд', () => {
    renderBoundary('Failed to fetch dynamically imported module: /assets/x.js')
    expect(reload).toHaveBeenCalledTimes(1)

    renderBoundary('Failed to fetch dynamically imported module: /assets/x.js')
    expect(reload).toHaveBeenCalledTimes(1)
  })
})
