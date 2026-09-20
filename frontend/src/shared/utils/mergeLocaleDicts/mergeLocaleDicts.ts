import { LocaleDictionary, UnionToIntersection } from '@shared/types'

export const mergeLocaleDicts = <T extends LocaleDictionary[]>(
  ...dicts: T
): UnionToIntersection<T[number]> => {
  return dicts.reduce(
    (merged, dict) => ({
      ru: { ...merged.ru, ...dict.ru },
      en: { ...merged.en, ...dict.en },
    }),
    dicts[0],
  ) as UnionToIntersection<T[number]>
}
