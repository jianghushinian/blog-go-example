package main

import "fmt"

// Car 代表一个汽车对象
type Car struct {
	// 品牌
	Brand string
	// 型号
	Model string
	// 颜色
	Color string
	// 发动机类型
	Engine string
}

// CarBuilder 用于构建 Car 对象的构建器
type CarBuilder struct {
	car Car
}

// NewCarBuilder 创建一个新的 CarBuilder 实例
func NewCarBuilder() *CarBuilder {
	return &CarBuilder{car: Car{}}
}

// SetBrand 设置汽车的品牌
func (b *CarBuilder) SetBrand(brand string) *CarBuilder {
	b.car.Brand = brand
	return b
}

// SetModel 设置汽车的型号
func (b *CarBuilder) SetModel(model string) *CarBuilder {
	b.car.Model = model
	return b
}

// SetColor 设置汽车的颜色
func (b *CarBuilder) SetColor(color string) *CarBuilder {
	b.car.Color = color
	return b
}

// SetEngine 设置汽车的发动机类型
func (b *CarBuilder) SetEngine(engine string) *CarBuilder {
	b.car.Engine = engine
	return b
}

// Build 构建并返回最终的 Car 对象
func (b *CarBuilder) Build() Car {
	return b.car
}

func main() {
	// 使用 CarBuilder 构建一个 Car 对象
	car := NewCarBuilder().
		SetBrand("Tesla").
		SetModel("Model S").
		SetColor("Red").
		SetEngine("Electric").
		Build()

	fmt.Printf("Car: %+v\n", car)
}
