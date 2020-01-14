// Definitions for the Kubernetes Controllers
package controller

import (
	. "github.com/solo-io/wasme/pkg/operator/api/wasme.io/v1"

	"github.com/pkg/errors"
	"github.com/solo-io/autopilot/pkg/events"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type FilterDeploymentEventHandler interface {
	Create(obj *FilterDeployment) error
	Update(old, new *FilterDeployment) error
	Delete(obj *FilterDeployment) error
	Generic(obj *FilterDeployment) error
}

type FilterDeploymentEventHandlerFuncs struct {
	OnCreate  func(obj *FilterDeployment) error
	OnUpdate  func(old, new *FilterDeployment) error
	OnDelete  func(obj *FilterDeployment) error
	OnGeneric func(obj *FilterDeployment) error
}

func (f *FilterDeploymentEventHandlerFuncs) Create(obj *FilterDeployment) error {
	if f.OnCreate == nil {
		return nil
	}
	return f.OnCreate(obj)
}

func (f *FilterDeploymentEventHandlerFuncs) Delete(obj *FilterDeployment) error {
	if f.OnDelete == nil {
		return nil
	}
	return f.OnDelete(obj)
}

func (f *FilterDeploymentEventHandlerFuncs) Update(objOld, objNew *FilterDeployment) error {
	if f.OnUpdate == nil {
		return nil
	}
	return f.OnUpdate(objOld, objNew)
}

func (f *FilterDeploymentEventHandlerFuncs) Generic(obj *FilterDeployment) error {
	if f.OnGeneric == nil {
		return nil
	}
	return f.OnGeneric(obj)
}

type FilterDeploymentController struct {
	watcher events.EventWatcher
}

func NewFilterDeploymentController(name string, mgr manager.Manager) (*FilterDeploymentController, error) {
	if err := AddToScheme(mgr.GetScheme()); err != nil {
		return nil, err
	}

	w, err := events.NewWatcher(name, mgr)
	if err != nil {
		return nil, err
	}
	return &FilterDeploymentController{
		watcher: w,
	}, nil
}

func (c *FilterDeploymentController) AddEventHandler(h FilterDeploymentEventHandler, predicates ...predicate.Predicate) error {
	handler := genericFilterDeploymentHandler{handler: h}
	if err := c.watcher.Watch(&FilterDeployment{}, handler, predicates...); err != nil {
		return err
	}
	return nil
}

// genericFilterDeploymentHandler implements a generic events.EventHandler
type genericFilterDeploymentHandler struct {
	handler FilterDeploymentEventHandler
}

func (h genericFilterDeploymentHandler) Create(object runtime.Object) error {
	obj, ok := object.(*FilterDeployment)
	if !ok {
		return errors.Errorf("internal error: FilterDeployment handler received event for %T")
	}
	return h.handler.Create(obj)
}

func (h genericFilterDeploymentHandler) Delete(object runtime.Object) error {
	obj, ok := object.(*FilterDeployment)
	if !ok {
		return errors.Errorf("internal error: FilterDeployment handler received event for %T")
	}
	return h.handler.Delete(obj)
}

func (h genericFilterDeploymentHandler) Update(old, new runtime.Object) error {
	objOld, ok := old.(*FilterDeployment)
	if !ok {
		return errors.Errorf("internal error: FilterDeployment handler received event for %T")
	}
	objNew, ok := new.(*FilterDeployment)
	if !ok {
		return errors.Errorf("internal error: FilterDeployment handler received event for %T")
	}
	return h.handler.Update(objOld, objNew)
}

func (h genericFilterDeploymentHandler) Generic(object runtime.Object) error {
	obj, ok := object.(*FilterDeployment)
	if !ok {
		return errors.Errorf("internal error: FilterDeployment handler received event for %T")
	}
	return h.handler.Generic(obj)
}
