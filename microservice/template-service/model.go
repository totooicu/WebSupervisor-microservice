package main

// EchoParameter echo服务参数
type EchoParameter struct {
	Message string `json:"message"`
}

// AddParameter 加法服务参数
type AddParameter struct {
	Num1 float64 `json:"num1"`
	Num2 float64 `json:"num2"`
}

// SubtractParameter 减法服务参数
type SubtractParameter struct {
	Num1 float64 `json:"num1"`
	Num2 float64 `json:"num2"`
}

// MultiplyParameter 乘法服务参数
type MultiplyParameter struct {
	Num1 float64 `json:"num1"`
	Num2 float64 `json:"num2"`
}

// DivideParameter 除法服务参数
type DivideParameter struct {
	Num1 float64 `json:"num1"`
	Num2 float64 `json:"num2"`
}