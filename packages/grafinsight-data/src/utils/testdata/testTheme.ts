import { GrafInsightTheme, GrafInsightThemeType } from '../../types/theme';

export function getTestTheme(type: GrafInsightThemeType = GrafInsightThemeType.Dark): GrafInsightTheme {
  return ({
    type,
    isDark: type === GrafInsightThemeType.Dark,
    isLight: type === GrafInsightThemeType.Light,
    colors: {
      panelBg: 'white',
    },
  } as unknown) as GrafInsightTheme;
}
