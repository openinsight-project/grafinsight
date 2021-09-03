import { configure } from 'enzyme';
import Adapter from '@wojtekmaj/enzyme-adapter-react-17';
import $ from 'jquery';
import 'mutationobserver-shim';

const global = window as any;
global.$ = global.jQuery = $;

import '../vendor/flot/jquery.flot';
import '../vendor/flot/jquery.flot.time';
import angular from 'angular';

angular.module('grafinsight', ['ngRoute']);
angular.module('grafinsight.services', ['ngRoute', '$strap.directives']);
angular.module('grafinsight.panels', []);
angular.module('grafinsight.controllers', []);
angular.module('grafinsight.directives', []);
angular.module('grafinsight.filters', []);
angular.module('grafinsight.routes', ['ngRoute']);

jest.mock('app/core/core', () => ({}));
jest.mock('app/features/plugins/plugin_loader', () => ({}));

configure({ adapter: new Adapter() });

const localStorageMock = (() => {
  let store: any = {};
  return {
    getItem: (key: string) => {
      return store[key];
    },
    setItem: (key: string, value: any) => {
      store[key] = value.toString();
    },
    clear: () => {
      store = {};
    },
    removeItem: (key: string) => {
      delete store[key];
    },
  };
})();

global.localStorage = localStorageMock;

const throwUnhandledRejections = () => {
  process.on('unhandledRejection', (err) => {
    throw err;
  });
};

throwUnhandledRejections();
