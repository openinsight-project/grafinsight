import { config, GrafInsightBootConfig } from '@grafinsight/runtime';
// Legacy binding paths
export { config, GrafInsightBootConfig as Settings };

let grafinsightConfig: GrafInsightBootConfig = config;

export default grafinsightConfig;

export const getConfig = () => {
  return grafinsightConfig;
};

export const updateConfig = (update: Partial<GrafInsightBootConfig>) => {
  grafinsightConfig = {
    ...grafinsightConfig,
    ...update,
  };
};
