/**
 * Все конечные наборы значений резюме. Значение enum совпадает с кодом, который присылает бэкенд
 * (например, `WorkFormat.Remote === 'remote'`), а подписи для интерфейса лежат в словаре по этим кодам.
 */

/** Формат работы */
export enum WorkFormat {
  FullTime = 'fullTime',
  PartTime = 'partTime',
  Office = 'office',
  Remote = 'remote',
  Hybrid = 'hybrid',
}

/** Язык, которым владеет человек */
export enum LanguageType {
  Ru = 'ru',
  En = 'en',
  Fr = 'fr',
  De = 'de',
  Kz = 'kz',
  Zh = 'zh',
  Es = 'es',
}

/** Уровень знания языка */
export enum LanguageGrade {
  A1 = 'A1',
  A2 = 'A2',
  B1 = 'B1',
  B2 = 'B2',
  C1 = 'C1',
  C2 = 'C2',
}

/** Личные качества */
export enum SoftSkill {
  Communication = 'communication',
  Teamwork = 'teamwork',
  ProblemSolving = 'problemSolving',
  TimeManagement = 'timeManagement',
  Mentoring = 'mentoring',
  Adaptability = 'adaptability',
  Ownership = 'ownership',
}

/** Профессиональные навыки */
export enum HardSkill {
  Architecture = 'architecture',
  Testing = 'testing',
  Cicd = 'cicd',
  Performance = 'performance',
  Accessibility = 'accessibility',
  CodeReview = 'codeReview',
  UiDesign = 'uiDesign',
  Api = 'api',
}

/** Уровень образования (в JSON приходит числом) */
export enum EducationLevel {
  Secondary = 0,
  Bachelor = 1,
  Specialist = 2,
  Master = 3,
  Postgraduate = 4,
  Candidate = 5,
  Doctor = 6,
}

/** Способ связи */
export enum ContactKey {
  Telegram = 'telegram',
  Email = 'email',
  Phone = 'phone',
}

/** Вид дополнительной активности */
export enum ActivityType {
  OpenSource = 'openSource',
  Hackathon = 'hackathon',
  Speaker = 'speaker',
  Volunteering = 'volunteering',
  Mentoring = 'mentoring',
  Community = 'community',
}

/** Готовность к предложениям: ищу, рассматриваю, не ищу */
export enum AvailabilityStatus {
  Open = 'open',
  Considering = 'considering',
  Busy = 'busy',
}

/** Секции страницы (hero всегда первый и не настраивается) */
export enum SectionId {
  About = 'about',
  TechStack = 'techStack',
  Experience = 'experience',
  Skills = 'skills',
  Education = 'education',
  Courses = 'courses',
  Projects = 'projects',
  Activities = 'activities',
  Contacts = 'contacts',
}
