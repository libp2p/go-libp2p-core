package metrics

import (
	"fmt"
	"sync"

	"go.opencensus.io/stats/view"
)

// A map of application namespace strings to registered view slices
// Can be considered map[string][]*views.View
var registeredViews = sync.Map{}

type ErrNamespace struct {
	Namespace string
}

// ErrUnregisteredNamespace is an error for lookup of requested unregistered Namespace
type ErrUnregisteredNamespace ErrNamespace

func (e ErrUnregisteredNamespace) Error() string {
	return fmt.Sprintf("no views found registered under Namespace %s", e.Namespace)
}

// ErrDuplicateNamespaceRegistration is an error for a Namespace that has already
// registered views
type ErrDuplicateNamespaceRegistration ErrNamespace

func (e ErrDuplicateNamespaceRegistration) Error() string {
	return fmt.Sprintf("duplicate registration of views by Namespace %s", e.Namespace)
}

// RegisterViews accepts a namespace and a slice of Views, which will be registered
// with opencensus and maintained in the global registered views map
func RegisterViews(namespace string, views ...*view.View) error {
	_, loaded := registeredViews.LoadOrStore(namespace, views)
	if loaded {
		return ErrDuplicateNamespaceRegistration{Namespace: namespace}
	}

	if err := view.Register(views...); err != nil {
		registeredViews.Delete(namespace)
		return err
	}

	return nil
}

// LookupViews returns all views for a Namespace name. Returns an error if the
// Namespace has not been registered.
func LookupViews(name string) ([]*view.View, error) {
	views, ok := registeredViews.Load(name)
	if !ok {
		return nil, ErrUnregisteredNamespace{Namespace: name}
	}
	return views.([]*view.View), nil
}

// AllViews returns all registered views as a single slice
func AllViews() []*view.View {
	var views []*view.View
	registeredViews.Range(func(key interface{}, values interface{}) bool {
		views = append(views, values.([]*view.View)...)
		return true
	})
	return views
}

func AllViews() (views []*view.View) {
	for _, v := range registeredViews {
		views = append(views, v...)
	}
	return
}
