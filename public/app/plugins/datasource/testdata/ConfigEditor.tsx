// Libraries
import React, { PureComponent } from 'react';

import { DataSourcePluginOptionsEditorProps } from '@grafinsight/data';

type Props = DataSourcePluginOptionsEditorProps<any>;

/**
 * Empty Config Editor -- settings to save
 */
export class ConfigEditor extends PureComponent<Props> {
  render() {
    return <div />;
  }
}
