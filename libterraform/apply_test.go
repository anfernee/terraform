package libterraform

import (
	"fmt"
	"testing"

	tf "github.com/hashicorp/terraform/terraform"
)

type mockContext struct {
	destroy     bool
	applyCalled int
	planCalled  int
	stopCalled  int
}

func (c *mockContext) Apply() (*tf.State, error) {
	c.applyCalled += 1
	return nil, nil
}

func (c *mockContext) Plan() (*tf.Plan, error) {
	c.planCalled += 1
	return nil, nil
}

func (c *mockContext) Stop() {
	c.stopCalled += 1
}

func createMockContext(path string, destroy bool, state *tf.State, stateCh chan<- tf.State) (Context, error) {
	return &mockContext{destroy: destroy}, nil
}

func createTestTerraform() *libterraform {
	return &libterraform{
		createContext: createMockContext,
	}
}

func TestApply(t *testing.T) {
	tt := createTestTerraform()

	tsk := tt.Apply("/path/to/blueprint", nil, nil)
	tsk.Wait()
	ctx := tsk.(*task).ctx.(*mockContext)
	if ctx.applyCalled != 1 {
		t.Errorf("expect apply to be called 1 time; got %d", ctx.applyCalled)
	}
	if ctx.planCalled != 1 {
		t.Errorf("expect plan to be called 1 time; got %d", ctx.planCalled)
	}
	if ctx.stopCalled != 0 {
		t.Errorf("expect stop to be called 0 time; got %d", ctx.stopCalled)
	}
	if ctx.destroy {
		t.Errorf("expect destroy to be false; got true")
	}
}

func TestDestroy(t *testing.T) {
	tt := createTestTerraform()

	tsk := tt.Destroy("/path/to/blueprint", nil, nil)
	tsk.Wait()
	ctx := tsk.(*task).ctx.(*mockContext)
	if ctx.applyCalled != 1 {
		t.Errorf("expect apply to be called 1 time; got %d", ctx.applyCalled)
	}
	if ctx.planCalled != 1 {
		t.Errorf("expect plan to be called 1 time; got %d", ctx.planCalled)
	}
	if ctx.stopCalled != 0 {
		t.Errorf("expect stop to be called 0 time; got %d", ctx.stopCalled)
	}
	if !ctx.destroy {
		t.Errorf("expect destroy to be true; got false")
	}
}

// sample
func testApplySample(t *testing.T) {
	tt := NewTerraform()
	stateCh := make(chan tf.State)

	go func() {
		for state := range stateCh {
			fmt.Println("xxxxxxxxxxxxxxx", state)
		}
	}()

	task := tt.Apply("./test_files/example/", nil, stateCh)
	state, err := task.Wait()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(state)

	// Destroy
	stateCh = make(chan tf.State)
	go func() {
		for state := range stateCh {
			fmt.Println("yyyyyyyyyyyyyyyy", state)
		}
	}()

	task = tt.Destroy("./test_files/example/", state, stateCh)
	state, err = task.Wait()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(state)
}
