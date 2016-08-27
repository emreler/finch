package channel

import "testing"

func TestNotify(t *testing.T) {
	h := &HttpChannel{}

	h.Notify("http://requestb.in/10fumyy1", "testing")

}
