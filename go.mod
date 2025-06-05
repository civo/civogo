module github.com/civo/civogo

go 1.23.0

toolchain go1.23.8

require (
	github.com/google/go-querystring v1.1.0
	github.com/onsi/gomega v1.27.4
	k8s.io/api v0.27.1
	k8s.io/apimachinery v0.27.1
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	golang.org/x/mod v0.17.0
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/utils v0.0.0-20230209194617-a36077c30491 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace github.com/onsi/gomega => github.com/onsi/gomega v1.19.0
