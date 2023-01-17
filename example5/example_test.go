package example

import (
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestHandler_PingPong(t *testing.T) {
	PatchConvey("", t, func() {
		h := Handler{service: NewService()}

		//Mock(Service.PingPong).To(func(ping string) (pong string) {
		//	pong = ping
		//	return
		//}).Build()

		Mock(GetMethod(h.service, "PingPong")).To(func(ping string) (pong string) {
			pong = ping
			return
		}).Build()

		So(h.PingPong("111"), ShouldEqual, "111")
	})
}
