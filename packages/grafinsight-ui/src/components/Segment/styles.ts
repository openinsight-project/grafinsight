import { GrafInsightTheme } from '@grafinsight/data';
import { css } from 'emotion';

export const getSegmentStyles = (theme: GrafInsightTheme) => {
  return {
    segment: css`
      cursor: pointer;
      width: auto;
    `,

    queryPlaceholder: css`
      color: ${theme.palette.gray2};
    `,

    disabled: css`
      cursor: not-allowed;
      opacity: 0.65;
      box-shadow: none;
    `,
  };
};
