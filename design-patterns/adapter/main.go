package main

import "fmt"

// 统一的支付接口

// PaymentProcessor 是一个统一的支付接口
type PaymentProcessor interface {
	Pay(amount float64)
}

// 旧支付系统

// OldPaymentSystem 是旧的支付系统
type OldPaymentSystem struct{}

func (ops *OldPaymentSystem) Pay(amount float64) {
	fmt.Printf("Processing payment of %.2f using old payment system\n", amount)
}

// 新支付系统

// NewPaymentSystem 是新的支付系统
type NewPaymentSystem struct{}

func (nps *NewPaymentSystem) MakePayment(amount float64) {
	fmt.Printf("Making payment of %.2f using new payment system\n", amount)
}

// 新支付系统的适配器

// NewPaymentAdapter 是新支付系统的适配器
type NewPaymentAdapter struct {
	// 内部持有新支付系统
	NewSystem *NewPaymentSystem
}

func (npa *NewPaymentAdapter) Pay(amount float64) {
	npa.NewSystem.MakePayment(amount)
}

func main() {
	// 声明支付接口
	var processor PaymentProcessor

	// 使用旧支付系统
	processor = &OldPaymentSystem{}
	processor.Pay(100)

	// 使用新支付系统
	newPayment := &NewPaymentSystem{}
	// 使用适配器模式
	processor = &NewPaymentAdapter{NewSystem: newPayment}
	processor.Pay(200)
}
