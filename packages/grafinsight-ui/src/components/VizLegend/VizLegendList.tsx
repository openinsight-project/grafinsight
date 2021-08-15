import React from 'react';
import { VizLegendBaseProps, VizLegendItem } from './types';
import { InlineList } from '../List/InlineList';
import { List } from '../List/List';
import { css, cx } from 'emotion';
import { useStyles } from '../../themes';
import { GrafInsightTheme } from '@grafinsight/data';
import { VizLegendListItem } from './VizLegendListItem';

export interface Props extends VizLegendBaseProps {}

/**
 * @internal
 */
export const VizLegendList: React.FunctionComponent<Props> = ({
  items,
  itemRenderer,
  onSeriesColorChange,
  onLabelClick,
  placement,
  className,
}) => {
  const styles = useStyles(getStyles);

  if (!itemRenderer) {
    /* eslint-disable-next-line react/display-name */
    itemRenderer = (item) => (
      <VizLegendListItem item={item} onLabelClick={onLabelClick} onSeriesColorChange={onSeriesColorChange} />
    );
  }

  const getItemKey = (item: VizLegendItem) => `${item.label}`;

  switch (placement) {
    case 'right': {
      const renderItem = (item: VizLegendItem, index: number) => {
        return <span className={styles.itemRight}>{itemRenderer!(item, index)}</span>;
      };

      return (
        <div className={cx(styles.rightWrapper, className)}>
          <List items={items} renderItem={renderItem} getItemKey={getItemKey} />
        </div>
      );
    }
    case 'bottom':
    default: {
      const renderItem = (item: VizLegendItem, index: number) => {
        return <span className={styles.itemBottom}>{itemRenderer!(item, index)}</span>;
      };

      return (
        <div className={cx(styles.bottomWrapper, className)}>
          <div className={styles.section}>
            <InlineList
              items={items.filter((item) => item.yAxis === 1)}
              renderItem={renderItem}
              getItemKey={getItemKey}
            />
          </div>
          <div className={cx(styles.section, styles.sectionRight)}>
            <InlineList
              items={items.filter((item) => item.yAxis !== 1)}
              renderItem={renderItem}
              getItemKey={getItemKey}
            />
          </div>
        </div>
      );
    }
  }
};

VizLegendList.displayName = 'VizLegendList';

const getStyles = (theme: GrafInsightTheme) => {
  const itemStyles = css`
    padding-right: 10px;
    display: flex;
    font-size: ${theme.typography.size.sm};
    white-space: nowrap;
  `;

  return {
    itemBottom: itemStyles,
    itemRight: cx(
      itemStyles,
      css`
        margin-bottom: ${theme.spacing.xs};
      `
    ),
    rightWrapper: css`
      padding-left: ${theme.spacing.sm};
    `,
    bottomWrapper: css`
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      width: 100%;
      padding-left: ${theme.spacing.md};
    `,
    section: css`
      display: flex;
    `,
    sectionRight: css`
      justify-content: flex-end;
      flex-grow: 1;
    `,
  };
};
