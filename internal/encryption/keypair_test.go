package encryption

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeypairStorage(t *testing.T) {
	kp, err := NewKeyPair()
	assert.NoError(t, err)
	err = kp.StoreToLocation("/home/fred/Workspaces/go/src/configserver/samples/home/certs")
	assert.NoError(t, err)
}
