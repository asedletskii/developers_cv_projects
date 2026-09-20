import { useEffect } from 'react'

/** Ставит заголовок вкладки браузера и возвращает прежний, когда компонент уходит с экрана. */
export const useDocumentTitle = (title?: string) => {
  useEffect(() => {
    if (!title) return

    const previousTitle = document.title
    document.title = title

    return () => {
      document.title = previousTitle
    }
  }, [title])
}
