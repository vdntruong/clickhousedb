package concurrency

import "context"

// OrDone wraps a channel with context-aware cancellation.
// It allows for clean shutdown of channel operations when the context is cancelled.
// Parameters:
//   - ctx: Context for cancellation control
//   - ch: Source channel to wrap with cancellation capability
//
// Returns:
//   - <-chan T: Output channel that respects context cancellation
//
// Generic type T allows this function to work with channels of any type.
func OrDone[T any](ctx context.Context, ch <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)

		if ch == nil {
			return
		}

		for {
			select {
			case <-ctx.Done():
				return

			case v, ok := <-ch:
				if !ok {
					return
				}

				select {
				case <-ctx.Done():
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}
