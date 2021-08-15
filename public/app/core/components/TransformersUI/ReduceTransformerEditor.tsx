import React, { useCallback } from 'react';
import { LegacyForms, Select, StatsPicker } from '@grafinsight/ui';
import {
  DataTransformerID,
  ReducerID,
  SelectableValue,
  standardTransformers,
  TransformerRegistryItem,
  TransformerUIProps,
} from '@grafinsight/data';

import { ReduceTransformerMode, ReduceTransformerOptions } from '@grafinsight/data/src/transformations/transformers/reduce';
import { selectors } from '@grafinsight/e2e-selectors/src';

// TODO:  Minimal implementation, needs some <3
export const ReduceTransformerEditor: React.FC<TransformerUIProps<ReduceTransformerOptions>> = ({
  options,
  onChange,
}) => {
  const modes: Array<SelectableValue<ReduceTransformerMode>> = [
    {
      label: 'Series to rows',
      value: ReduceTransformerMode.SeriesToRows,
      description: 'Create a table with one row for each series value',
    },
    {
      label: 'Reduce fields',
      value: ReduceTransformerMode.ReduceFields,
      description: 'Collapse each field into a single value',
    },
  ];

  const onSelectMode = useCallback(
    (value: SelectableValue<ReduceTransformerMode>) => {
      const mode = value.value!;
      onChange({
        ...options,
        mode,
        includeTimeField: mode === ReduceTransformerMode.ReduceFields ? !!options.includeTimeField : false,
      });
    },
    [onChange, options]
  );

  const onToggleTime = useCallback(() => {
    onChange({
      ...options,
      includeTimeField: !options.includeTimeField,
    });
  }, [onChange, options]);

  return (
    <>
      <div>
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label width-8" aria-label={selectors.components.Transforms.Reduce.modeLabel}>
            Mode
          </div>
          <Select
            options={modes}
            value={modes.find((v) => v.value === options.mode) || modes[0]}
            onChange={onSelectMode}
            className="flex-grow-1"
          />
        </div>
      </div>
      <div className="gf-form-inline">
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label width-8" aria-label={selectors.components.Transforms.Reduce.calculationsLabel}>
            Calculations
          </div>
          <StatsPicker
            className="flex-grow-1"
            placeholder="Choose Stat"
            allowMultiple
            stats={options.reducers || []}
            onChange={(stats) => {
              onChange({
                ...options,
                reducers: stats as ReducerID[],
              });
            }}
          />
        </div>
      </div>
      {options.mode === ReduceTransformerMode.ReduceFields && (
        <div className="gf-form-inline">
          <div className="gf-form">
            <LegacyForms.Switch
              label="Include time"
              labelClass="width-8"
              checked={!!options.includeTimeField}
              onChange={onToggleTime}
            />
          </div>
        </div>
      )}
    </>
  );
};

export const reduceTransformRegistryItem: TransformerRegistryItem<ReduceTransformerOptions> = {
  id: DataTransformerID.reduce,
  editor: ReduceTransformerEditor,
  transformation: standardTransformers.reduceTransformer,
  name: standardTransformers.reduceTransformer.name,
  description: standardTransformers.reduceTransformer.description,
};
