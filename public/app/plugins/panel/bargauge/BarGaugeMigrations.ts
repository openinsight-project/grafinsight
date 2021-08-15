import { PanelModel } from '@grafinsight/data';
import { sharedSingleStatMigrationHandler } from '@grafinsight/ui';
import { BarGaugeOptions } from './types';

export const barGaugePanelMigrationHandler = (panel: PanelModel<BarGaugeOptions>): Partial<BarGaugeOptions> => {
  return sharedSingleStatMigrationHandler(panel);
};
