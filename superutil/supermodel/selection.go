package supermodel

// Operator 操作符
type Operator string

// 操作符枚举
const (
	Equals  Operator = "equals"
	Belongs Operator = "belongs"
)

// TokenType 标记类型
type TokenType string

// 标记类型枚举
const (
	Table TokenType = "table"
	Const TokenType = "const"
)

// Token 标记
type Token struct {
	Type   TokenType `json:",omitempty"`
	Table  string    `json:",omitempty"`
	Key    string    `json:",omitempty"`
	Consts []string  `json:",omitempty"`
}

// Requirement 匹配条件
type Requirement struct {
	Not      bool     `json:",omitempty"`
	Operator Operator `json:",omitempty"`
	Left     Token    `json:",omitempty"`
	Right    Token    `json:",omitempty"`
}

// LabelSelector 标签选择器
type LabelSelector []Requirement
