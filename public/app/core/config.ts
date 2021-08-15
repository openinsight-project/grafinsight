import { config, GrafInsightBootConfig } from '@grafinsight/runtime';
// Legacy binding paths
export { config, GrafInsightBootConfig as Settings };

let grafanaConfig: GrafInsightBootConfig = config;

export default grafanaConfig;

export const getConfig = () => {
  return grafanaConfig;
};

export const updateConfig = (update: Partial<GrafInsightBootConfig>) => {
  grafanaConfig = {
    ...grafanaConfig,
    ...update,
  };
};
