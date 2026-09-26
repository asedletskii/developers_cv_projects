import { useCallback, useState } from 'react'

export type ResumeLink = {
  /** Собственный UUID ссылки: адрес — `/cv/<id>`, а не резюме с меткой, чтобы адресат не мог её
   *  стереть, а бэк по нему узнавал, откуда пришли (счётчик просмотров) */
  id: string
  name: string
}

const storageKey = (resumeId: string) => `cv:links:${resumeId}`

const read = (resumeId: string): ResumeLink[] => {
  try {
    const parsed = JSON.parse(localStorage.getItem(storageKey(resumeId)) ?? '[]')
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

/**
 * Именованные ссылки на резюме. Бэка под них пока нет, поэтому лежат в localStorage браузера
 * (по резюме) — на другом устройстве их не будет. Когда появятся ручки, заменить только этот хук.
 */
export const useResumeLinks = (resumeId: string) => {
  const [links, setLinks] = useState(() => read(resumeId))

  const update = useCallback(
    (next: (links: ResumeLink[]) => ResumeLink[]) => {
      setLinks((current) => {
        const updated = next(current)

        try {
          localStorage.setItem(storageKey(resumeId), JSON.stringify(updated))
        } catch {
          // Хранилище недоступно (приватный режим) — ссылки живут до перезагрузки страницы
        }

        return updated
      })
    },
    [resumeId],
  )

  const create = (name: string) =>
    update((current) => [...current, { id: crypto.randomUUID(), name }])

  const rename = (id: string, name: string) =>
    update((current) => current.map((link) => (link.id === id ? { ...link, name } : link)))

  const remove = (id: string) => update((current) => current.filter((link) => link.id !== id))

  return { links, create, rename, remove }
}
