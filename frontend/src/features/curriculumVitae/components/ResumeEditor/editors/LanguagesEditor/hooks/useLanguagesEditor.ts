import { LanguageType, Person } from '@features/curriculumVitae/types'
import { useTranslate } from '@shared/hooks'
import { useFieldArray, useFormContext, useWatch } from 'react-hook-form'

/**
 * Логика вкладки «Языки». Язык в резюме уникален (повтор бэкенд отклоняет `409`), поэтому в списке
 * каждой строки не предлагаем языки, уже выбранные в других строках, а новая строка получает
 * первый ещё не занятый язык.
 */
export const useLanguagesEditor = () => {
  const translate = useTranslate('CurriculumVitae.About')
  const { control } = useFormContext<Person>()
  const { fields, append, remove, move } = useFieldArray({ control, name: 'languages' })
  const languages = useWatch({ control, name: 'languages' }) ?? []

  const usedLanguages = languages.map(({ language }) => language)
  const freeLanguage = Object.values(LanguageType).find((code) => !usedLanguages.includes(code))

  const getLanguageOptions = (index: number) =>
    Object.values(LanguageType)
      .filter((code) => code === languages[index]?.language || !usedLanguages.includes(code))
      .map((language) => ({ value: language, label: translate(`Languages.${language}`) }))

  return {
    fields,
    canAdd: !!freeLanguage,
    addLanguage: () => freeLanguage && append({ language: freeLanguage }),
    getLanguageOptions,
    canMoveUp: (index: number) => index > 0,
    canMoveDown: (index: number) => index < fields.length - 1,
    moveUp: (index: number) => move(index, index - 1),
    moveDown: (index: number) => move(index, index + 1),
    remove,
  }
}
