package global

var TextConfig map[string]string

func InitGlobalTextConfig(tconf map[string]string) {
	TextConfig = tconf
}

var AfterGameExp int

func InitAfterGameExp(i int) {
	AfterGameExp = i
}
