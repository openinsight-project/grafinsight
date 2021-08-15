import { ValueMapping, MappingType, ValueMap, RangeMap } from '../types';

type TimeSeriesValue = string | number | null;

const addValueToTextMappingText = (
  allValueMappings: ValueMapping[],
  valueToTextMapping: ValueMap,
  value: TimeSeriesValue
) => {
  if (valueToTextMapping.value === undefined) {
    return allValueMappings;
  }

  if (value === null && isNullValueMap(valueToTextMapping)) {
    return allValueMappings.concat(valueToTextMapping);
  }

  let valueAsNumber, valueToTextMappingAsNumber;

  if (isNumeric(value as string) && isNumeric(valueToTextMapping.value)) {
    valueAsNumber = parseFloat(value as string);
    valueToTextMappingAsNumber = parseFloat(valueToTextMapping.value as string);

    if (valueAsNumber === valueToTextMappingAsNumber) {
      return allValueMappings.concat(valueToTextMapping);
    }
    return allValueMappings;
  }

  if (value === valueToTextMapping.value) {
    return allValueMappings.concat(valueToTextMapping);
  }

  return allValueMappings;
};

const addRangeToTextMappingText = (
  allValueMappings: ValueMapping[],
  rangeToTextMapping: RangeMap,
  value: TimeSeriesValue
) => {
  if (rangeToTextMapping.from === undefined || rangeToTextMapping.to === undefined || value === undefined) {
    return allValueMappings;
  }

  if (
    value === null &&
    rangeToTextMapping.from &&
    rangeToTextMapping.to &&
    rangeToTextMapping.from.toLowerCase() === 'null' &&
    rangeToTextMapping.to.toLowerCase() === 'null'
  ) {
    return allValueMappings.concat(rangeToTextMapping);
  }

  const valueAsNumber = parseFloat(value as string);
  const fromAsNumber = parseFloat(rangeToTextMapping.from as string);
  const toAsNumber = parseFloat(rangeToTextMapping.to as string);

  if (isNaN(valueAsNumber) || isNaN(fromAsNumber) || isNaN(toAsNumber)) {
    return allValueMappings;
  }

  if (valueAsNumber >= fromAsNumber && valueAsNumber <= toAsNumber) {
    return allValueMappings.concat(rangeToTextMapping);
  }

  return allValueMappings;
};

const getAllFormattedValueMappings = (valueMappings: ValueMapping[], value: TimeSeriesValue) => {
  const allFormattedValueMappings = valueMappings.reduce((allValueMappings, valueMapping) => {
    if (valueMapping.type === MappingType.ValueToText) {
      allValueMappings = addValueToTextMappingText(allValueMappings, valueMapping as ValueMap, value);
    } else if (valueMapping.type === MappingType.RangeToText) {
      allValueMappings = addRangeToTextMappingText(allValueMappings, valueMapping as RangeMap, value);
    }

    return allValueMappings;
  }, [] as ValueMapping[]);

  allFormattedValueMappings.sort((t1, t2) => {
    return t1.id - t2.id;
  });

  return allFormattedValueMappings;
};

export const getMappedValue = (valueMappings: ValueMapping[], value: TimeSeriesValue): ValueMapping => {
  return getAllFormattedValueMappings(valueMappings, value)[0];
};

const isNullValueMap = (mapping: ValueMap): boolean => {
  if (!mapping || !mapping.value) {
    return false;
  }
  return mapping.value.toLowerCase() === 'null';
};

// Ref https://stackoverflow.com/a/58550111

export function isNumeric(num: any) {
  return (typeof num === 'number' || (typeof num === 'string' && num.trim() !== '')) && !isNaN(num as number);
}
