import { PanelPlugin } from '@grafinsight/data';
import { WelcomeBanner } from './Welcome';

export const plugin = new PanelPlugin(WelcomeBanner).setNoPadding();
