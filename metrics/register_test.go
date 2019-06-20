package metrics

import (
	"fmt"
	"go.opencensus.io/stats"
	"testing"

	"go.opencensus.io/stats/view"
)

func newTestMeasure(name string) stats.Measure {
	return stats.Int64(fmt.Sprintf("test/measure/%s", name),
		fmt.Sprintf("Test measurement %s", name),
		stats.UnitDimensionless,
	)
}

func newTestView(name string) *view.View {
	return &view.View{
		Name:        fmt.Sprintf("test/%s", name),
		Description: fmt.Sprintf("Test view %s", name),
		Measure:     newTestMeasure(name),
		Aggregation: view.LastValue(),
	}
}

func TestRegisteringViews(t *testing.T) {
	registeredViews = make(map[string][]*view.View)

	t.Run("test registering first views", func(t *testing.T) {
		testView := newTestView("empty-map-0")

		if err := RegisterViews("test", testView); err != nil {
			t.Fatal("unable to register view in empty map", err)
		}
	})

	t.Run("test registering with existing views", func(t *testing.T) {
		testView := newTestView("empty-map-1")
		testView2 := newTestView("existing-entity-0")

		if err := RegisterViews("test2", testView); err != nil {
			t.Fatal("unable to register view in empty map", err)
		}
		if err := RegisterViews("test3", testView2); err != nil {
			t.Fatal("unable to register view in map", err)
		}
	})

	t.Run("test registering duplicate views", func(t *testing.T) {
		testView := newTestView("empty-map-2")
		testView2 := newTestView("existing-entity-1")

		if err := RegisterViews("test4", testView); err != nil {
			t.Fatal("unable to register view in empty map", err)
		}
		if err := RegisterViews("test4", testView2); err == nil {
			t.Fatal("allowed duplicate view registration")
		}
	})

	t.Run("test looking up views", func(t *testing.T) {
		testView := newTestView("empty-map-3")

		if err := RegisterViews("test5", testView); err != nil {
			t.Fatal("unable to register view in empty map", err)
		}

		views, err := LookupViews("test5")
		if err != nil {
			t.Fatal("error looking up views", err)
		}

		if views[0].Name != testView.Name {
			t.Fatal("incorrect view lookup, received name:", views[0].Name)
		}
	})
}

func TestAllViews(t *testing.T) {
	registeredViews = make(map[string][]*view.View)
	t.Run("test retrieving all views", func(t *testing.T) {
		views := []*view.View{newTestView("all-views-0"), newTestView("all-views-1"), newTestView("all-views-2")}

		if err := RegisterViews("test6", views...); err != nil {
			t.Fatal("unable to register multiple views at once", err)
		}

		allViews := AllViews()
		if len(allViews) != len(views) {
			t.Fatalf("didn't receive equal number of views: %d %d", len(views), len(allViews))
		}
	})
}
