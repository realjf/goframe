package api

import (
	"kboard/config"
	"kboard/k8s"
	"kboard/k8s/resource"
	"net/http"
)

type IOrder struct {
	IApi
}

func NewIOrder(config config.IConfig, w http.ResponseWriter, r *http.Request) *IOrder {
	order := &IOrder{
		IApi: *NewIApi(config, w, r),
	}
	order.Module = "order"
	return order
}

func (this *IOrder) Index() {

	this.TplEngine.Response(100, "", "")
}

// @todo 创建工单
func (this *IOrder) Edit() {

}

func (this *IOrder) Save() {

}

func (this *IOrder) List() {
	resReplicaSet := resource.NewResReplicaSet()
	resReplicaSet.SetMetadataName("hello")
	resReplicaSet.SetReplicas(3)
	resReplicaSet.SetNamespace("myapp")

	labels := map[string]string{
		"app": "nginx",
	}

	resReplicaSet.SetSelector(resource.Selector{
		MatchLabels: labels,
	})

	container := resource.NewContainer("mycontainer", "image")
	container.SetResource(resource.ContainerResources{
		Limits: resource.Limits{
			Cpu:    "0.5",
			Memory: "100Mi",
		},
		Requests: resource.Request{
			Cpu:    "0.1",
			Memory: "50Mi",
		},
	})
	container.SetLivenessProbe(resource.LivenessProbe{
		ProbeAction: resource.ProbeAction{
			Exec: &resource.ExecAction{
				Command: []string{
					"/bin/sh",
					"-c",
				},
			},
		},
		InitialDelaySeconds: 50,
		PeriodSeconds:       10,
		TimeoutSeconds:      10,
		SuccessThreshold:    1,
		FailureThreshold:    10,
	})
	resReplicaSet.AddContainer(container)
	resReplicaSet.SetTemplateLabel(labels)
	resReplicaSet.SetLabels(labels)
	yamlData, err := resReplicaSet.ToYamlFile()
	if err != nil {
		this.TplEngine.Response(99, err, "错误")
	}
	lib := k8s.NewReplicaSet(this.Config)

	res := lib.WriteToEtcd("myapp", "hello", yamlData)
	this.TplEngine.Response(100, res, "数据")
}
