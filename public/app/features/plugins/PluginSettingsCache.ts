import { getBackendSrv } from '@grafinsight/runtime/src';
import { PluginMeta } from '@grafinsight/data';

type PluginCache = {
  [key: string]: PluginMeta;
};

const pluginInfoCache: PluginCache = {};

export function getPluginSettings(pluginId: string): Promise<PluginMeta> {
  const v = pluginInfoCache[pluginId];
  if (v) {
    return Promise.resolve(v);
  }
  return getBackendSrv()
    .get(`/api/plugins/${pluginId}/settings`)
    .then((settings: any) => {
      pluginInfoCache[pluginId] = settings;
      return settings;
    })
    .catch((err: any) => {
      // err.isHandled = true;
      return Promise.reject('Unknown Plugin');
    });
}
