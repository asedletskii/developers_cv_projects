const normalize = (name: string) => name.trim().toLowerCase()

/** Одна ли это технология: без учёта регистра и пробелов по краям. */
export const isSameTechnology = (a: string, b: string) => normalize(a) === normalize(b)
