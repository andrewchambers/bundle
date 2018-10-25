# bundle - structured concurrency

This go package provides structured concurrency and garbage collection of goroutines.

## Examples

Structured cleanup

```
b := bundle.New(context.Background())
// Cancellation occurs when you call b.Close or b.Cancel
defer b.Cancel()

b.Go(func(ctx context.Context) {
    // Do some work...
})

b.Go(func(ctx context.Context) {
    // Do some more work...
})

// After b.Wait() returns we know all goroutines have exited.
b.Wait()
```

Garbage collection of goroutines:

```

b := bundle.New(context.Background()).Go

b.Go(func(ctx context.Context) {
    <- ctx.Done()
    fmt.Println("Bundle was garbage collected\n")
})

// Remove reference to b so it is garbage collected.
b = nil

// ... If your program allocates enough memory to trigger the
// garbage collector, eventually you will see the printed output. 
```


## How it works

With some indirection and runtime.SetFinalizer, we can cancel our context in a safe way. sync.Waitgroup
allows us to block until all workers have exited.

## Is this worth a library?

Maybe not, I just found myself using sync.WaitGroup a lot, and this just wraps up how I use it.
Let me know in the issues if this has been done before, or in a better way.

# Use cases

Using the garbage collector to close contexts might mean is useful for things like infinite lazy
lists, be creative.

## Influences

- https://vorpus.org/blog/notes-on-structured-concurrency-or-go-statement-considered-harmful/
- http://libdill.org/