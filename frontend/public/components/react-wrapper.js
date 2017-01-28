import React from 'react';

import {k8sEnum} from '../module/k8s';

export const angulars = {
  store: null,
  ModalLauncherSvc: null,
  $location: null,
  $log: null,
  $interval: null,
  $timeout: null,
};

const app = angular.module('bridge.react-wrapper', ['bridge']);

const toRegister = [];
export const register = (name, Component) => {
  if (app && app.value) {
    return app.value(name, Component);
  }
  toRegister.push({name, Component});
};

app.value('nop', () => <div/>);

app.service('angularBridge', function ($ngRedux, $location, $routeParams, $timeout, $interval, $log, ModalLauncherSvc, errorMessageSvc) {
  // "Export" angular modules to the outside world via ref through 'angulars'...
  // NOTE: this only exist after the app has loaded!

  this.expose = () => {
    _.map(toRegister, ({name, Component}) => {
      app.value(name, Component);
    });

    angulars.store = $ngRedux;
    angulars.ModalLauncherSvc = ModalLauncherSvc;
    angulars.modal = (...args) => () => ModalLauncherSvc.open(...args),
    angulars.$location = $location;
    angulars.routeParams = $routeParams;
    angulars.$log = $log;
    angulars.$interval= $interval;
    angulars.$timeout = $timeout;
    angulars.errorMessageSvc = errorMessageSvc;
  };
});

// see https://github.com/ngReact/ngReact#the-react-component-directive
app.directive('reactiveK8sList', function () {
  return {
    template: '<react-component name="{{component}}" props="props"></react-component>',
    restrict: 'E',
    scope: {
      kind: '=',
      // A React Component that has been registered with angular
      component: '=',
      canCreate: '=',
      selector: '=',
      fieldSelector: '=',
      selectorRequired: '=',
    },
    controller: function ($routeParams, $scope) {
      const { kind, canCreate, selector, fieldSelector, component, selectorRequired } = $scope;

      $scope.props = {
        kind, canCreate, selector, fieldSelector, component, selectorRequired,
        namespace: $routeParams.ns,
        defaultNS: k8sEnum.DefaultNS,
        name: $routeParams.name,
        location: location.pathname,
      };
    }
  };
});
