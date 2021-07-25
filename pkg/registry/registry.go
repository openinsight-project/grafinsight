package registry

import (
	"context"
	"reflect"
	"sort"

	"github.com/openinsight-project/grafinsight/pkg/services/sqlstore/migrator"
)

type Descriptor struct {
	Name         string
	Instance     Service
	InitPriority Priority
}

var services []*Descriptor

func RegisterServiceWithPriority(instance Service, priority Priority) {
	Register(&Descriptor{
		Name:         reflect.TypeOf(instance).Elem().Name(),
		Instance:     instance,
		InitPriority: priority,
	})
}

func RegisterService(instance Service) {
	Register(&Descriptor{
		Name:         reflect.TypeOf(instance).Elem().Name(),
		Instance:     instance,
		InitPriority: Medium,
	})
}

func Register(descriptor *Descriptor) {
	if descriptor == nil {
		return
	}
	// Overwrite any existing equivalent service
	for i, svc := range services {
		if svc.Name == descriptor.Name {
			services[i] = descriptor
			return
		}
	}

	services = append(services, descriptor)
}

// GetService gets the registered service descriptor with a certain name.
// If none is found, nil is returned.
func GetService(name string) *Descriptor {
	for _, svc := range services {
		if svc.Name == name {
			return svc
		}
	}

	return nil
}

func GetServices() []*Descriptor {
	slice := getServicesWithOverrides()

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].InitPriority > slice[j].InitPriority
	})

	return slice
}

type OverrideServiceFunc func(descriptor Descriptor) (*Descriptor, bool)

var overrides []OverrideServiceFunc

func RegisterOverride(fn OverrideServiceFunc) {
	overrides = append(overrides, fn)
}

func ClearOverrides() {
	overrides = nil
}

func getServicesWithOverrides() []*Descriptor {
	slice := []*Descriptor{}
	for _, s := range services {
		var descriptor *Descriptor
		for _, fn := range overrides {
			if newDescriptor, override := fn(*s); override {
				descriptor = newDescriptor
				break
			}
		}

		if descriptor != nil {
			slice = append(slice, descriptor)
		} else {
			slice = append(slice, s)
		}
	}

	return slice
}

// Service interface is the lowest common shape that services
// are expected to fulfill to be started within Grafana.
type Service interface {
	// Init is called by Grafana main process which gives the service
	// the possibility do some initial work before its started. Things
	// like adding routes, bus handlers should be done in the Init function
	Init() error
}

// CanBeDisabled allows the services to decide if it should
// be started or not by itself. This is useful for services
// that might not always be started, ex alerting.
// This will be called after `Init()`.
type CanBeDisabled interface {
	// IsDisabled should return a bool saying if it can be started or not.
	IsDisabled() bool
}

// BackgroundService should be implemented for services that have
// long running tasks in the background.
type BackgroundService interface {
	// Run starts the background process of the service after `Init` have been called
	// on all services. The `context.Context` passed into the function should be used
	// to subscribe to ctx.Done() so the service can be notified when Grafana shuts down.
	Run(ctx context.Context) error
}

// DatabaseMigrator allows the caller to add migrations to
// the migrator passed as argument
type DatabaseMigrator interface {
	// AddMigrations allows the service to add migrations to
	// the database migrator.
	AddMigration(mg *migrator.Migrator)
}

// IsDisabled returns whether a service is disabled.
func IsDisabled(srv Service) bool {
	canBeDisabled, ok := srv.(CanBeDisabled)
	return ok && canBeDisabled.IsDisabled()
}

type Priority int

const (
	High       Priority = 100
	MediumHigh Priority = 75
	Medium     Priority = 50
	Low        Priority = 0
)
