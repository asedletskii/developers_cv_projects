export const sharedDict = {
  ru: {
    Header: {
      logout: 'Выход',

      Tabs: {
        about: 'О проекте',
        lk: 'Личный кабинет',
      },
    },

    NotFound: {
      title: 'Страница потерялась',
      text: 'Такого адреса нет: возможно, ссылка устарела или в ней опечатка.',
      home: 'На главную',
      back: 'Назад',
    },

    Footer: {
      info: 'Pet-проект, который агрегирует все web-проекты двух разработчиков',
      developers: 'Разработчики',
      madeWith: 'Сервис для создания резюме',
      backToTop: 'Наверх',

      Developers: {
        sedletskii: 'Седлецкий Антон',
        kosterin: 'Костерин Дмитрий',

        Roles: {
          frontend: 'Frontend разработчик',
          backend: 'Backend разработчик',
        },
      },
    },

    ThemeToggle: {
      label: 'Сменить тему',
    },

    LanguageSwitch: {
      label: 'Язык интерфейса',
    },

    ContentContainer: {
      backButton: 'Назад',
    },

    Toasts: {
      copied: 'Скопировано',
      copyFailed: 'Не удалось скопировать',
      cannotGetPersonData: 'Не удалось загрузить резюме',
      cannotGetProject: 'Не удалось загрузить проект',
      revealFailed: 'Не удалось показать контакт',
      downloadFailed: 'Не удалось скачать резюме',
      downloadWithoutContacts: 'Контакты получить не удалось — файл собран без них',
      statusChangeFailed: 'Не удалось изменить статус резюме',
      linkActionFailed: 'Не удалось сохранить ссылку',
    },
  },
  en: {
    Header: {
      logout: 'Log out',

      Tabs: {
        about: 'About',
        lk: 'Personal account',
      },
    },

    NotFound: {
      title: 'Page not found',
      text: 'There is no such address: the link may be outdated or contain a typo.',
      home: 'Go home',
      back: 'Go back',
    },

    Footer: {
      info: 'A Pet project that aggregates all the web projects of two developers',
      developers: 'Developers',
      madeWith: 'From Russia with love',
      backToTop: 'Back to top',

      Developers: {
        sedletskii: 'Sedletskii Anton',
        kosterin: 'Kosterin Dmitry',

        Roles: {
          frontend: 'Frontend developer',
          backend: 'Backend developer',
        },
      },
    },

    ThemeToggle: {
      label: 'Switch theme',
    },

    LanguageSwitch: {
      label: 'Interface language',
    },

    ContentContainer: {
      backButton: '',
    },

    Toasts: {
      copied: 'Copied',
      copyFailed: "Couldn't copy",
      cannotGetPersonData: "Couldn't load the résumé",
      cannotGetProject: "Couldn't load the project",
      revealFailed: "Couldn't reveal the contact",
      downloadFailed: "Couldn't download the résumé",
      downloadWithoutContacts: "Couldn't get the contacts — the file was built without them",
      statusChangeFailed: "Couldn't change the résumé status",
      linkActionFailed: "Couldn't save the link",
    },
  },
} as const
