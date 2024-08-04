package main

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
)

// 生产实践：模型训练平台

func BuildDeployment() *appsv1.Deployment {
	// ...
	return nil
}

func DeployPredictService(deployment *appsv1.Deployment) error {
	// ...
	return nil
}

type Predictor interface {
	Deploy(deployment *appsv1.Deployment) error
	Scale(namespace, name string, replicas int) error
	Delete(namespace, name string) error
	// ...
}

type OpenPAIAdapter struct {
	// ...
}

func NewOpenPAIAdapter() *OpenPAIAdapter {
	return &OpenPAIAdapter{}
}

func (o *OpenPAIAdapter) Deploy(deployment *appsv1.Deployment) error {
	// 将 K8s Deployment 资源转换成 OpenPAI 的 RESTful 接口调用
	panic("implement me")
}

func (o *OpenPAIAdapter) Scale(namespace, name string, replicas int) error {
	panic("implement me")
}

func (o *OpenPAIAdapter) Delete(namespace, name string) error {
	panic("implement me")
}

func main() {
	// 部署推理服务原有流程
	{
		deployment := BuildDeployment()
		err := DeployPredictService(deployment)
		fmt.Println(err)
	}

	// 使用适配器部署推理服务流程
	{
		// 推理任务统一接口
		var predictor Predictor

		// 根据业务逻辑构造不同的适配器
		// switch expr {
		// case:
		predictor = NewOpenPAIAdapter()
		// }

		// 部署推理服务
		deployment := BuildDeployment()
		err := predictor.Deploy(deployment)
		fmt.Println(err)
	}

	// 将来想要改回来，只需要将旧代码也进行适配，写一个适配器
}
