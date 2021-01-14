package captcha

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Captcha struct {
	pattern      int
	leftOperand  int
	operator     int
	rightOperand int
}

func New(p, l, o, r int) Captcha {
	return Captcha{p, l, o, r}
}

var operandString = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
var operators = []string{"", "+", "-", "*"}

func (cc Captcha) String() string {
	if cc.pattern == 1 {
		return fmt.Sprintf("%d %s %s", cc.leftOperand, operators[cc.operator], operandString[cc.rightOperand])
	}
	return fmt.Sprintf("%s %s %d", operandString[cc.leftOperand], operators[cc.operator], cc.rightOperand)
}

var src = rand.NewSource(time.Now().UnixNano())
var rnd = rand.New(src)
var store = map[string]int{}

func KeyQuestion() (string, string) {
	pattern, leftOperand, operator, rightOperand := rnd.Intn(2)+1, rand.Intn(9)+1, rand.Intn(3)+1, rand.Intn(9)+1
	ans := 0
	switch operator {
	case 1:
		ans = leftOperand + rightOperand
	case 2:
		ans = leftOperand - rightOperand
	case 3:
		ans = leftOperand * rightOperand
	}
	cc := New(pattern, leftOperand, operator, rightOperand)
	key := uuid.New().String()
	store[key] = ans
	return key, cc.String()
}

var mux sync.Mutex

func Answer(key string, ans int) bool {
	mux.Lock()
	defer mux.Unlock()

	if v, ok := store[key]; ok {
		delete(store, key)
		return v == ans
	}

	log.Printf("not found % key in store\n", key)
	return false
}
