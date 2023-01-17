package example

import (
	"fmt"
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// function is too short to patch 要 Mock 的函数太短了，加几行或者用 Debug 模式跑
// unexpected signal during runtime execution macOS 版本太低了（10.x / 11.x），使用该命令解决 go env -w CGO_ENABLED=0

func Test_isSuperUser(t *testing.T) {
	PatchConvey("", t, func() {
		Mock(getSuperUsers).To(func() []string {
			var superUsers []string
			for i := 1; i <= 3; i++ {
				superUsers = append(superUsers, fmt.Sprintf("super%d", i))
			}
			return superUsers
		}).Build()

		So(isSuperUser("super1"), ShouldBeTrue)
		So(isSuperUser("super4"), ShouldBeFalse)
	})
}
