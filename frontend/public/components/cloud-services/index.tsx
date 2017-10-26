/* eslint-disable no-undef */

import * as React from 'react';
import * as _ from 'lodash';

import './ocs-templates';

export { ClusterServiceVersionsDetailsPage, ClusterServiceVersionsPage } from './clusterserviceversion';
export { ClusterServiceVersionResourcesDetailsPage } from './clusterserviceversion-resource';
export { CatalogsDetailsPage } from './catalog';

export enum ALMStatusDescriptors {
  metrics = 'urn:alm:descriptor:com.tectonic.ui:metrics',
  w3Link = 'urn:alm:descriptor:org.w3:link',
  tectonicLink = 'urn:alm:descriptor:com.tectonic.ui:important.link',
  conditions = 'urn:alm:descriptor:io.kubernetes.conditions',
  importantMetrics = 'urn:alm:descriptor:com.tectonic.ui:metrics',
  text = 'urn:alm:descriptor:text',
  prometheus = 'urn:alm:descriptor:io.prometheus:api.v1',

  // Prefix for all kubernetes resource status descriptors.
  k8sResourcePrefix = 'urn:alm:descriptor:io.kubernetes:',
}

export enum ClusterServiceVersionPhase {
  CSVPhaseNone = '',
  CSVPhasePending = 'Pending',
  CSVPhaseInstalling = 'Installing',
  CSVPhaseSucceeded = 'Succeeded',
  CSVPhaseFailed = 'Failed',
  CSVPhaseUnknown = 'Unknown',
}

export enum CSVConditionReason {
  CSVReasonRequirementsUnknown = 'RequirementsUnknown',
  CSVReasonRequirementsNotMet = 'RequirementsNotMet',
  CSVReasonRequirementsMet = 'AllRequirementsMet',
  CSVReasonComponentFailed = 'InstallComponentFailed',
  CSVReasonInstallSuccessful = 'InstallSucceeded',
  CSVReasonInstallCheckFailed = 'InstallCheckFailed',
}

export enum InstallPlanApproval {
  Automatic = 'Automatic',
  Manual = 'Manual',
  UpdateOnly = 'Update-Only',
}

export type K8sResourceKind = {
  apiVersion: string;
  kind: string;
  metadata: {
    annotations?: {[key: string]: string},
    name: string,
    namespace?: string,
    labels?: {[key: string]: string},
    [key: string]: any,
  };
  spec?: {
    selector?: {
      matchLabels?: {[key: string]: any},
    },
    [key: string]: any
  };
  status?: {[key: string]: any};
};

export type CustomResourceDefinitionKind = {

} & K8sResourceKind;

export type ClusterServiceVersionKind = {
  spec: {
    customresourcedefinitions: {owned?: any[], required?: any[]};
  };
  status?: {
    phase: ClusterServiceVersionPhase;
    reason: CSVConditionReason;
  };
} & K8sResourceKind;

export type ClusterServiceVersionResourceKind = {
  status: {[name: string]: any};
} & K8sResourceKind;

export type CatalogEntryKind = {

} & K8sResourceKind;

export type InstallPlanKind = {
  spec: {
    clusterServiceVersionNames: string[];
    approval: 'Automatic' | 'Manual' | 'Update-Only';
  };
  status?: {
    status: 'Planning' | 'Requires Approval' | 'Installing' | 'Complete';
    plan: {
      resolving: string;
      resource: CustomResourceDefinitionKind;
    }[];
  }
} & K8sResourceKind;

export const ClusterServiceVersionLogo: React.StatelessComponent<ClusterServiceVersionLogoProps> = (props) => {
  const {icon, displayName, provider, version} = props;

  return <div className="co-clusterserviceversion-logo">
    <div className="co-clusterserviceversion-logo__icon">{ _.isEmpty(icon)
      ? <span className="fa ci-appcube" style={{fontSize: '40px'}} />
      : <img src={`data:${icon.mediatype};base64,${icon.base64data}`} height="40" width="40" /> }
    </div>
    <div className="co-clusterserviceversion-logo__name">
      <h1 className="co-clusterserviceversion-logo__name__clusterserviceversion">{displayName}</h1>
      { provider && <span className="co-clusterserviceversion-logo__name__provider">{`${version || ''} by ${_.get(provider, 'name', provider)}`}</span> }
    </div>
  </div>;
};

export type ClusterServiceVersionLogoProps = {
  displayName: string;
  icon: {base64data: string, mediatype: string};
  provider: {name: string} | string;
  version?: string;
};

ClusterServiceVersionLogo.displayName = 'ClusterServiceVersionLogo';
