/** Запасной экран на случай, когда сломались сами провайдеры: тема и переводы недоступны,
 *  поэтому только чистая разметка и два языка вручную */
export const RootFallback = () => {
  const isRu = navigator.language.toLowerCase().startsWith('ru')

  return (
    <div
      role="alert"
      style={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        gap: 16,
        padding: 24,
        textAlign: 'center',
        background: '#151517',
        color: '#f4f4f6',
        fontFamily: 'system-ui, -apple-system, sans-serif',
      }}
    >
      <h1 style={{ margin: 0, fontSize: 28 }}>
        {isRu ? 'Что-то пошло не так' : 'Something went wrong'}
      </h1>
      <p style={{ margin: 0, color: '#a1a1a8' }}>
        {isRu
          ? 'Обновите страницу — обычно этого достаточно.'
          : 'Reload the page — that usually helps.'}
      </p>
      <button
        onClick={() => window.location.reload()}
        style={{
          padding: '12px 28px',
          border: 'none',
          borderRadius: 999,
          background: '#ffcc3c',
          color: '#1f1f1d',
          fontSize: 16,
          fontWeight: 600,
          cursor: 'pointer',
        }}
      >
        {isRu ? 'Обновить страницу' : 'Reload'}
      </button>
    </div>
  )
}
