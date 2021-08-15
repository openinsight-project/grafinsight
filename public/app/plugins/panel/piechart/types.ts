import { PieChartType, SingleStatBaseOptions, PieChartLabels, PieChartLegendOptions } from '@grafinsight/ui';

export interface PieChartOptions extends SingleStatBaseOptions {
  pieType: PieChartType;
  displayLabels: PieChartLabels[];
  legend: PieChartLegendOptions;
}
