import { DataFrameFieldIndex, FieldMatcher } from '@grafinsight/data';

/**
 * Mode to describe if a legend is isolated/selected or being appended to an existing
 * series selection.
 * @alpha
 */
export enum GraphNGLegendEventMode {
  ToggleSelection = 'select',
  AppendToSelection = 'append',
}

/**
 * Event being triggered when the user interact with the Graph legend.
 * @alpha
 */
export interface GraphNGLegendEvent {
  fieldIndex: DataFrameFieldIndex;
  mode: GraphNGLegendEventMode;
}

export interface XYFieldMatchers {
  x: FieldMatcher; // first match
  y: FieldMatcher;
}
