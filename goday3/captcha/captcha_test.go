package captcha_test

import (
	"fmt"
	"testing"

	"github.com/pallat/todos/captcha"

	"github.com/stretchr/testify/assert"
)

func TestCaptchaPattern1(t *testing.T) {

	operands := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for givenOperand, want := range operands {
		t.Run(fmt.Sprintf("right operand %d", givenOperand), func(t *testing.T) {
			givenPattern := 1
			lo := 1
			op := 1
			ro := givenOperand
			want := fmt.Sprintf("1 + %s", want)

			cc := captcha.New(givenPattern, lo, op, ro)

			get := cc.String()
			assert.Equal(t, want, get, fmt.Sprintf("given %d %d %d %d want %q but get %q", givenPattern, ro, op, lo, want, get))
		})
	}

	for givenOperand := 1; givenOperand <= 9; givenOperand++ {
		t.Run(fmt.Sprintf("right operand %d", givenOperand), func(t *testing.T) {
			givenPattern := 1
			lo := givenOperand
			op := 1
			ro := 1
			want := fmt.Sprintf("%d + one", lo)

			cc := captcha.New(givenPattern, lo, op, ro)

			get := cc.String()
			assert.Equal(t, want, get, fmt.Sprintf("given %d %d %d %d want %q but get %q", givenPattern, ro, op, lo, want, get))
		})
	}
}

func TestCaptchaPattern2(t *testing.T) {

	operands := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

	for givenOperand, want := range operands {
		t.Run(fmt.Sprintf("left operand %d", givenOperand), func(t *testing.T) {
			givenPattern := 2
			lo := givenOperand
			op := 1
			ro := 1
			want := fmt.Sprintf("%s + 1", want)

			cc := captcha.New(givenPattern, lo, op, ro)

			get := cc.String()
			assert.Equal(t, want, get, fmt.Sprintf("given %d %d %d %d want %q but get %q", givenPattern, ro, op, lo, want, get))
		})
	}
}
