# Go Mockito Examples

提供 examples 以对 [ByteDance Go Mockey](https://github.com/bytedance/mockey) 的简单使用，**_若遇到复杂情况或与 examples 表现不一的情况或需要使用进阶功能，请查询官网_**

可以 clone 该仓库，对一些 Mock 的方式做调整来进行实验

```shell
$ go get github.com/bytedance/mockey@latest
```

```go
import (
	. "github.com/bytedance/mockey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)
```

## Examples

### Example 1 - Mock 变量

假设有一个变量 `superUsers`，其值是通过 RPC 来获取的

我们要测 `isSuperUser` 函数的话，应该屏蔽掉 superUsers 这个外部依赖，可以使用 `MockValue` 函数，传入变量的地址即可

```go
// example1/example.go

var superUsers []string = getSuperUsers()

func getSuperUsers() []string {
	// Get From RPC
	return nil
}

func isSuperUser(user string) bool {
	for _, superUser := range superUsers {
		if superUser == user {
			return true
		}
	}
	return false
}
```

```go
// example1/example_test.go

func Test_isSuperUser(t *testing.T) {
	PatchConvey("", t, func() {
		MockValue(&superUsers).To([]string{"super1", "super2", "super3"})

		So(isSuperUser("super1"), ShouldBeTrue)
		So(isSuperUser("super4"), ShouldBeFalse)
	})
}
```

### Example 2 - Mock 函数行为

假设 `getSuperUsers` 函数将通过 RPC 来获取数据，那么我们要测 `isSuperUser` 函数时必须 Mock 掉这个外部依赖

可以使用 `Mock(..).To(..).Build()` 来 Mock 函数的行为

* `To(..)` 中传的参数是跟你要 Mock 的函数相同签名的函数（参考示例 `func() []string`）
* 记得调用 `.Build(..)`

```go
// example2/example.go

func getSuperUsers() []string {
	// Get From RPC
	var superUsers []string
	for i := 1; i <= 3; i++ {
		superUsers = append(superUsers, fmt.Sprintf("super%d from rpc", i))
	}
	return superUsers
}

func isSuperUser(user string) bool {
	for _, superUser := range getSuperUsers() {
		if superUser == user {
			return true
		}
	}
	return false
}
```

```go
// example2/example_test.go

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
```

### Example 3 - 直接 Mock 函数返回值

需求同 Example 2，需要对 `isSuperUser` 函数进行单测，并 Mock 掉 `getSuperUsers` 函数

可以直接 Mock `getSuperUsers` 函数的返回值，更加方便快捷

```go
// example3/example_test.go

func Test_isSuperUser(t *testing.T) {
    PatchConvey("", t, func() {
        Mock(getSuperUsers).Return([]string{"super1", "super2", "super3"}).Build()
        
        So(isSuperUser("super1"), ShouldBeTrue)
        So(isSuperUser("super4"), ShouldBeFalse)
    })
}
```

### Example 4 - 破除协程的干扰

在项目过程中，我曾尝试 Mock `isSuperUser` 函数中使用 `go` 来调用的 `refresh` 函数，但在 Mock 的时候出现了一些问题（虽然该 example 不一定能复现）

* 可以通过控制 `Patch()` 与 `UnPatch()` 方法解决，但不适用于我当时的代码结构
* 可以通过 `time.Sleep(..)` 来对齐协程执行顺序，但这很不优雅

```go
// example4/example.go

func refresh() {
	log.Println("refreshing...")
	log.Println("refreshing...")
}

func isSuperUser(user string) bool {
	go refresh()

	return strings.HasPrefix(user, "super")
}
```

一个比较取巧的方法是对 `go refresh()` 语句进行封装，然后 Mock 掉封装函数的行为

```go
// example4/example.go

func refresh() {
	log.Println("refreshing...")
	log.Println("refreshing...")
}

func refreshAsync() {
	go refresh()
}

func isSuperUser(user string) bool {
	refreshAsync()

	return strings.HasPrefix(user, "super")
}
```

```go
// example4/example_test.go

func Test_isSuperUser(t *testing.T) {
	PatchConvey("", t, func() {
		//Mock(refresh).To(func() {}).Build()
		Mock(refreshAsync).To(func() {}).Build()

		So(isSuperUser("super1"), ShouldBeTrue)
		So(isSuperUser("user1"), ShouldBeFalse)
	})
}
```

### Example 5 - Mock 接口的私有实现结构体的方法

在实际开发中，经常对 Handler、Service 这类进行接口抽象，只暴露接口签名，将实例化该接口私有实现结构体的步骤封装为一个 `NewXxxx` 方法

但是 Mockito **没法直接 Mock 一个 interface 的方法**

如下代码，我们要测 `func (h Handler) PingPong(ping string) (pong string)` 方法，但是**其依赖的 `func (e serviceImpl) PingPong(ping string) (pong string)` 方法涉及了 RPC 调用，我们要屏蔽该外部依赖**：

```go
// example5/example.go

type Service interface {
	PingPong(ping string) (pong string)
}

func NewService() Service {
	return &serviceImpl{}
}

type serviceImpl struct{}

func (e serviceImpl) PingPong(ping string) (pong string) {
	log.Println("Begin Do some RPC Call !!!")
	pong = ping
	log.Println("End   Do some RPC Call !!!")
	return
}

// --------------------------------------------------------

type Handler struct {
	service Service
}

// 测这个方法！！！！！！！！！！！！！！！！！！！！！！！
func (h Handler) PingPong(ping string) (pong string) {
	pong = h.service.PingPong(ping)
	return
}
```

如果直接 `Mock(Service.PingPong).To(..).Build()`，则会发现并没有 Mock 成功，`service.PingPong(..)` 还是走到了原函数（会出现 RPC Call 的信息）

```go
// example5/example_test.go

func TestHandler_PingPong(t *testing.T) {
	PatchConvey("", t, func() {
		h := Handler{service: NewService()}

		Mock(Service.PingPong).To(func(ping string) (pong string) {
			pong = ping
			return
		}).Build()

		So(h.PingPong("111"), ShouldEqual, "111")
	})
}
```

但是我们又不能 Mock `Service` 的实现类 `serviceImpl`，因为其是私有的，包外不可访问的（在本例中由于在同一个 package 下可以，但实际情况中会遇到要使用另一个 package 中的 Service 但情况）

这时，我们可以使用 Mockito 提供的工具方法 `func GetMethod(instance interface{}, methodName string) interface{}`，先获得 `Service` 实现的实例，通过 `GetMethod` 拿到我们要 Mock 的方法

```go
// example5/example_test.go

func TestHandler_PingPong(t *testing.T) {
    PatchConvey("", t, func() {
        h := Handler{service: NewService())
        
        Mock(GetMethod(h.service, "PingPong")).To(func(ping string) (pong string) {
            pong = ping
            return
        }).Build()
        
        So(h.PingPong("111"), ShouldEqual, "111")
    })
}
```

## FAQ

### Mock  函数后仍走入了原函数

* 未禁用内联或者编译优化，可以用 Debug 模式跑试试
* Mock 函数后没调用 `Build()`
* Mock 函数的签名不一致，例如：
  * `Mock(f).To(...).Build()`
  * `Mock(A.f).To(...).Build()`
  * `Mock((*A).f).To(...).Build()`
* 协程相关问题

### function is too short to patch 

要 Mock 的函数太短了，可能只有一行，导致编译后机器码太短，一般两行及以上不会有这个问题

可以进行**内联**或者**用 Debug 模式**跑测试

### unexpected signal during runtime execution

Mac OS 10.x / 11.x 版本太低，使用以下命令解决

```shell
go env -w CGO_ENABLED=0
```