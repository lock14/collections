package treemap

import (
	"runtime"
	"sync/atomic"
	"testing"
)

type Dummy struct {
	v int
}

func TestTreeMap_MemoryLeak(t *testing.T) {
	tm := NewOrdered[int, *Dummy]()

	var collected int32

	// Create a closure that creates the dummy object so that there are no
	// local variables holding references to it in the main test function.
	func() {
		d := &Dummy{v: 42}
		runtime.SetFinalizer(d, func(obj *Dummy) {
			atomic.AddInt32(&collected, 1)
		})

		tm.Put(1, d)
	}()

	// Remove the element. If there's a memory leak (e.g., trailing element in slice
	// not being zeroed out), the pointer will remain in the underlying slice capacity,
	// and the object won't be collected.
	tm.Remove(1)

	// Force GC
	runtime.GC()
	runtime.GC() // Sometimes needs two passes

	// Check if it was collected
	if atomic.LoadInt32(&collected) == 0 {
		t.Error("Memory leak detected: removed element was not garbage collected")
	}
}
