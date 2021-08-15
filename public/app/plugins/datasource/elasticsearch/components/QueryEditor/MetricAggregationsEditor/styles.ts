import { GrafInsightTheme } from '@grafinsight/data';
import { stylesFactory } from '@grafinsight/ui';
import { css } from 'emotion';

export const getStyles = stylesFactory((theme: GrafInsightTheme, hidden: boolean) => ({
  color:
    hidden &&
    css`
      &,
      &:hover,
      label,
      a {
        color: ${hidden ? theme.colors.textFaint : theme.colors.text};
      }
    `,
}));
