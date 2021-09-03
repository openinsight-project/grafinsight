import angular, { ILocationService } from 'angular';
import { ChangeTracker } from './ChangeTracker';
import { ContextSrv } from 'app/core/services/context_srv';
import { DashboardSrv } from './DashboardSrv';
import { GrafInsightRootScope } from 'app/routes/GrafInsightCtrl';

/** @ngInject */
export function unsavedChangesSrv(
  this: any,
  $rootScope: GrafInsightRootScope,
  $location: ILocationService,
  $timeout: any,
  contextSrv: ContextSrv,
  dashboardSrv: DashboardSrv,
  $window: any
) {
  this.init = function (dashboard: any, scope: any) {
    this.tracker = new ChangeTracker(dashboard, scope, 1000, $location, $window, $timeout, contextSrv, $rootScope);
    return this.tracker;
  };
}

angular.module('grafinsight.services').service('unsavedChangesSrv', unsavedChangesSrv);
