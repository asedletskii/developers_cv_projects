export const curriculumVitaeDict = {
  ru: {
    CurriculumVitae: {
      Error: {
        notFoundTitle: 'Резюме не найдено',
        notFoundText: 'Возможно, ссылка устарела или в ней опечатка.',
        failedTitle: 'Не удалось загрузить резюме',
        failedText: 'Попробуйте обновить страницу чуть позже.',
        errorCode: 'Код ошибки: {{status}}',
        requestId: 'ID запроса: {{requestId}}',
        home: 'На главную',
        retry: 'Попробовать снова',
        copyDetails: 'Скопировать детали',
        copied: 'Скопировано',
      },

      Availability: {
        open: 'Открыт к предложениям',
        considering: 'Рассматриваю предложения',
        busy: 'Сейчас не ищу',
        from: 'с {{date}}',
        updated: 'Обновлено {{time}}',
      },

      Filter: {
        technology: 'Технология',
        experience: 'Опыт',
        projects: 'Проекты',
        reset: 'Сбросить фильтр',
        select: 'Показать, где используется {{name}}',
      },

      Controls: {
        download: 'Скачать резюме',
      },

      // Подписи в файле, который собирается по кнопке «Скачать» (язык файла — язык резюме)
      HtmlExport: {
        technologies: 'Стек технологий',
        languages: 'Языки',
        hoursSuffix: 'ч',
        fileSuffix: 'резюме',
        contactsHidden: 'Контакты скрыты автором — их можно получить на странице резюме.',
        generated: 'Сформировано {{date}}',
        actualVersion: 'Актуальная версия',
      },

      Nav: {
        about: 'Обо мне',
        experience: 'Опыт',
        skills: 'Навыки',
        education: 'Образование',
        courses: 'Курсы',
        projects: 'Проекты',
        activities: 'Активность',
        contacts: 'Контакты',
      },

      About: {
        title: 'Обо мне',
        relocation: 'Готов к релокации',
        noRelocation: 'Без релокации',

        WorkFormat: {
          fullTime: 'Полная занятость',
          partTime: 'Частичная занятость',
          office: 'Офис',
          remote: 'Удалённо',
          hybrid: 'Гибрид',
        },

        Languages: {
          ru: 'Русский',
          en: 'Английский',
          fr: 'Французский',
          de: 'Немецкий',
          kz: 'Казахский',
          zh: 'Китайский',
          es: 'Испанский',
        },
      },

      Experience: {
        title: 'Опыт работы',
        present: 'по настоящее время',
        showMore: 'Показать ещё {{count}}',
        showLess: 'Скрыть',
        tasks: 'Задачи',
        achievements: 'Достижения',
      },

      TechStack: {
        title: 'Стек технологий',
      },

      Skills: {
        title: 'Навыки',
        hard: 'Профессиональные',
        soft: 'Личные качества',

        HardSkills: {
          architecture: 'Архитектура фронтенда',
          testing: 'Тестирование',
          cicd: 'CI/CD',
          performance: 'Оптимизация производительности',
          accessibility: 'Доступность (a11y)',
          codeReview: 'Code review',
          uiDesign: 'Работа с дизайн-системами',
          api: 'Проектирование API',
        },

        SoftSkills: {
          communication: 'Коммуникация',
          teamwork: 'Командная работа',
          problemSolving: 'Решение проблем',
          timeManagement: 'Тайм-менеджмент',
          mentoring: 'Наставничество',
          adaptability: 'Гибкость',
          ownership: 'Ответственность',
        },
      },

      Education: {
        title: 'Образование',
        course: '{{course}} курс',
        graduated: 'Выпуск {{year}}',

        Levels: {
          level0: 'Среднее образование',
          level1: 'Бакалавриат',
          level2: 'Специалитет',
          level3: 'Магистратура',
          level4: 'Аспирантура',
          level5: 'Кандидат наук',
          level6: 'Доктор наук',
        },
      },

      Courses: {
        title: 'Курсы',
        hours: '{{hours}} ч',
        certificate: 'Сертификат',
      },

      Activities: {
        title: 'Дополнительная активность',
        open: 'Подробнее',

        Types: {
          openSource: 'Open source',
          hackathon: 'Хакатон',
          speaker: 'Выступление',
          volunteering: 'Волонтёрство',
          mentoring: 'Менторство',
          community: 'Сообщество',
        },
      },

      Projects: {
        title: 'Проекты',
        open: 'Открыть',
        code: 'Код',
      },

      AnalyticsNotice: {
        title: 'Анонимная статистика',
        text: 'Автор резюме видит, сколько человек его открыло: время на странице, просмотренные разделы и клики. Имя и почта не собираются, IP-адрес в открытом виде не сохраняется, а повторные заходы различаются по обезличенному идентификатору в cookie.',
        dnt: 'Если в браузере включено «Не отслеживать», статистика не собирается.',
        ok: 'Понятно',
      },

      Contacts: {
        title: 'Давайте работать вместе',
        subtitle: 'Напишите мне — отвечу в ближайшее время',
        telegram: 'Telegram',
        email: 'Почта',
        phone: 'Позвонить',
        copy: 'Скопировать',

        Challenge: {
          title: 'Подтвердите, что вы не робот',
        },

        Reveal: {
          telegram: 'Показать Telegram',
          email: 'Показать почту',
          phone: 'Показать телефон',
        },
      },
    },
  },

  en: {
    CurriculumVitae: {
      Error: {
        notFoundTitle: 'Résumé not found',
        notFoundText: 'The link may be outdated or contain a typo.',
        failedTitle: "Couldn't load the résumé",
        failedText: 'Please try refreshing the page in a moment.',
        errorCode: 'Error code: {{status}}',
        requestId: 'Request ID: {{requestId}}',
        home: 'Go home',
        retry: 'Try again',
        copyDetails: 'Copy details',
        copied: 'Copied',
      },

      Availability: {
        open: 'Open to offers',
        considering: 'Considering offers',
        busy: 'Not looking right now',
        from: 'from {{date}}',
        updated: 'Updated {{time}}',
      },

      Filter: {
        technology: 'Technology',
        experience: 'Experience',
        projects: 'Projects',
        reset: 'Reset filter',
        select: 'Show where {{name}} is used',
      },

      Controls: {
        download: 'Download résumé',
      },

      HtmlExport: {
        technologies: 'Tech stack',
        languages: 'Languages',
        hoursSuffix: 'h',
        fileSuffix: 'résumé',
        contactsHidden:
          'The author keeps the contacts hidden — you can get them on the résumé page.',
        generated: 'Generated on {{date}}',
        actualVersion: 'Latest version',
      },

      Nav: {
        about: 'About',
        experience: 'Experience',
        skills: 'Skills',
        education: 'Education',
        courses: 'Courses',
        projects: 'Projects',
        activities: 'Activities',
        contacts: 'Contacts',
      },

      About: {
        title: 'About me',
        relocation: 'Open to relocation',
        noRelocation: 'Not relocating',

        WorkFormat: {
          fullTime: 'Full-time',
          partTime: 'Part-time',
          office: 'On-site',
          remote: 'Remote',
          hybrid: 'Hybrid',
        },

        Languages: {
          ru: 'Russian',
          en: 'English',
          fr: 'French',
          de: 'German',
          kz: 'Kazakh',
          zh: 'Chinese',
          es: 'Spanish',
        },
      },

      Experience: {
        title: 'Work experience',
        present: 'Present',
        showMore: 'Show {{count}} more',
        showLess: 'Show less',
        tasks: 'Tasks',
        achievements: 'Achievements',
      },

      TechStack: {
        title: 'Tech stack',
      },

      Skills: {
        title: 'Skills',
        hard: 'Professional',
        soft: 'Personal qualities',

        HardSkills: {
          architecture: 'Frontend architecture',
          testing: 'Testing',
          cicd: 'CI/CD',
          performance: 'Performance optimization',
          accessibility: 'Accessibility (a11y)',
          codeReview: 'Code review',
          uiDesign: 'Design systems',
          api: 'API design',
        },

        SoftSkills: {
          communication: 'Communication',
          teamwork: 'Teamwork',
          problemSolving: 'Problem solving',
          timeManagement: 'Time management',
          mentoring: 'Mentoring',
          adaptability: 'Adaptability',
          ownership: 'Ownership',
        },
      },

      Education: {
        title: 'Education',
        course: 'Year {{course}}',
        graduated: 'Class of {{year}}',

        Levels: {
          level0: 'Secondary education',
          level1: "Bachelor's degree",
          level2: 'Specialist degree',
          level3: "Master's degree",
          level4: 'Postgraduate',
          level5: 'PhD',
          level6: 'Doctor of Sciences',
        },
      },

      Courses: {
        title: 'Courses',
        hours: '{{hours}} h',
        certificate: 'Certificate',
      },

      Activities: {
        title: 'Extracurricular activities',
        open: 'Learn more',

        Types: {
          openSource: 'Open source',
          hackathon: 'Hackathon',
          speaker: 'Talk',
          volunteering: 'Volunteering',
          mentoring: 'Mentoring',
          community: 'Community',
        },
      },

      Projects: {
        title: 'Projects',
        open: 'Open',
        code: 'Code',
      },

      AnalyticsNotice: {
        title: 'Anonymous statistics',
        text: 'The résumé author sees how many people opened it: time on the page, sections viewed and clicks. Names and emails are not collected, IP addresses are not stored in plain form, and repeat visits are told apart by an anonymous cookie id.',
        dnt: 'If “Do Not Track” is on in your browser, detailed statistics are not collected.',
        ok: 'Got it',
      },

      Contacts: {
        title: "Let's work together",
        subtitle: "Drop me a message — I'll get back to you soon",
        telegram: 'Telegram',
        email: 'Email',
        phone: 'Call',
        copy: 'Copy',

        Challenge: {
          title: 'Confirm you are not a robot',
        },

        Reveal: {
          telegram: 'Show Telegram',
          email: 'Show email',
          phone: 'Show phone',
        },
      },
    },
  },
} as const
