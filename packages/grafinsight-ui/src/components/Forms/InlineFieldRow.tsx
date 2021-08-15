import React, { FC, HTMLProps, ReactNode } from 'react';
import { css, cx } from 'emotion';
import { useStyles } from '../../themes';
import { GrafInsightTheme } from '@grafinsight/data';

export interface Props extends Omit<HTMLProps<HTMLDivElement>, 'css'> {
  children: ReactNode | ReactNode[];
}

export const InlineFieldRow: FC<Props> = ({ children, className, ...htmlProps }) => {
  const styles = useStyles(getStyles);
  return (
    <div className={cx(styles.container, className)} {...htmlProps}>
      {children}
    </div>
  );
};

const getStyles = (theme: GrafInsightTheme) => {
  return {
    container: css`
      label: InlineFieldRow;
      display: flex;
      flex-direction: row;
      flex-wrap: wrap;
      align-content: flex-start;
    `,
  };
};
