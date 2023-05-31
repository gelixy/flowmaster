package blocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoSimpleBlocks(t *testing.T) {
	action1 := NewPath("step 1", nil)
	action2 := NewPath("step 2", nil)

	action1.LinkDown(action2)

	action1.In <- "hello world"

	message := <-action2.Out

	assert.Equal(t, "hello world", message)
}

func TestTwoSimpleBlocksWithConverter(t *testing.T) {
	action1 := NewPath("step 1", func(message any) any {
		return message.(string) + " :: step 1"
	})
	action2 := NewPath("step 2", func(message any) any {
		return message.(string) + " :: step 2"
	})

	action1.LinkDown(action2)

	action1.In <- "hello world"

	message := <-action2.Out

	assert.Equal(t, "hello world :: step 1", message)
}
