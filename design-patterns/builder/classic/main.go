package main

import "fmt"

// House 代表一个由多个部分构建的复杂对象
type House struct {
	// 地基
	Foundation string
	// 墙壁
	Walls string
	// 屋顶
	Roof string
}

// NewHouse 一个普通的 House 构造函数
func NewHouse(foundation, walls, roof string) *House {
	return &House{
		Foundation: foundation,
		Walls:      walls,
		Roof:       roof,
	}
}

// Builder 用于构建 House 对象的接口
type Builder interface {
	BuildFoundation()
	BuildWalls()
	BuildRoof()
	GetResult() *House
}

// ConcreteBuilder Builder 接口的具体实现，用于构建具体的 House
type ConcreteBuilder struct {
	house *House
}

// NewConcreteBuilder 创建一个新的 ConcreteBuilder 实例
func NewConcreteBuilder() *ConcreteBuilder {
	return &ConcreteBuilder{house: &House{}}
}

// BuildFoundation 构建地基
func (b *ConcreteBuilder) BuildFoundation() {
	b.house.Foundation = "Concrete Foundation"
}

// BuildWalls 构建墙壁
func (b *ConcreteBuilder) BuildWalls() {
	b.house.Walls = "Wooden Walls"
}

// BuildRoof 构建屋顶
func (b *ConcreteBuilder) BuildRoof() {
	b.house.Roof = "Shingle Roof"
}

// GetResult 返回构建完成的 House 对象
func (b *ConcreteBuilder) GetResult() *House {
	return b.house
}

// Director 用于控制构建过程的指挥者
type Director struct {
	builder Builder
}

// NewDirector 创建一个新的 Director 实例
func NewDirector(builder Builder) *Director {
	return &Director{builder: builder}
}

// Construct 构建 House 的方法
func (d *Director) Construct() {
	d.builder.BuildFoundation()
	d.builder.BuildWalls()
	d.builder.BuildRoof()
}

func main() {
	// 普通方式
	{
		house := NewHouse("Concrete Foundation", "Wooden Walls", "Shingle Roof")
		fmt.Printf("%+v\n", house)
	}

	// builder 模式
	{
		// 创建具体的 Builder
		builder := NewConcreteBuilder()

		// 创建 Director 并传入具体的 Builder
		director := NewDirector(builder)

		// 通过 Director 控制构建过程
		director.Construct()

		// 获取构建的最终产品
		house := builder.GetResult()

		// 输出构建好的 House
		fmt.Printf("%+v\n", house)
	}
}
