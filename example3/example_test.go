package example

import (
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_isSuperUser(t *testing.T) {
	PatchConvey("", t, func() {
		Mock(getSuperUsers).Return([]string{"super1", "super2", "super3"}).Build()

		So(isSuperUser("super1"), ShouldBeTrue)
		So(isSuperUser("super4"), ShouldBeFalse)
	})
}
