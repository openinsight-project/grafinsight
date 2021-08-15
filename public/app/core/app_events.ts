import { EventBusSrv, EventBusExtended } from '@grafinsight/data';

export const appEvents: EventBusExtended = new EventBusSrv();

export default appEvents;
