package go_yapi

import (
	"testing"
)

func TestJsonParam(t *testing.T) {
	do := func(in interface{}, need string) {
		if esc := jsonParam(in); esc != need {
			t.Errorf("'%s' escaped as '%s' but need '%s'", in, esc, need)
		}
	}

	do(`quotes "word"`, `"quotes \"word\""`)
	do(`backslash \ \string\`, `"backslash \\ \\string\\"`)
	do("one line\r\ntwo line", `"one line\r\ntwo line"`)
	do(5, "5")
	do(1.2, "1.2")
	do(nil, "null")
}
