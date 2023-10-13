module github.com/hyunkim13/kt-driver

go 1.18

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible

require (
	github.com/docker/machine v0.16.2
	github.com/rancher/machine v0.13.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.7.2 // indirect
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/term v0.0.0-20220722155259-a9ba230a4035 // indirect
	gotest.tools v2.2.0+incompatible // indirect
)
