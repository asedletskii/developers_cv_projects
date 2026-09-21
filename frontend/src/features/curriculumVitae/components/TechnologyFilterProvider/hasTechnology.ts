import { Technology } from '@features/curriculumVitae/types'

const normalize = (name: string) => name.trim().toLowerCase()

/** Есть ли технология в стеке (без учёта регистра и пробелов по краям). */
export const hasTechnology = (stack: Technology[], name: string) =>
  stack.some((technology) => normalize(technology.name) === normalize(name))

export const isSameTechnology = (a: string, b: string) => normalize(a) === normalize(b)
