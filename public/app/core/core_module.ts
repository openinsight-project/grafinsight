import angular from 'angular';

const coreModule = angular.module('grafinsight.core', ['ngRoute']);

// legacy modules
const angularModules = [
  coreModule,
  angular.module('grafinsight.controllers', []),
  angular.module('grafinsight.directives', []),
  angular.module('grafinsight.factories', []),
  angular.module('grafinsight.services', []),
  angular.module('grafinsight.filters', []),
  angular.module('grafinsight.routes', []),
];

export { angularModules, coreModule };

export default coreModule;
