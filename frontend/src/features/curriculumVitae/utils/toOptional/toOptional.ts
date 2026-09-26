/** Пустая строка из формы — «не задано»: бэк принимает необязательные поля указателями. */
export const toOptionalString = (value?: string) => value || undefined

/** Не число (пустое поле формы даёт NaN) — «не задано». */
export const toOptionalNumber = (value?: number) => (Number.isFinite(value) ? value : undefined)
