import { DocsId } from '@grafinsight/data';

// TODO: Documentation links
const DOCS_LINKS: Record<DocsId, string> = {
  [DocsId.Transformations]: 'https://grafinsight.com/docs/grafinsight/latest/panels/transformations',
  [DocsId.FieldConfig]: 'https://grafinsight.com/docs/grafinsight/latest/panels/field-configuration-options/',
  [DocsId.FieldConfigOverrides]:
    'https://grafinsight.com/docs/grafinsight/latest/panels/field-configuration-options/#override-a-field',
};

export const getDocsLink = (id: DocsId) => DOCS_LINKS[id];
