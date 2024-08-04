package main

import (
	"fmt"
)

// 生产实践：多云管理平台

type ProviderType string

const (
	ProviderTypeAliyun  ProviderType = "aliyun"
	ProviderTypeTencent ProviderType = "tencent"
	// ...
)

// RunInstanceRequest 创建云主机参数
type RunInstanceRequest struct {
	// ...
}

// RunInstanceResponse 创建云主机返回结果
type RunInstanceResponse struct {
	// ...
}

// Provider 定义云厂商统一接口
type Provider interface {
	// Type 返回 Provider 类型
	Type() ProviderType
	// RunInstance 创建云主机
	RunInstance(r *RunInstanceRequest) (*RunInstanceResponse, error)
	// ...
}

func NewProvider(typ ProviderType) Provider {
	switch typ {
	case ProviderTypeAliyun:
		return NewAliCloudProvider(typ)
	case ProviderTypeTencent:
		return NewTencentCloudProvider(typ)
	default:
		panic("unknown provider")
	}
	return nil
}

// AliCloudProvider 阿里云 Provider
type AliCloudProvider struct {
	typ ProviderType
	// 包装 alibaba-cloud-sdk-go
}

func NewAliCloudProvider(typ ProviderType) *AliCloudProvider {
	return &AliCloudProvider{
		typ: typ,
		// ...
	}
}

func (a AliCloudProvider) Type() ProviderType {
	return a.typ
}

// RunInstance https://help.aliyun.com/zh/ecs/developer-reference/api-ecs-2014-05-26-runinstances
func (a AliCloudProvider) RunInstance(r *RunInstanceRequest) (*RunInstanceResponse, error) {
	panic("implement me")
}

// TencentCloudProvider 腾讯云 Provider
type TencentCloudProvider struct {
	typ ProviderType
	// 包装 tencentcloud-sdk-go
}

func NewTencentCloudProvider(typ ProviderType) *TencentCloudProvider {
	return &TencentCloudProvider{
		typ: typ,
		// ...
	}
}

func (t TencentCloudProvider) Type() ProviderType {
	return t.typ
}

// RunInstance https://cloud.tencent.com/document/api/213/15730
func (t TencentCloudProvider) RunInstance(r *RunInstanceRequest) (*RunInstanceResponse, error) {
	panic("implement me")
}

func main() {
	var p Provider

	// 阿里云
	p = NewProvider(ProviderTypeAliyun)
	resp, err := p.RunInstance(&RunInstanceRequest{})
	fmt.Println(resp, err)

	// 腾讯云
	p = NewProvider(ProviderTypeTencent)
	resp, err = p.RunInstance(&RunInstanceRequest{})
	fmt.Println(resp, err)
}
