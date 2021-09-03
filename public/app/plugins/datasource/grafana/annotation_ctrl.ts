import { SelectableValue } from '@grafinsight/data';
import { GrafInsightAnnotationType } from './types';

export const annotationTypes: Array<SelectableValue<GrafInsightAnnotationType>> = [
  { text: 'Dashboard', value: GrafInsightAnnotationType.Dashboard },
  { text: 'Tags', value: GrafInsightAnnotationType.Tags },
];

export class GrafInsightAnnotationsQueryCtrl {
  annotation: any;

  types = annotationTypes;

  constructor() {
    this.annotation.type = this.annotation.type || GrafInsightAnnotationType.Tags;
    this.annotation.limit = this.annotation.limit || 100;
  }

  static templateUrl = 'partials/annotations.editor.html';
}
