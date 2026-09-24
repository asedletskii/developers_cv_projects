// import { UndefinedScreen } from '@shared/screens/UndefinedScreen'
import { curriculumVitaeRoutes } from '@features/curriculumVitae/routes'
import { lkRoutes } from '@features/lk/routes'
import { Route, Routes, useLocation } from 'react-router'

export const PageRouter = () => {
  const location = useLocation()
  const state = location.state as { backgroundLocation?: Location }

  return (
    <Routes location={state?.backgroundLocation || location}>
      <Route path="/">
        {/* Общий префикс /cv на обе фичи: curriculumVitae — /cv/:cvId, lk — /cv/lk/:userId(/:cvId) */}
        <Route path="cv">
          {curriculumVitaeRoutes.pages}
          {lkRoutes.pages}
        </Route>
      </Route>

      {/* <Route path="*" element={<UndefinedScreen />} /> */}
    </Routes>
  )
}
