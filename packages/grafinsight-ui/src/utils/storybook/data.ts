import { applyFieldOverrides, DataFrame, GrafInsightTheme } from '@grafinsight/data';

export function prepDataForStorybook(data: DataFrame[], theme: GrafInsightTheme) {
  return applyFieldOverrides({
    data: data,
    fieldConfig: {
      overrides: [],
      defaults: {},
    },
    theme,
    replaceVariables: (value: string) => value,
  });
}
