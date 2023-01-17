package example

import (
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_isSuperUser(t *testing.T) {
	PatchConvey("", t, func() {
		MockValue(&superUsers).To([]string{"super1", "super2", "super3"})

		So(isSuperUser("super1"), ShouldBeTrue)
		So(isSuperUser("super4"), ShouldBeFalse)
	})
}
