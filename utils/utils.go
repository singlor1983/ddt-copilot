package utils

import (
	"github.com/lxn/win"
	"time"
)

type Empty struct{}

type Set[T comparable] map[T]Empty

func (s Set[T]) Add(item T) {
	s[item] = Empty{}
}

func (s Set[T]) Del(item T) {
	delete(s, item)
}

func (s Set[T]) Contain(item T) bool {
	_, exists := s[item]
	return exists
}

func (s Set[T]) Range(fn func(T)) {
	for t := range s {
		fn(t)
	}
}

func NewSet[T comparable]() Set[T] {
	return make(Set[T])
}

func LeftClick(hWnd win.HWND, x int, y int, duration time.Duration) {
	lParam := uintptr(y<<16 | x)
	win.PostMessage(hWnd, win.WM_LBUTTONDOWN, win.MK_LBUTTON, lParam)
	time.Sleep(duration)
	win.PostMessage(hWnd, win.WM_LBUTTONUP, 0, lParam)
}

func KeyBoard(hWnd win.HWND, key uintptr, duration time.Duration) {
	win.PostMessage(hWnd, win.WM_KEYDOWN, key, 0)
	time.Sleep(duration)
	win.PostMessage(hWnd, win.WM_KEYUP, key, 0)
}
