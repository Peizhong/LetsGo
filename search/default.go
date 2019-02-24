package search

// 实现默认匹配器
type defaultMatcher struct{}

func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

// 声明函数的时候带有接收者。这个方法会和接收者类型绑定在一起
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}

// 调用
// dm := new(defaultMatch)
// dm.Search(..)

// 声明指针为接收者
func (m *defaultMatcher) Search2(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}

// 调用
// var dm defaultMatch
// dm.Search(..)

// 使用指针作为接收者声明的方法，只能在接口类型的值是一个指针的时候被调用。
// 使用值作为接收者声明的方法，在接口类型的值为值或者指针时，都可以被调用。
