import { useState } from 'react'
import { curriculumVitaeApi } from '@features/curriculumVitae/api'
import { ContactLink } from '@features/curriculumVitae/utils'
import { ToastType } from '@shared/types'
import { showToast } from '@shared/utils'

/**
 * Значение скрытого контакта бэк не отдаёт в публичном ответе — запрашиваем его по клику
 * (`POST /cv/:id/contact-views`). У открытых контактов `link.value` уже есть, `reveal` для них
 * ничего не делает.
 *
 * Исключение — предпросмотр в ЛК: там резюме читается владельческой ручкой, которая скрытые
 * контакты не вырезает, `link.value` приходит сразу. Кнопку всё равно показываем как «Показать»
 * (иначе предпросмотр выглядит не так, как настоящая публичная страница), но по клику просто
 * достаём уже пришедшее значение — без лишнего запроса на бэк.
 */
export const useContactButton = (resumeId: string, link: ContactLink) => {
  const [value, setValue] = useState(link.isHidden ? undefined : link.value)
  const [isRevealing, setIsRevealing] = useState(false)

  const reveal = async () => {
    if (value || isRevealing) return

    if (link.value) {
      setValue(link.value)
      return
    }

    setIsRevealing(true)
    try {
      setValue(await curriculumVitaeApi.revealContact(resumeId, link.key))
    } catch {
      showToast({ messageKey: 'revealFailed', type: ToastType.Error })
    } finally {
      setIsRevealing(false)
    }
  }

  const copy = async () => {
    if (!value) return

    try {
      await navigator.clipboard.writeText(value)
      showToast({ messageKey: 'copied', type: ToastType.Success })
    } catch {
      showToast({ messageKey: 'copyFailed', type: ToastType.Error })
    }
  }

  return { value, isRevealing, reveal, copy }
}
