import React, { FC, HTMLAttributes, ReactNode } from 'react';
import { css } from 'emotion';
import { GrafInsightTheme } from '@grafinsight/data';
import { selectors } from '@grafinsight/e2e-selectors/src';
import { useTheme } from '../../themes';
import { Icon } from '../Icon/Icon';
import { IconName } from '../../types/icon';
import { getColorsFromSeverity } from '../../utils/colors';

export type AlertVariant = 'success' | 'warning' | 'error' | 'info';

export interface Props extends HTMLAttributes<HTMLDivElement> {
  title: string;
  /** On click handler for alert button, mostly used for dismissing the alert */
  onRemove?: (event: React.MouseEvent) => void;
  severity?: AlertVariant;
  children?: ReactNode;
  /** Custom component or text for alert button */
  buttonContent?: ReactNode | string;
  /** @deprecated */
  /** Deprecated use onRemove instead */
  onButtonClick?: (event: React.MouseEvent) => void;
  /** @deprecated */
  /** Deprecated use buttonContent instead */
  buttonText?: string;
}

function getIconFromSeverity(severity: AlertVariant): string {
  switch (severity) {
    case 'error':
    case 'warning':
      return 'exclamation-triangle';
    case 'info':
      return 'info-circle';
    case 'success':
      return 'check';
    default:
      return '';
  }
}

export const Alert: FC<Props> = React.forwardRef<HTMLDivElement, Props>(
  ({ title, buttonText, onButtonClick, onRemove, children, buttonContent, severity = 'error', ...restProps }, ref) => {
    const theme = useTheme();
    const styles = getStyles(theme, severity, !!buttonContent);

    return (
      <div ref={ref} className={styles.alert} aria-label={selectors.components.Alert.alert(severity)} {...restProps}>
        <div className={styles.icon}>
          <Icon size="xl" name={getIconFromSeverity(severity) as IconName} />
        </div>
        <div className={styles.body}>
          <div className={styles.title}>{title}</div>
          {children && <div>{children}</div>}
        </div>
        {/* If onRemove is specified, giving preference to onRemove */}
        {onRemove ? (
          <button type="button" className={styles.close} onClick={onRemove}>
            {buttonContent || <Icon name="times" size="lg" />}
          </button>
        ) : onButtonClick ? (
          <button type="button" className="btn btn-outline-danger" onClick={onButtonClick}>
            {buttonText}
          </button>
        ) : null}
      </div>
    );
  }
);

Alert.displayName = 'Alert';

const getStyles = (theme: GrafInsightTheme, severity: AlertVariant, outline: boolean) => {
  const { white } = theme.palette;
  const severityColors = getColorsFromSeverity(severity, theme);
  const background = css`
    background: linear-gradient(90deg, ${severityColors[0]}, ${severityColors[0]});
  `;

  return {
    alert: css`
      flex-grow: 1;
      padding: 15px 20px;
      margin-bottom: ${theme.spacing.xs};
      position: relative;
      color: ${white};
      text-shadow: 0 1px 0 rgba(0, 0, 0, 0.2);
      border-radius: ${theme.border.radius.md};
      display: flex;
      flex-direction: row;
      align-items: center;
      ${background}
    `,
    icon: css`
      padding: 0 ${theme.spacing.md} 0 0;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 35px;
    `,
    title: css`
      font-weight: ${theme.typography.weight.semibold};
    `,
    body: css`
      flex-grow: 1;
      margin: 0 ${theme.spacing.md} 0 0;
      overflow-wrap: break-word;
      word-break: break-word;

      a {
        color: ${white};
        text-decoration: underline;
      }
    `,
    close: css`
      background: none;
      display: flex;
      align-items: center;
      border: ${outline ? `1px solid ${white}` : 'none'};
      border-radius: ${theme.border.radius.sm};
    `,
  };
};
