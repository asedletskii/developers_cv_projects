import { ROUTES } from '@shared/constants'
import { generatePath } from 'react-router'

/** Публичный адрес страницы `/cv/<id>`: `id` — идентификатор резюме или собственный UUID ссылки. */
export const getResumeLinkUrl = (id: string, origin = window.location.origin) =>
  `${origin}${generatePath(ROUTES.cvResume, { cvId: id })}`
