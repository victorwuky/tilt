package buildcontrol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/windmilleng/tilt/internal/k8s/testyaml"
	"github.com/windmilleng/tilt/internal/store"
	"github.com/windmilleng/tilt/internal/testutils/manifestbuilder"
	"github.com/windmilleng/tilt/internal/testutils/tempdir"
	"github.com/windmilleng/tilt/pkg/model"
)

func TestNextTargetToBuildDoesntReturnCurrentlyBuildingTarget(t *testing.T) {
	f := newTestFixture(t)
	defer f.TearDown()

	mt := f.manifestNeedingCrashRebuild()
	f.st.UpsertManifestTarget(mt)

	// Verify this target is normally next-to-build
	f.assertNextTargetToBuild(mt.Manifest.Name)

	// If target is currently building, should NOT be next-to-build
	mt.State.CurrentBuild = model.BuildRecord{StartTime: time.Now()}
	f.assertNoTargetNextToBuild()
}

func TestCurrentBuildRecursive(t *testing.T) {
	f := newTestFixture(t)
	defer f.TearDown()

	m1 := manifestbuilder.New(f, "m1").WithK8sYAML(testyaml.SanchoYAML).WithResourceDeps("m2").Build()
	mt1 := store.NewManifestTarget(m1)
	f.st.UpsertManifestTarget(mt1)

	m2 := manifestbuilder.New(f, "m2").WithK8sYAML(testyaml.SanchoYAML).WithResourceDeps("m3").Build()
	mt2 := store.NewManifestTarget(m2)
	f.st.UpsertManifestTarget(mt2)

	m3 := manifestbuilder.New(f, "m3").WithK8sYAML(testyaml.SanchoYAML).Build()
	mt3 := store.NewManifestTarget(m3)
	f.st.UpsertManifestTarget(mt3)

	f.assertNextTargetToBuild(m3.Name)

	// If target is currently building, none of its deps should be next to build.
	mt3.State.CurrentBuild = model.BuildRecord{StartTime: time.Now()}

	f.assertNoTargetNextToBuild()
}

type testFixture struct {
	*tempdir.TempDirFixture
	t  *testing.T
	st *store.EngineState
}

func newTestFixture(t *testing.T) testFixture {
	f := tempdir.NewTempDirFixture(t)
	return testFixture{
		TempDirFixture: f,
		t:              t,
		st:             store.NewState(),
	}
}

func (f *testFixture) assertNextTargetToBuild(expected model.ManifestName) {
	next := NextTargetToBuild(*f.st)
	require.NotNil(f.t, next, "expected next target %s but got: nil", expected)
	actual := next.Manifest.Name
	assert.Equal(f.t, expected, actual, "expected next target to be %s but got %s", expected, actual)
}

func (f *testFixture) assertNoTargetNextToBuild() {
	next := NextTargetToBuild(*f.st)
	if next != nil {
		f.t.Fatalf("expected no next target to build, but got %s", next.Manifest.Name)
	}
}

func (f *testFixture) manifestNeedingCrashRebuild() *store.ManifestTarget {
	m := manifestbuilder.New(f, "needs-crash-rebuild").
		WithK8sYAML(testyaml.SanchoYAML).
		Build()
	mt := store.NewManifestTarget(m)
	mt.State.BuildHistory = []model.BuildRecord{
		model.BuildRecord{
			StartTime:  time.Now().Add(-5 * time.Second),
			FinishTime: time.Now(),
		},
	}
	mt.State.NeedsRebuildFromCrash = true
	return mt
}
