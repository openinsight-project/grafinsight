import React, { useCallback } from 'react';
import { DataFrame, DisplayValue, fieldReducers, getFieldDisplayName, reduceField } from '@grafinsight/data';
import { UPlotConfigBuilder } from './config/UPlotConfigBuilder';
import { VizLegendItem, VizLegendOptions } from '../VizLegend/types';
import { AxisPlacement } from './config';
import { VizLayout, VizLayoutLegendProps } from '../VizLayout/VizLayout';
import { mapMouseEventToMode } from '../GraphNG/utils';
import { VizLegend } from '../VizLegend/VizLegend';
import { GraphNGLegendEvent } from '..';

const defaultFormatter = (v: any) => (v == null ? '-' : v.toFixed(1));

interface PlotLegendProps extends VizLegendOptions, Omit<VizLayoutLegendProps, 'children'> {
  data: DataFrame[];
  config: UPlotConfigBuilder;
  onSeriesColorChange?: (label: string, color: string) => void;
  onLegendClick?: (event: GraphNGLegendEvent) => void;
}

export const PlotLegend: React.FC<PlotLegendProps> = ({
  data,
  config,
  onSeriesColorChange,
  onLegendClick,
  placement,
  calcs,
  displayMode,
  ...vizLayoutLegendProps
}) => {
  const onLegendLabelClick = useCallback(
    (legend: VizLegendItem, event: React.MouseEvent) => {
      const { fieldIndex } = legend;

      if (!onLegendClick || !fieldIndex) {
        return;
      }

      onLegendClick({
        fieldIndex,
        mode: mapMouseEventToMode(event),
      });
    },
    [onLegendClick]
  );

  const legendItems = config
    .getSeries()
    .map<VizLegendItem | undefined>((s) => {
      const seriesConfig = s.props;
      const fieldIndex = seriesConfig.dataFrameFieldIndex;
      const axisPlacement = config.getAxisPlacement(s.props.scaleKey);

      if (seriesConfig.hideInLegend || !fieldIndex) {
        return undefined;
      }

      const field = data[fieldIndex.frameIndex]?.fields[fieldIndex.fieldIndex];

      if (!field) {
        return undefined;
      }

      const label = getFieldDisplayName(field, data[fieldIndex.frameIndex]!);

      return {
        disabled: !seriesConfig.show ?? false,
        fieldIndex,
        color: seriesConfig.lineColor!,
        label,
        yAxis: axisPlacement === AxisPlacement.Left ? 1 : 2,
        getDisplayValues: () => {
          if (!calcs?.length) {
            return [];
          }

          const fmt = field.display ?? defaultFormatter;
          const fieldCalcs = reduceField({
            field,
            reducers: calcs,
          });

          return calcs.map<DisplayValue>((reducer) => {
            return {
              ...fmt(fieldCalcs[reducer]),
              title: fieldReducers.get(reducer).name,
            };
          });
        },
        getItemKey: () => `${label}-${fieldIndex.frameIndex}-${fieldIndex.fieldIndex}`,
      };
    })
    .filter((i) => i !== undefined) as VizLegendItem[];

  return (
    <VizLayout.Legend placement={placement} {...vizLayoutLegendProps}>
      <VizLegend
        onLabelClick={onLegendLabelClick}
        placement={placement}
        items={legendItems}
        displayMode={displayMode}
        onSeriesColorChange={onSeriesColorChange}
      />
    </VizLayout.Legend>
  );
};

PlotLegend.displayName = 'PlotLegend';
