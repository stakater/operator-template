package actions

import (
	"context"
	"github.com/stakater/operator-boilerplate/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IAction interface {
	Run() error
	GetName() string
}

type GenericAction struct {
	Name string
	Fn   func() error
}

func (a *GenericAction) GetName() string {
	return a.Name
}

func (a *GenericAction) Run() error {
	return a.Fn()
}

type CreateResourceAction struct {
	*utils.ReconcilerBase
	context.Context
	Name   string
	Object client.Object
}

func (a *CreateResourceAction) GetName() string {
	return a.Name
}

func (a *CreateResourceAction) Run() error {
	return a.CreateResource(a.Context, a.Object)
}

type CreateOwnedResourceAction struct {
	CreateResourceAction
	Owner client.Object
}

func (a *CreateOwnedResourceAction) Run() error {
	return a.CreateOwnedResource(a.Context, a.Owner, a.Object)
}
