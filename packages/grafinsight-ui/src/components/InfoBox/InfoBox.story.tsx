import React from 'react';
import { FeatureState } from '@grafinsight/data';
import { InfoBox, FeatureInfoBox } from '@grafinsight/ui';
import mdx from './InfoBox.mdx';
import { Story } from '@storybook/react';
import { FeatureInfoBoxProps } from './FeatureInfoBox';
import { InfoBoxProps } from './InfoBox';

export default {
  title: 'Layout/InfoBox',
  component: InfoBox,
  decorators: [],
  parameters: {
    docs: {
      page: mdx,
    },
  },
  argTypes: {
    onDismiss: { action: 'Dismissed' },
    featureState: {
      control: { type: 'select', options: ['alpha', 'beta', undefined] },
    },
    children: {
      table: {
        disable: true,
      },
    },
  },
};

const defaultProps: FeatureInfoBoxProps = {
  title: 'A title',
  severity: 'info',
  url: 'http://www.grafinsight.com',
  featureState: FeatureState.beta,

  children: (
    <p>
      The database user should only be granted SELECT permissions on the specified database &amp; tables you want to
      query. GrafInsight does not validate that queries are safe so queries can contain any SQL statement. For example,
      statements like <code>USE otherdb;</code> and <code>DROP TABLE user;</code> would be executed. To protect against
      this we <strong>Highly</strong> recommend you create a specific MySQL user with restricted permissions.
    </p>
  ),
};

const InfoBoxTemplate: Story<InfoBoxProps> = (args) => <InfoBox {...args} />;
export const infoBox = InfoBoxTemplate.bind({});
infoBox.args = defaultProps;

const FeatureInfoBoxTemplate: Story<FeatureInfoBoxProps> = (args) => <FeatureInfoBox {...args}></FeatureInfoBox>;
export const featureInfoBox = FeatureInfoBoxTemplate.bind({});
featureInfoBox.args = defaultProps;
