export const DURATION = '0.4s'

/** Кривая как у листов iOS: быстрый старт, мягкая остановка */
export const EASE = 'cubic-bezier(0.32, 0.72, 0, 1)'

export const getTitleId = (id: string) => `${id}-title`

export const getContentId = (id: string) => `${id}-content`

/** Отступы вокруг содержимого: живут внутри сворачиваемой части, чтобы у свёрнутой секции остался компактный зазор */
export const contentSx = { pt: 3, pb: { xs: 4, md: 6 } }
