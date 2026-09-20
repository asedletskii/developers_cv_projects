import { useEffect } from 'react'

/** Закрывает страницу от индексации, пока она на экране: личные ссылки, кабинет, служебные экраны */
export const useNoIndex = () => {
  useEffect(() => {
    const meta = document.createElement('meta')
    meta.name = 'robots'
    meta.content = 'noindex, nofollow'
    document.head.appendChild(meta)

    return () => meta.remove()
  }, [])
}
