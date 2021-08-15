import { LiveChannel, LiveChannelAddress } from '@grafinsight/data';
import { Observable } from 'rxjs';

/**
 * @alpha -- experimental
 */
export interface GrafInsightLiveSrv {
  /**
   * Is the server currently connected
   */
  isConnected(): boolean;

  /**
   * Listen for changes to the main service
   */
  getConnectionState(): Observable<boolean>;

  /**
   * Get a channel.  If the scope, namespace, or path is invalid, a shutdown
   * channel will be returned with an error state indicated in its status.
   *
   * This is a singleton instance that stays active until explicitly shutdown.
   * Multiple requests for this channel will return the same object until
   * the channel is shutdown
   */
  getChannel<TMessage, TPublish = any>(address: LiveChannelAddress): LiveChannel<TMessage, TPublish>;
}

let singletonInstance: GrafInsightLiveSrv;

/**
 * Used during startup by GrafInsight to set the GrafInsightLiveSrv so it is available
 * via the {@link getGrafInsightLiveSrv} to the rest of the application.
 *
 * @internal
 */
export const setGrafInsightLiveSrv = (instance: GrafInsightLiveSrv) => {
  singletonInstance = instance;
};

/**
 * Used to retrieve the GrafInsightLiveSrv that allows you to subscribe to
 * server side events and streams
 *
 * @alpha -- experimental
 */
export const getGrafInsightLiveSrv = (): GrafInsightLiveSrv => singletonInstance;
