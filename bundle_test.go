package bundle

import (
	"context"
	"fmt"
	"testing"
)

func TestBundle(t *testing.T) {
	b := New(context.Background())

	c1 := make(chan struct{}, 1)
	c2 := make(chan struct{}, 1)

	b.Go(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			c1 <- struct{}{}
			return
		}

	})

	b.Go(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			c2 <- struct{}{}
			return
		}
	})

	b.Close()

	select {
	case <-c1:
	default:
		t.FailNow()
	}

	select {
	case <-c2:
	default:
		t.FailNow()
	}
}

func TestBundleAutomaticCancel(t *testing.T) {

	c := make(chan struct{}, 1)

	func(c chan struct{}) {
		b := New(context.Background())

		b.Go(func(ctx context.Context) {
			select {
			case <-ctx.Done():
				c <- struct{}{}
				return
			}
		})

	}(c)

	func() {
		i := 0
		for {
			// Allocate like crazy
			_ = fmt.Sprintf("%d", i)
			i += 1

			select {
			case <-c:
				// Our context was automatically cancelled, test passed.
				t.Logf("cancelled after %d iterations", i)
				return
			default:
			}
		}
	}()

}
