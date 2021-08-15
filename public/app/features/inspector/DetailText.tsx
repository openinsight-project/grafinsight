import React, { FC } from 'react';
import { useStyles } from '@grafinsight/ui';
import { GrafInsightTheme } from '@grafinsight/data';
import { css } from 'emotion';

const getStyles = (theme: GrafInsightTheme) => css`
  margin: 0;
  margin-left: ${theme.spacing.md};
  font-size: ${theme.typography.size.sm};
  color: ${theme.colors.textWeak};
`;

export const DetailText: FC = ({ children }) => {
  const collapsedTextStyles = useStyles(getStyles);
  return <p className={collapsedTextStyles}>{children}</p>;
};
