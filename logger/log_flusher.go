package logger

import (
	"sync"
)

/*
结构体中的嵌入
嵌入式类型（例如 mutex）应位于结构体内的字段列表的顶部，并且必须有一个空行将嵌入式字段与常规字段分隔开。

内嵌应该提供切实的好处，比如以语义上合适的方式添加或增强功能。 它应该在对用户没有任何不利影响的情况下使用。（另请参见：避免在公共结构中嵌入类型）。

例外：即使在未导出类型中，Mutex 也不应该作为内嵌字段。另请参见：零值 Mutex 是有效的。

嵌入 不应该:
	纯粹是为了美观或方便。
	使外部类型更难构造或使用。
	影响外部类型的零值。如果外部类型有一个有用的零值，则在嵌入内部类型之后应该仍然有一个有用的零值。
	作为嵌入内部类型的副作用，从外部类型公开不相关的函数或字段。
	公开未导出的类型。
	影响外部类型的复制形式。
	更改外部类型的 API 或类型语义。
	嵌入内部类型的非规范形式。
	公开外部类型的实现详细信息。
	允许用户观察或控制类型内部。
	通过包装的方式改变内部函数的一般行为，这种包装方式会给用户带来一些意料之外情况。
	简单地说，有意识地和有目的地嵌入。一种很好的测试体验是， "是否所有这些导出的内部方法/字段都将直接添加到外部类型" 如果答案是some或no，不要嵌入内部类型 - 而是使用字段。
*/

var flusherQueue = &struct {
	sync.RWMutex

	flushers []*fileWriter
}{}

func addFlusher(lf *fileWriter) {
	flusherQueue.Lock()
	flusherQueue.flushers = append(flusherQueue.flushers, lf)
	flusherQueue.Unlock()
}

func flushAllLog() {
	flusherQueue.RLock()
	for _, v := range flusherQueue.flushers {
		_ = v.Flush()
	}
	flusherQueue.RUnlock()
}
