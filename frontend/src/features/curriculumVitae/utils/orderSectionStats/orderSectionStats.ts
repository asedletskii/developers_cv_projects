/** Сортирует строки статистики секций в порядке, в котором секции стоят на странице резюме сейчас.
 *  Порядок секций владелец меняет в редакторе, поэтому полагаться на порядок ответа сервера нельзя.
 *  Секции, которых на странице уже нет (скрыли или опустошили), уходят в конец — история не теряется. */
export const orderSectionStats = <T extends { section: string }>(
  items: T[],
  orderedIds: string[],
): T[] => {
  const position = (section: string) => {
    const index = orderedIds.indexOf(section)
    return index === -1 ? orderedIds.length : index
  }

  return [...items].sort((a, b) => position(a.section) - position(b.section))
}
