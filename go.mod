module github.com/civo/civogo

go 1.16

require (
	github.com/google/go-querystring v1.1.0
	github.com/onsi/gomega v1.27.4
	k8s.io/api v0.27.1
)

replace (
	github.com/onsi/gomega => github.com/onsi/gomega v1.19.0
)
