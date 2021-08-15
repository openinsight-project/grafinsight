import React, { FC } from 'react';
import { css } from 'emotion';
import { stylesFactory, useStyles } from '../../themes';
import { GrafInsightTheme, locale } from '@grafinsight/data';

const getStyles = stylesFactory((theme: GrafInsightTheme) => {
  return {
    counter: css`
      label: counter;
      margin-left: ${theme.spacing.sm};
      border-radius: ${theme.spacing.lg};
      background-color: ${theme.colors.bg2};
      padding: ${theme.spacing.xxs} ${theme.spacing.sm};
      color: ${theme.colors.textWeak};
      font-weight: ${theme.typography.weight.semibold};
      font-size: ${theme.typography.size.sm};
    `,
  };
});

export interface CounterProps {
  value: number;
}

export const Counter: FC<CounterProps> = ({ value }) => {
  const styles = useStyles(getStyles);

  return <span className={styles.counter}>{locale(value, 0).text}</span>;
};
