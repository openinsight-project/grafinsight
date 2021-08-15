import { PanelPlugin } from '@grafinsight/data';
import { NodeGraphPanel } from './NodeGraphPanel';
import { Options } from './types';

export const plugin = new PanelPlugin<Options>(NodeGraphPanel);
