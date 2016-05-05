package libterraform

import (
	"sync"

	tf "github.com/hashicorp/terraform/terraform"
)

// Task represents a running terraform operation
type Task interface {
	Cancel()

	Error() error

	Wait() (*tf.State, error)

	State() *tf.State
}

type task struct {
	sync.Mutex

	ctx    Context
	doneCh chan struct{}
	state  *tf.State
	err    error
}

func newTask(ctx Context) *task {
	return &task{
		doneCh: make(chan struct{}),
		ctx:    ctx,
	}
}

func (t *task) setError(err error) *task {
	t.Lock()
	defer t.Unlock()

	t.err = err
	return t
}

func (t *task) setState(state *tf.State) *task {
	t.Lock()
	defer t.Unlock()

	t.state = state
	return t
}

func (t *task) done() *task {
	t.Lock()
	defer t.Unlock()

	close(t.doneCh)
	t.doneCh = nil
	return t
}

func (t *task) State() *tf.State {
	t.Lock()
	defer t.Unlock()

	return t.state
}

func (t *task) Error() error {
	t.Lock()
	defer t.Unlock()

	return t.err
}

func (t *task) Cancel() {
	t.Lock()
	ctx := t.ctx
	defer t.Unlock()

	ctx.Stop()
}

func (t *task) Wait() (*tf.State, error) {
	t.Lock()
	done := t.doneCh
	t.Unlock()

	if done != nil {
		<-done
	}
	return t.state, t.err
}
