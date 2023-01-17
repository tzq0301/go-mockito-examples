package example

import (
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_isSuperUser(t *testing.T) {
	PatchConvey("", t, func() {
		//Mock(refresh).To(func() {}).Build()
		Mock(refreshAsync).To(func() {}).Build()

		So(isSuperUser("super1"), ShouldBeTrue)
		So(isSuperUser("user1"), ShouldBeFalse)
	})
}
