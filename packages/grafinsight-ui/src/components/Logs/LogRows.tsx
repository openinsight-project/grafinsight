import React, { PureComponent } from 'react';
import memoizeOne from 'memoize-one';
import { TimeZone, LogsDedupStrategy, LogRowModel, Field, LinkModel, LogsSortOrder, sortLogRows } from '@grafinsight/data';

import { Themeable } from '../../types/theme';
import { withTheme } from '../../themes/index';
import { getLogRowStyles } from './getLogRowStyles';

//Components
import { LogRow } from './LogRow';
import { RowContextOptions } from './LogRowContextProvider';

export const PREVIEW_LIMIT = 100;

export interface Props extends Themeable {
  logRows?: LogRowModel[];
  deduplicatedRows?: LogRowModel[];
  dedupStrategy: LogsDedupStrategy;
  highlighterExpressions?: string[];
  showLabels: boolean;
  showTime: boolean;
  wrapLogMessage: boolean;
  timeZone: TimeZone;
  logsSortOrder?: LogsSortOrder | null;
  allowDetails?: boolean;
  previewLimit?: number;
  // Passed to fix problems with inactive scrolling in Logs Panel
  // Can be removed when we unify scrolling for Panel and Explore
  disableCustomHorizontalScroll?: boolean;
  forceEscape?: boolean;
  showDetectedFields?: string[];
  showContextToggle?: (row?: LogRowModel) => boolean;
  onClickFilterLabel?: (key: string, value: string) => void;
  onClickFilterOutLabel?: (key: string, value: string) => void;
  getRowContext?: (row: LogRowModel, options?: RowContextOptions) => Promise<any>;
  getFieldLinks?: (field: Field, rowIndex: number) => Array<LinkModel<Field>>;
  onClickShowDetectedField?: (key: string) => void;
  onClickHideDetectedField?: (key: string) => void;
}

interface State {
  renderAll: boolean;
}

class UnThemedLogRows extends PureComponent<Props, State> {
  renderAllTimer: number | null = null;

  static defaultProps = {
    previewLimit: PREVIEW_LIMIT,
  };

  state: State = {
    renderAll: false,
  };

  componentDidMount() {
    // Staged rendering
    const { logRows, previewLimit } = this.props;
    const rowCount = logRows ? logRows.length : 0;
    // Render all right away if not too far over the limit
    const renderAll = rowCount <= previewLimit! * 2;
    if (renderAll) {
      this.setState({ renderAll });
    } else {
      this.renderAllTimer = window.setTimeout(() => this.setState({ renderAll: true }), 2000);
    }
  }

  componentWillUnmount() {
    if (this.renderAllTimer) {
      clearTimeout(this.renderAllTimer);
    }
  }

  makeGetRows = memoizeOne((orderedRows: LogRowModel[]) => {
    return () => orderedRows;
  });

  sortLogs = memoizeOne((logRows: LogRowModel[], logsSortOrder: LogsSortOrder): LogRowModel[] =>
    sortLogRows(logRows, logsSortOrder)
  );

  render() {
    const {
      dedupStrategy,
      showContextToggle,
      showLabels,
      showTime,
      wrapLogMessage,
      logRows,
      deduplicatedRows,
      highlighterExpressions,
      timeZone,
      onClickFilterLabel,
      onClickFilterOutLabel,
      theme,
      allowDetails,
      previewLimit,
      getFieldLinks,
      disableCustomHorizontalScroll,
      logsSortOrder,
      showDetectedFields,
      onClickShowDetectedField,
      onClickHideDetectedField,
      forceEscape,
    } = this.props;
    const { renderAll } = this.state;
    const { logsRowsTable, logsRowsHorizontalScroll } = getLogRowStyles(theme);
    const dedupedRows = deduplicatedRows ? deduplicatedRows : logRows;
    const hasData = logRows && logRows.length > 0;
    const dedupCount = dedupedRows
      ? dedupedRows.reduce((sum, row) => (row.duplicates ? sum + row.duplicates : sum), 0)
      : 0;
    const showDuplicates = dedupStrategy !== LogsDedupStrategy.none && dedupCount > 0;

    // For horizontal scrolling we can't use CustomScrollbar as it causes the problem with logs context - it is not visible
    // for top log rows. Therefore we use CustomScrollbar only in LogsPanel and for Explore, we use custom css styling.
    const horizontalScrollWindow = wrapLogMessage || disableCustomHorizontalScroll ? '' : logsRowsHorizontalScroll;

    // Staged rendering
    const processedRows = dedupedRows ? dedupedRows : [];
    const orderedRows = logsSortOrder ? this.sortLogs(processedRows, logsSortOrder) : processedRows;
    const firstRows = orderedRows.slice(0, previewLimit!);
    const lastRows = orderedRows.slice(previewLimit!, orderedRows.length);

    // React profiler becomes unusable if we pass all rows to all rows and their labels, using getter instead
    const getRows = this.makeGetRows(orderedRows);
    const getRowContext = this.props.getRowContext ? this.props.getRowContext : () => Promise.resolve([]);

    return (
      <div className={horizontalScrollWindow}>
        <table className={logsRowsTable}>
          <tbody>
            {hasData &&
              firstRows.map((row, index) => (
                <LogRow
                  key={row.uid}
                  getRows={getRows}
                  getRowContext={getRowContext}
                  highlighterExpressions={highlighterExpressions}
                  row={row}
                  showContextToggle={showContextToggle}
                  showDuplicates={showDuplicates}
                  showLabels={showLabels}
                  showTime={showTime}
                  showDetectedFields={showDetectedFields}
                  wrapLogMessage={wrapLogMessage}
                  timeZone={timeZone}
                  allowDetails={allowDetails}
                  onClickFilterLabel={onClickFilterLabel}
                  onClickFilterOutLabel={onClickFilterOutLabel}
                  onClickShowDetectedField={onClickShowDetectedField}
                  onClickHideDetectedField={onClickHideDetectedField}
                  getFieldLinks={getFieldLinks}
                  logsSortOrder={logsSortOrder}
                  forceEscape={forceEscape}
                />
              ))}
            {hasData &&
              renderAll &&
              lastRows.map((row, index) => (
                <LogRow
                  key={row.uid}
                  getRows={getRows}
                  getRowContext={getRowContext}
                  row={row}
                  showContextToggle={showContextToggle}
                  showDuplicates={showDuplicates}
                  showLabels={showLabels}
                  showTime={showTime}
                  showDetectedFields={showDetectedFields}
                  wrapLogMessage={wrapLogMessage}
                  timeZone={timeZone}
                  allowDetails={allowDetails}
                  onClickFilterLabel={onClickFilterLabel}
                  onClickFilterOutLabel={onClickFilterOutLabel}
                  onClickShowDetectedField={onClickShowDetectedField}
                  onClickHideDetectedField={onClickHideDetectedField}
                  getFieldLinks={getFieldLinks}
                  logsSortOrder={logsSortOrder}
                  forceEscape={forceEscape}
                />
              ))}
            {hasData && !renderAll && (
              <tr>
                <td colSpan={5}>Rendering {orderedRows.length - previewLimit!} rows...</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    );
  }
}

export const LogRows = withTheme(UnThemedLogRows);
LogRows.displayName = 'LogsRows';
