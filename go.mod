module gitlab.pizza/nalbury/bailer

go 1.12

require (
	github.com/Azure/go-autorest/autorest/adal v0.6.0 // indirect
	github.com/contribsys/faktory v1.0.1-1.0.20190721213733-1cd2b871d15e
	github.com/contribsys/faktory_worker_go v0.0.0-20190428165239-86cdd9bae9d4
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/elazarl/goproxy v0.0.0-20170405201442-c4fc26588b6e // indirect
	github.com/evanphx/json-patch v4.2.0+incompatible // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/google/uuid v1.1.1 // indirect
	github.com/gophercloud/gophercloud v0.4.0 // indirect
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.7 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/robfig/cron/v3 v3.0.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/net v0.0.0-20190812203447-cdfb69ac37fc // indirect
	golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/inf.v0 v0.9.0 // indirect
	k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apimachinery v0.0.0-20190313205120-d7deff9243b1
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf // indirect
	k8s.io/utils v0.0.0-20190829053155-3a4a5477acf8 // indirect
)

//replace gopkg.in/robfig/cron.v3 => github.com/robfig/cron/v3 v3.0.0

replace github.com/contribsys/faktory => github.com/ClaytonNorthey92/faktory v1.0.1-1.0.20190721213733-1cd2b871d15e
