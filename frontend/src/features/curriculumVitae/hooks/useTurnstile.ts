import { useCallback, useEffect, useRef, useState } from 'react'
import {
  TURNSTILE_ACTION,
  TURNSTILE_SCRIPT_SRC,
  TURNSTILE_TOKEN_TIMEOUT_MS,
} from '@features/curriculumVitae/constants'

const TURNSTILE_SITE_KEY = import.meta.env.VITE_TURNSTILE_SITE_KEY

type TurnstileApi = {
  render: (container: HTMLElement, options: Record<string, unknown>) => string
  execute: (widgetId: string) => void
  reset: (widgetId: string) => void
  remove: (widgetId: string) => void
}

declare global {
  interface Window {
    turnstile?: TurnstileApi
  }
}

let scriptPromise: Promise<void> | undefined

const loadTurnstile = () => {
  scriptPromise ??= new Promise<void>((resolve, reject) => {
    if (window.turnstile) return resolve()

    const script = document.createElement('script')
    script.src = TURNSTILE_SCRIPT_SRC
    script.async = true
    script.onload = () => resolve()
    script.onerror = () => {
      scriptPromise = undefined
      reject(new Error('Turnstile script failed to load'))
    }
    document.head.appendChild(script)
  })

  return scriptPromise
}

type Pending = { resolve: (token: string) => void; reject: (error: Error) => void }

export const useTurnstile = () => {
  // Контейнер лежит в состоянии, а не в ref: он может смонтироваться позже хука (например, внутри Dialog)
  const [container, setContainer] = useState<HTMLElement | null>(null)
  const [isChallengeVisible, setIsChallengeVisible] = useState(false)
  const widgetIdRef = useRef<string | undefined>(undefined)
  const pendingRef = useRef<Pending | undefined>(undefined)

  useEffect(() => {
    if (!TURNSTILE_SITE_KEY || !container) return

    let isCancelled = false

    const fail = (message: string) => {
      setIsChallengeVisible(false)
      pendingRef.current?.reject(new Error(message))
    }

    loadTurnstile()
      .then(() => {
        if (isCancelled || !window.turnstile) return

        widgetIdRef.current = window.turnstile.render(container, {
          sitekey: TURNSTILE_SITE_KEY,
          action: TURNSTILE_ACTION,
          execution: 'execute',
          appearance: 'interaction-only',
          theme: 'auto',
          'before-interactive-callback': () => setIsChallengeVisible(true),
          'after-interactive-callback': () => setIsChallengeVisible(false),
          callback: (token: string) => {
            setIsChallengeVisible(false)
            pendingRef.current?.resolve(token)
          },
          'error-callback': () => fail('Turnstile error'),
          'timeout-callback': () => fail('Turnstile timeout'),
        })
      })
      .catch(() => undefined)

    return () => {
      isCancelled = true
      if (widgetIdRef.current) window.turnstile?.remove(widgetIdRef.current)
      widgetIdRef.current = undefined
    }
  }, [container])

  const getToken = useCallback(async (): Promise<string | undefined> => {
    if (!TURNSTILE_SITE_KEY) return undefined

    await loadTurnstile()

    const widgetId = widgetIdRef.current
    if (!widgetId || !window.turnstile) throw new Error('Turnstile is not ready')

    const turnstile = window.turnstile

    return new Promise<string>((resolve, reject) => {
      const timer = window.setTimeout(
        () => reject(new Error('Turnstile timeout')),
        TURNSTILE_TOKEN_TIMEOUT_MS,
      )
      pendingRef.current = {
        resolve: (token) => (window.clearTimeout(timer), resolve(token)),
        reject: (error) => (window.clearTimeout(timer), reject(error)),
      }

      // execute() второй раз без reset() падает, а токен всё равно одноразовый
      turnstile.reset(widgetId)
      turnstile.execute(widgetId)
    }).finally(() => {
      pendingRef.current = undefined
    })
  }, [])

  const cancel = useCallback(() => {
    setIsChallengeVisible(false)
    pendingRef.current?.reject(new Error('Turnstile cancelled'))
  }, [])

  return { setContainer, getToken, isChallengeVisible, cancel }
}
