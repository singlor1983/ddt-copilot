package defs

type ProcessName string

const (
	ProcessTg3   ProcessName = "Tango3.exe"
	ProcessTgWeb ProcessName = "TangoWeb.exe"
)

const (
	VK_0      uintptr = 48
	VK_1      uintptr = 49
	VK_2      uintptr = 50
	VK_3      uintptr = 51
	VK_4      uintptr = 52
	VK_5      uintptr = 53
	VK_6      uintptr = 54
	VK_7      uintptr = 55
	VK_8      uintptr = 56
	VK_9      uintptr = 57
	VK_B      uintptr = 66
	VK_Q      uintptr = 81
	VK_E      uintptr = 69
	VK_T      uintptr = 84
	VK_Y      uintptr = 89
	VK_U      uintptr = 85
	VK_P      uintptr = 80
	VK_F      uintptr = 70
	VK_SPACE  uintptr = 32
	VK_LEFT   uintptr = 37
	VK_UP     uintptr = 38
	VK_RIGHT  uintptr = 39
	VK_DOWN   uintptr = 40
	VK_ESCAPE uintptr = 27
)

var vkMap = map[string]uintptr{
	"0":      VK_0,
	"1":      VK_1,
	"2":      VK_2,
	"3":      VK_3,
	"4":      VK_4,
	"5":      VK_5,
	"6":      VK_6,
	"7":      VK_7,
	"8":      VK_8,
	"9":      VK_9,
	"B":      VK_B,
	"Q":      VK_Q,
	"E":      VK_E,
	"T":      VK_T,
	"Y":      VK_Y,
	"U":      VK_U,
	"P":      VK_P,
	"F":      VK_F,
	"SPACE":  VK_SPACE,
	"LEFT":   VK_LEFT,
	"UP":     VK_UP,
	"RIGHT":  VK_RIGHT,
	"DOWN":   VK_DOWN,
	"ESCAPE": VK_ESCAPE,
}

func GetVkFromStr(key string) uintptr {
	v, ok := vkMap[key]
	if !ok {
		return 0
	}
	return v
}
