package conc

import (
	"context"
	"runtime"
	"sync"
)

// A bundle is a set of goroutines that
// run together under a shared context.
//
// When the bundle is garbage collected,
// it is cancelled automatically.
type Bundle struct {
	// These layers of indirection are needed
	// to make the bundle copyable.
	selfDestruct *selfDestruct
	bundle       *bundle
}

type selfDestruct struct {
	bundle *bundle
}

func newSelfDestruct(b *bundle) *selfDestruct {
	s := &selfDestruct{
		bundle: b,
	}

	runtime.SetFinalizer(s, func(s *selfDestruct) {
		s.bundle.Cancel()
	})

	return s
}

type bundle struct {
	wg              sync.WaitGroup
	bundleCtx       context.Context
	cancelOnce      sync.Once
	cancelBundleCtx func()
}

func New(parentContext context.Context) *Bundle {
	bundleCtx, cancelBundleCtx := context.WithCancel(parentContext)

	b := &bundle{
		wg:              sync.WaitGroup{},
		bundleCtx:       bundleCtx,
		cancelBundleCtx: cancelBundleCtx,
	}

	return &Bundle{
		selfDestruct: newSelfDestruct(b),
		bundle:       b,
	}
}

func (b *bundle) Go(task func(ctx context.Context)) {
	select {
	case <-b.bundleCtx.Done():
		// Do not start a new task into a cancelled
		// context.
		return
	default:
	}

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		task(b.bundleCtx)
	}()
}

func (b *bundle) Cancel() {
	b.cancelOnce.Do(func() {
		b.cancelBundleCtx()
	})
}

func (b *bundle) Wait() {
	b.wg.Wait()
}

// Go starts a function within the bundles context.
func (b *Bundle) Go(task func(ctx context.Context)) {
	b.bundle.Go(task)
}

// Wait waits for all goroutines in the bundle to exit.
func (b *Bundle) Wait() {
	b.bundle.Wait()
}

// Cancel cancels all tasks in the bundle.
func (b *Bundle) Cancel() {
	b.bundle.Cancel()
}

// Close cancels the bundle
// and waits until they are completed.
func (b *Bundle) Close() {
	b.bundle.Cancel()
	b.bundle.Wait()
}
