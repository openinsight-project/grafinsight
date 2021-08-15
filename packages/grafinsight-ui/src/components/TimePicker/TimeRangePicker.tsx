// Libraries
import React, { PureComponent, memo, FormEvent } from 'react';
import { css } from 'emotion';

// Components
import { Tooltip } from '../Tooltip/Tooltip';
import { TimePickerContent } from './TimeRangePicker/TimePickerContent';
import { ClickOutsideWrapper } from '../ClickOutsideWrapper/ClickOutsideWrapper';

// Utils & Services
import { stylesFactory } from '../../themes/stylesFactory';
import { withTheme, useTheme } from '../../themes/ThemeContext';

// Types
import {
  isDateTime,
  rangeUtil,
  GrafInsightTheme,
  dateTimeFormat,
  timeZoneFormatUserFriendly,
  TimeRange,
  TimeZone,
  dateMath,
} from '@grafinsight/data';
import { Themeable } from '../../types';
import { otherOptions, quickOptions } from './rangeOptions';
import { ButtonGroup, ToolbarButton } from '../Button';

export interface TimeRangePickerProps extends Themeable {
  hideText?: boolean;
  value: TimeRange;
  timeZone?: TimeZone;
  timeSyncButton?: JSX.Element;
  isSynced?: boolean;
  onChange: (timeRange: TimeRange) => void;
  onChangeTimeZone: (timeZone: TimeZone) => void;
  onMoveBackward: () => void;
  onMoveForward: () => void;
  onZoom: () => void;
  history?: TimeRange[];
  hideQuickRanges?: boolean;
}

export interface State {
  isOpen: boolean;
}

export class UnthemedTimeRangePicker extends PureComponent<TimeRangePickerProps, State> {
  state: State = {
    isOpen: false,
  };

  onChange = (timeRange: TimeRange) => {
    this.props.onChange(timeRange);
    this.setState({ isOpen: false });
  };

  onOpen = (event: FormEvent<HTMLButtonElement>) => {
    const { isOpen } = this.state;
    event.stopPropagation();
    event.preventDefault();
    this.setState({ isOpen: !isOpen });
  };

  onClose = () => {
    this.setState({ isOpen: false });
  };

  render() {
    const {
      value,
      onMoveBackward,
      onMoveForward,
      onZoom,
      timeZone,
      timeSyncButton,
      isSynced,
      theme,
      history,
      onChangeTimeZone,
      hideQuickRanges,
    } = this.props;

    const { isOpen } = this.state;
    const styles = getStyles(theme);
    const hasAbsolute = isDateTime(value.raw.from) || isDateTime(value.raw.to);
    const variant = isSynced ? 'active' : 'default';

    return (
      <ButtonGroup className={styles.container}>
        {hasAbsolute && <ToolbarButton variant={variant} onClick={onMoveBackward} icon="angle-left" narrow />}

        <Tooltip content={<TimePickerTooltip timeRange={value} timeZone={timeZone} />} placement="bottom">
          <ToolbarButton
            aria-label="TimePicker Open Button"
            onClick={this.onOpen}
            icon="clock-nine"
            isOpen={isOpen}
            variant={variant}
          >
            <TimePickerButtonLabel {...this.props} />
          </ToolbarButton>
        </Tooltip>
        {isOpen && (
          <ClickOutsideWrapper includeButtonPress={false} onClick={this.onClose}>
            <TimePickerContent
              timeZone={timeZone}
              value={value}
              onChange={this.onChange}
              otherOptions={otherOptions}
              quickOptions={quickOptions}
              history={history}
              showHistory
              onChangeTimeZone={onChangeTimeZone}
              hideQuickRanges={hideQuickRanges}
            />
          </ClickOutsideWrapper>
        )}

        {timeSyncButton}

        {hasAbsolute && <ToolbarButton onClick={onMoveForward} icon="angle-right" narrow variant={variant} />}

        <Tooltip content={ZoomOutTooltip} placement="bottom">
          <ToolbarButton onClick={onZoom} icon="search-minus" variant={variant} />
        </Tooltip>
      </ButtonGroup>
    );
  }
}

const ZoomOutTooltip = () => (
  <>
    Time range zoom out <br /> CTRL+Z
  </>
);

const TimePickerTooltip = ({ timeRange, timeZone }: { timeRange: TimeRange; timeZone?: TimeZone }) => {
  const theme = useTheme();
  const styles = getLabelStyles(theme);

  return (
    <>
      {dateTimeFormat(timeRange.from, { timeZone })}
      <div className="text-center">to</div>
      {dateTimeFormat(timeRange.to, { timeZone })}
      <div className="text-center">
        <span className={styles.utc}>{timeZoneFormatUserFriendly(timeZone)}</span>
      </div>
    </>
  );
};

type LabelProps = Pick<TimeRangePickerProps, 'hideText' | 'value' | 'timeZone'>;

export const TimePickerButtonLabel = memo<LabelProps>(({ hideText, value, timeZone }) => {
  const theme = useTheme();
  const styles = getLabelStyles(theme);

  if (hideText) {
    return null;
  }

  return (
    <span className={styles.container}>
      <span>{formattedRange(value, timeZone)}</span>
      <span className={styles.utc}>{rangeUtil.describeTimeRangeAbbreviation(value, timeZone)}</span>
    </span>
  );
});

TimePickerButtonLabel.displayName = 'TimePickerButtonLabel';

const formattedRange = (value: TimeRange, timeZone?: TimeZone) => {
  const adjustedTimeRange = {
    to: dateMath.isMathString(value.raw.to) ? value.raw.to : value.to,
    from: dateMath.isMathString(value.raw.from) ? value.raw.from : value.from,
  };
  return rangeUtil.describeTimeRange(adjustedTimeRange, timeZone);
};

export const TimeRangePicker = withTheme(UnthemedTimeRangePicker);

const getStyles = stylesFactory((theme: GrafInsightTheme) => {
  return {
    container: css`
      position: relative;
      display: flex;
      vertical-align: middle;
    `,
  };
});

const getLabelStyles = stylesFactory((theme: GrafInsightTheme) => {
  return {
    container: css`
      display: flex;
      align-items: center;
      white-space: nowrap;
    `,
    utc: css`
      color: ${theme.palette.orange};
      font-size: ${theme.typography.size.sm};
      padding-left: 6px;
      line-height: 28px;
      vertical-align: bottom;
      font-weight: ${theme.typography.weight.semibold};
    `,
  };
});
