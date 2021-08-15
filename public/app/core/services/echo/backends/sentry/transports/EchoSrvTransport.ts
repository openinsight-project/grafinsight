import { getEchoSrv, EchoEventType } from '@grafinsight/runtime/src';
import { BaseTransport } from '@sentry/browser/dist/transports';
import { Event } from '@sentry/browser';
import { Status } from '@sentry/types';

export class EchoSrvTransport extends BaseTransport {
  sendEvent(event: Event) {
    getEchoSrv().addEvent({
      type: EchoEventType.Sentry,
      payload: event,
    });
    return Promise.resolve({ status: Status.Success, event });
  }
}
