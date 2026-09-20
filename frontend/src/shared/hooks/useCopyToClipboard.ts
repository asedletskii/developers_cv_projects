import { useCallback, useState } from 'react'

export const useCopyToClipboard = (resetMs = 2000) => {
  const [copiedKey, setCopiedKey] = useState<string | null>(null)

  const copy = useCallback(
    async (key: string, text: string) => {
      try {
        await navigator.clipboard.writeText(text)
        setCopiedKey(key)
        window.setTimeout(
          () => setCopiedKey((current) => (current === key ? null : current)),
          resetMs,
        )
      } catch {
        //TODO: буфер недоступен (не https/нет прав) — значение можно выделить и скопировать руками
      }
    },
    [resetMs],
  )

  return { copiedKey, copy }
}
