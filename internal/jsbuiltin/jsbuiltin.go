package jsbuiltin

// JsModule JS模块
type JsModule struct {
	File string
	Code string
}

// GetJs 获得JS模块列表
func GetJs() []JsModule {
	return modules[:]
}
