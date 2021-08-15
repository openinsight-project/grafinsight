import { ScopedVars } from '@grafinsight/data';
import { getTemplateSrv } from '@grafinsight/runtime/src';
import { variableAdapters } from './adapters';

export function getAllVariableValuesForUrl(scopedVars?: ScopedVars) {
  const params: Record<string, string | string[]> = {};
  const variables = getTemplateSrv().getVariables();

  // console.log(variables)
  for (let i = 0; i < variables.length; i++) {
    const variable = variables[i];
    if (scopedVars && scopedVars[variable.name] !== void 0) {
      if (scopedVars[variable.name].skipUrlSync) {
        continue;
      }
      params['var-' + variable.name] = scopedVars[variable.name].value;
    } else {
      // @ts-ignore
      if (variable.skipUrlSync) {
        continue;
      }
      params['var-' + variable.name] = variableAdapters.get(variable.type).getValueForUrl(variable as any);
    }
  }

  return params;
}
