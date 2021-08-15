import React, { PureComponent } from 'react';
import { QueryEditorHelpProps } from '@grafinsight/data';
import InfluxCheatSheet from './InfluxCheatSheet';

export default class InfluxStartPage extends PureComponent<QueryEditorHelpProps> {
  render() {
    return <InfluxCheatSheet onClickExample={this.props.onClickExample} />;
  }
}
