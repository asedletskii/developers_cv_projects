import { Technology } from '@features/curriculumVitae/types'
import { isSameTechnology } from '../isSameTechnology'

/** Есть ли технология в стеке (без учёта регистра и пробелов по краям). */
export const hasTechnology = (stack: Technology[], name: string) =>
  stack.some((technology) => isSameTechnology(technology.name, name))
