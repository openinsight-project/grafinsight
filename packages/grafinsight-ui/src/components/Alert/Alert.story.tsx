import React from 'react';
import { Story } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { Alert, AlertVariant, VerticalGroup } from '@grafinsight/ui';
import { Props } from './Alert';
import { withCenteredStory, withHorizontallyCenteredStory } from '../../utils/storybook/withCenteredStory';
import mdx from '../Alert/Alert.mdx';
import { StoryExample } from '../../utils/storybook/StoryExample';

const severities: AlertVariant[] = ['error', 'warning', 'info', 'success'];

export default {
  title: 'Overlays/Alert',
  component: Alert,
  decorators: [withCenteredStory, withHorizontallyCenteredStory],
  parameters: {
    docs: {
      page: mdx,
    },
    knobs: {
      disable: true,
    },
  },
  argTypes: {
    severity: { control: { type: 'select', options: severities } },
  },
};

export const Examples: Story<Props> = ({ severity, title, buttonContent }) => {
  return (
    <VerticalGroup>
      <StoryExample name="With buttonContent and children">
        <Alert
          title={title}
          severity={severity}
          buttonContent={<span>{buttonContent}</span>}
          onRemove={action('Remove button clicked')}
        >
          Child content that includes some alert details, like maybe what actually happened.
        </Alert>
      </StoryExample>
      <StoryExample name="No dismiss">
        <Alert title={title} severity={severity} />
      </StoryExample>
      <StoryExample name="Severities">
        <VerticalGroup>
          {severities.map((severity) => (
            <Alert
              title={`Severity: ${severity}`}
              severity={severity}
              key={severity}
              onRemove={action('Remove button clicked')}
            />
          ))}
        </VerticalGroup>
      </StoryExample>
    </VerticalGroup>
  );
};

Examples.args = {
  severity: 'error',
  title: 'Some very important message',
  buttonContent: 'Close',
};
