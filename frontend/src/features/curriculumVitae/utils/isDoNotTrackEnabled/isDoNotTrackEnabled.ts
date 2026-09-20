type NavigatorPrivacy = { doNotTrack?: string | null; globalPrivacyControl?: boolean }

export const isDoNotTrackEnabled = (
  nav: NavigatorPrivacy | undefined = typeof navigator === 'undefined' ? undefined : navigator,
) => nav?.doNotTrack === '1' || nav?.globalPrivacyControl === true
