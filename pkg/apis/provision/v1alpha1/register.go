package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "provision.mydomain.com"

var SchemeGroupVersion = schema.GroupVersion{
	Group:   GroupName,
	Version: "v1alpha1",
}

// Resource generate a group resource instance based on the given resource name.
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// var 的定义和 internal version 的 register 中基本类似，
// 只是创建 Builder 时多了一个中间产物local scheme builder，
// local builder 会在包括生成代码的init中去使用.
var (
	SchemeBuilder      runtime.SchemeBuilder
	localSchemeBuilder = &SchemeBuilder
	AddToScheme        = localSchemeBuilder.AddToScheme
)

func init() {
	// 这里去注册该 version 的类型，以及它们向 internal version 的转换函数.
	localSchemeBuilder.Register(addKnownTypes)
}

// addKnownTypes 被 SchemeBuilder 调用，从而把自己知道的 Object(Type) 注册到 scheme 中.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&ProvisionRequest{},
		&ProvisionRequestList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
