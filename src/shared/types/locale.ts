/** Язык интерфейса */
export enum Languages {
  Ru = 'ru',
  En = 'en',
}

export type OneLocaleDictionary = Record<string, any>

export type LocaleDictionary = Record<Languages, OneLocaleDictionary>

export type AppLocaleDictionary = (typeof import('@app/locale').appLocaleDict)['ru']

export type Primitive = string | number | boolean | null | undefined

export type NestedKeyOf<T> = {
  [K in keyof T & string]: T[K] extends Primitive
    ? K
    : T[K] extends object
      ? K | `${K}.${NestedKeyOf<T[K]>}`
      : K
}[keyof T & string]

export type PathValue<T, P extends string> = P extends `${infer Key}.${infer Rest}`
  ? Key extends keyof T
    ? PathValue<T[Key], Rest>
    : never
  : P extends keyof T
    ? T[P]
    : never

export type UnionToIntersection<U> = (U extends any ? (k: U) => void : never) extends (
  k: infer I,
) => void
  ? I
  : never
