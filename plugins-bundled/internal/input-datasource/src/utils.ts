import { toDataFrame, DataFrameDTO, toCSV } from '@grafinsight/data';

export function dataFrameToCSV(dto?: DataFrameDTO[]) {
  if (!dto || !dto.length) {
    return '';
  }
  return toCSV(dto.map((v) => toDataFrame(v)));
}
