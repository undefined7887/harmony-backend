package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	idToken := ""

	claims, err := NewAuthService().Auth(context.TODO(), idToken)
	fmt.Println(err)

	if assert.NoError(t, err) {
		assert.True(t, claims.EmailVerified)

		fmt.Println("email:", claims.Email)
		fmt.Println("email_verified:", claims.EmailVerified)
		fmt.Println("picture:", claims.Picture)
	}
}
