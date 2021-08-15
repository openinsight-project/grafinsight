import { GrafInsightThemeType } from '@grafinsight/data';

type VariantDescriptor = { [key in GrafInsightThemeType]: string | number };

/**
 * @deprecated use theme.isLight ? or theme.isDark instead
 */
export const selectThemeVariant = (variants: VariantDescriptor, currentTheme?: GrafInsightThemeType) => {
  return variants[currentTheme || GrafInsightThemeType.Dark];
};
