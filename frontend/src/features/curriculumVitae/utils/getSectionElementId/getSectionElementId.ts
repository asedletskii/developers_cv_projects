/** Id заголовка секции в разметке — на него ссылается `aria-labelledby` тела секции. */
export const getTitleId = (id: string) => `${id}-title`

/** Id тела секции в разметке — на него ссылается `aria-controls` заголовка. */
export const getContentId = (id: string) => `${id}-content`
