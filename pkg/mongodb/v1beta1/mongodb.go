// Package v1beta1 of mongodb provides a methods for the mongodb compont
// IBM Confidential
// OCO Source Materials
// 5737-E67
// (C) Copyright IBM Corporation 2019 All Rights Reserved
// The source code for this program is not published or otherwise divested of its trade secrets, irrespective of what has been deposited with the U.S. Copyright Office.
package v1beta1

import (
	"context"

	multicloudv1beta1 "github.ibm.com/IBMPrivateCloud/ibm-klusterlet-operator/pkg/apis/multicloud/v1beta1"

	"k8s.io/apimachinery/pkg/api/errors"

	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("mongodb")

// Create MongoDB CR
func Create(instance *multicloudv1beta1.Endpoint, cr *multicloudv1beta1.MongoDB, client client.Client) error {
	log.Info("Creating a new MongoDB", "MongoDB.Namespace", cr.Namespace, "MongoDB.Name", cr.Name)
	err := client.Create(context.TODO(), cr)
	if err != nil {
		log.Error(err, "Fail to CREATE MongoDB CR")
		return err
	}

	// Adding Finalizer to Instance
	instance.Finalizers = append(instance.Finalizers, cr.Name)
	return nil
}

// Update MongoDB CR
func Update(instance *multicloudv1beta1.Endpoint, cr *multicloudv1beta1.MongoDB, foundCR *multicloudv1beta1.MongoDB, client client.Client) error {
	foundCR.Spec = cr.Spec
	err := client.Update(context.TODO(), foundCR)
	if err != nil && !errors.IsConflict(err) {
		log.Error(err, "Fail to UPDATE MongoDB CR")
		return err
	}

	// Adding Finalizer to Instance if Finalizer does not exist
	// NOTE: This is to handle requeue due to failed instance update during creation
	for _, finalizer := range instance.Finalizers {
		if finalizer == cr.Name {
			return nil
		}
	}
	instance.Finalizers = append(instance.Finalizers, cr.Name)
	return nil
}

// Delete MongoDB CR
func Delete(foundCR *multicloudv1beta1.MongoDB, client client.Client) error {
	return client.Delete(context.TODO(), foundCR)
}

// Finalize MongoDB CR
func Finalize(instance *multicloudv1beta1.Endpoint, cr *multicloudv1beta1.MongoDB, client client.Client) error {
	for i, finalizer := range instance.Finalizers {
		if finalizer == cr.Name {
			// Remove finalizer
			instance.Finalizers = append(instance.Finalizers[0:i], instance.Finalizers[i+1:]...)
			return nil
		}
	}
	return nil
}