import React from 'react';
import { DataTransformerID, standardTransformers, TransformerRegistryItem, TransformerUIProps } from '@grafinsight/data';
import { SeriesToRowsTransformerOptions } from '@grafinsight/data/src/transformations/transformers/seriesToRows';

export const SeriesToRowsTransformerEditor: React.FC<TransformerUIProps<SeriesToRowsTransformerOptions>> = ({
  input,
  options,
  onChange,
}) => {
  return null;
};

export const seriesToRowsTransformerRegistryItem: TransformerRegistryItem<SeriesToRowsTransformerOptions> = {
  id: DataTransformerID.seriesToRows,
  editor: SeriesToRowsTransformerEditor,
  transformation: standardTransformers.seriesToRowsTransformer,
  name: 'Series to rows',
  description: `Merge many series and return a single series with time, metric and value as columns.
                Useful for showing multiple time series visualized in a table.`,
};
