import { lazy } from 'react'
import { Navigate, Route } from 'react-router'

const LkProfileScreen = lazy(() => import('./screens/LkProfileScreen/LkProfileScreen'))
const LkResumeEditScreen = lazy(() => import('./screens/LkResumeEditScreen/LkResumeEditScreen'))
const LkResumeScreen = lazy(() => import('./screens/LkResumeScreen/LkResumeScreen'))

// Смонтировано под общим `<Route path="cv">` из PageRouter, рядом с curriculumVitaeRoutes —
// поэтому "lk" объявлен как статический сегмент, а не часть :cvId (react-router отдаёт
// приоритет статическим сегментам над динамическими при совпадении путей)
export const lkRoutes = {
  pages: (
    <Route path="lk">
      {/* Без userId показывать нечего — уходим на главную */}
      <Route index element={<Navigate to="/" replace />} />

      <Route path=":userId">
        <Route index element={<LkProfileScreen />} />
        {/* "new" — статический сегмент, приоритетнее :cvId при совпадении путей (react-router) */}
        <Route path="new" element={<LkResumeEditScreen />} />
        {/* Клик по резюме открывает редактирование; просмотр как есть — отдельная кнопка */}
        <Route path=":cvId" element={<LkResumeEditScreen />} />
        <Route path=":cvId/preview" element={<LkResumeScreen />} />
      </Route>
    </Route>
  ),
}
