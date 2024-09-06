package mapper

// Option 用于设置 Config 结构体的字段。
type Option func(*Setting)

type Setting struct {
	EnabledTypeChecking      bool
	EnabledMapperStructField bool
	EnabledAutoTypeConvert   bool
	EnabledMapperTag         bool
	EnabledJsonTag           bool
	EnabledCustomTag         bool
	CustomTagName            string

	// in the version < 0.7.8, we use field name as the key when mapping structs if field tag is "-"
	// from 0.7.8, we add switch enableIgnoreFieldTag which is false in default
	// if caller enable this flag, the field will be ignored in the mapping process
	EnableFieldIgnoreTag bool
}

// getDefaultSetting return default mapper setting
// Use default value:
//
//	EnabledTypeChecking:      false,
//	EnabledMapperStructField: true,
//	EnabledAutoTypeConvert:   true,
//	EnabledMapperTag:         true,
//	EnabledJsonTag:           true,
//	EnabledCustomTag:         false,
//	EnableFieldIgnoreTag:     false
func getDefaultSetting() *Setting {
	return &Setting{
		EnabledTypeChecking:      false,
		EnabledMapperStructField: true,
		EnabledAutoTypeConvert:   true,
		EnabledMapperTag:         true,
		EnabledJsonTag:           true,
		EnabledCustomTag:         false,
		EnableFieldIgnoreTag:     false, // 保留老版本默认行为：对于tag = “-”的字段使用FieldName
	}
}

// NewSetting create new setting with multi option
func NewSetting(opts ...Option) *Setting {
	cfg := getDefaultSetting()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// CTypeChecking set EnabledTypeChecking value
//
// Default value: false
func CTypeChecking(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnabledTypeChecking = isEnabled
	}
}

// CMapperTag set EnabledMapperTag value
//
// Default value: true
func CMapperTag(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnabledMapperTag = isEnabled
	}
}

func CJsonTag(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnabledJsonTag = isEnabled
	}
}

func CAutoTypeConvert(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnabledAutoTypeConvert = isEnabled
	}
}

func CMapperStructField(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnabledMapperStructField = isEnabled
	}
}

func CCustomTag(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnabledCustomTag = isEnabled
	}
}

func CFieldIgnoreTag(isEnabled bool) Option {
	return func(c *Setting) {
		c.EnableFieldIgnoreTag = isEnabled
	}
}

func CCustomTagName(tagName string) Option {
	return func(c *Setting) {
		c.CustomTagName = tagName
	}
}
