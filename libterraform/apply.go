package libterraform

import (
	"log"

	tf "github.com/hashicorp/terraform/terraform"
)

func (t *libterraform) Apply(path string, state *tf.State, stateCh chan<- tf.State) Task {
	return t.apply(path, false, state, stateCh)
}

func (t *libterraform) Destroy(path string, state *tf.State, stateCh chan<- tf.State) Task {
	return t.apply(path, true, state, stateCh)
}

func (t *libterraform) apply(path string, destroy bool, state *tf.State, stateCh chan<- tf.State) Task {
	var op string
	if destroy {
		op = "destroy"
	} else {
		op = "apply"
	}

	// Create context
	ctx, err := t.createContext(path, destroy, state, stateCh)
	if err != nil {
		return newTask(ctx).setError(err)
	}
	task := newTask(ctx)

	// context.Plan
	log.Println("[INFO] start plan", path)
	if _, err := ctx.Plan(); err != nil {
		log.Println("[INFO] error during plan:", err)
		return task.setError(err).done()
	}

	// Start the apply in a goroutine
	log.Println("[INFO] start ", op, path)
	go func() {
		state, err := ctx.Apply()
		log.Println("[INFO] done ", op, err)
		if stateCh != nil {
			close(stateCh)
		}
		task.setState(state).setError(err).done()
	}()

	return task
}
