/** Доля высоты экрана, на которой лежит «линия чтения»: секция, пересекающая её, считается активной */
export const READING_LINE = 0.3

/** События пользователя, снимающие закрепление пункта навигации после клика */
export const UNLOCK_EVENTS = ['scrollend', 'wheel', 'touchstart', 'keydown'] as const
