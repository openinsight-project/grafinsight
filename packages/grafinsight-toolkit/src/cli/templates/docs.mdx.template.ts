export const docsTpl = `import { ArgsTable } from '@storybook/addon-docs/blocks';
import { <%= name %> } from './<%= name %>';

# <%= name %>

### Usage

\`\`\`jsx
import { <%= name %> } from '@grafinsight/ui';

<<%= name %> />
\`\`\`

### Props
<ArgsTable of={<%= name %>} />
`;
