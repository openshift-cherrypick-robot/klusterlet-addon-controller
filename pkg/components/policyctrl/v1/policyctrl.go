// (c) Copyright IBM Corporation 2019, 2020. All Rights Reserved.
// Note to U.S. Government Users Restricted Rights:
// U.S. Government Users Restricted Rights - Use, duplication or disclosure restricted by GSA ADP Schedule
// Contract with IBM Corp.
// Licensed Materials - Property of IBM
//
// Copyright (c) 2020 Red Hat, Inc.

package v1

import (
	agentv1 "github.com/open-cluster-management/endpoint-operator/pkg/apis/agent/v1"
	addonoperator "github.com/open-cluster-management/endpoint-operator/pkg/components/addon-operator/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// constants for policy controller
const (
	PolicyController = "klusterlet-addon-policyctrl"
	PolicyCtrl       = "policyctrl"
)

var log = logf.Log.WithName("policyctrl")

// IsEnabled - check whether policyctrl is enabled
func IsEnabled(instance *agentv1.KlusterletAddonConfig) bool {
	return instance.Spec.PolicyController.Enabled
}

// NewPolicyControllerCR - create CR for component poliicy controller
func NewPolicyControllerCR(instance *agentv1.KlusterletAddonConfig,
	client client.Client, namespace string) (*agentv1.PolicyController, error) {
	labels := map[string]string{
		"app": instance.Name,
	}

	gv := agentv1.GlobalValues{
		ImagePullPolicy: instance.Spec.ImagePullPolicy,
		ImagePullSecret: instance.Spec.ImagePullSecret,
		ImageOverrides:  make(map[string]string, 1),
	}

	imageKey, imageRepository, err := instance.GetImage("policy-controller")
	if err != nil {
		log.Error(err, "Fail to get Image", "Component.Name", "policy-controller")
		return nil, err
	}

	gv.ImageOverrides[imageKey] = imageRepository

	imageKey, imageRepository, err = instance.GetImage("governance-policy-spec-sync")
	if err != nil {
		log.Error(err, "Fail to get Image", "Component.Name", "governance-policy-spec-sync")
		return nil, err
	}

	gv.ImageOverrides[imageKey] = imageRepository

	imageKey, imageRepository, err = instance.GetImage("governance-policy-status-sync")
	if err != nil {
		log.Error(err, "Fail to get Image", "Component.Name", "governance-policy-status-sync")
		return nil, err
	}

	gv.ImageOverrides[imageKey] = imageRepository

	imageKey, imageRepository, err = instance.GetImage("governance-policy-template-sync")
	if err != nil {
		log.Error(err, "Fail to get Image", "Component.Name", "governance-policy-template-sync")
		return nil, err
	}

	gv.ImageOverrides[imageKey] = imageRepository

	return &agentv1.PolicyController{
		TypeMeta: metav1.TypeMeta{
			APIVersion: agentv1.SchemeGroupVersion.String(),
			Kind:       "PolicyController",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      PolicyController,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: agentv1.PolicyControllerSpec{
			FullNameOverride:            PolicyController,
			ClusterName:                 instance.Spec.ClusterName,
			ClusterNamespace:            instance.Spec.ClusterNamespace,
			HubKubeconfigSecret:         PolicyCtrl + "-hub-kubeconfig",
			GlobalValues:                gv,
			DeployedOnHub:               false,
			PostDeleteJobServiceAccount: addonoperator.KlusterletAddonOperator,
		},
	}, nil
}
