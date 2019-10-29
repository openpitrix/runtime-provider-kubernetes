module openpitrix.io/runtime-provider-kubernetes

go 1.13

require openpitrix.io/openpitrix v0.4.3-0.20191029152853-04510e51006a

require k8s.io/helm v2.14.3+incompatible

require k8s.io/kubernetes v1.14.1

require k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible

require k8s.io/cli-runtime v0.0.0-20190409023024-d644b00f3b79 // indirect

require k8s.io/cloud-provider v0.0.0-20190409023720-1bc0c81fa51d // indirect

replace github.com/gocraft/dbr => github.com/gocraft/dbr v0.0.0-20180507214907-a0fd650918f6

replace github.com/docker/docker => github.com/docker/engine v0.0.0-20190423201726-d2cfbce3f3b0

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190409021813-1ec86e4da56c

replace k8s.io/api => k8s.io/api v0.0.0-20190409021203-6e4e0e4f393b

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190404173353-6a84e37a896d

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190409022649-727a075fdec8

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/DATA-DOG/go-sqlmock v1.3.3 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/chai2010/gettext-go v0.0.0-20170215093142-bf70f2a70fb1 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000 // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20191011121108-aa519ddbe484 // indirect
	github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gobuffalo/packr v1.30.1 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/rubenv/sql-migrate v0.0.0-20191025130928-9355dd04f4b3 // indirect
	github.com/xlab/handysort v0.0.0-20150421192137-fb3537ed64a1 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	google.golang.org/grpc v1.24.0
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.0.0-20191025225708-5524a3672fbb
	k8s.io/apiextensions-apiserver v0.0.0-20191028232452-c47e10e6d5a3 // indirect
	k8s.io/apimachinery v0.0.0-20191025225532-af6325b3a843
	k8s.io/apiserver v0.0.0-00010101000000-000000000000 // indirect
	k8s.io/kube-openapi v0.0.0-20190918143330-0270cf2f1c1d // indirect
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4 // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	vbom.ml/util v0.0.0-20180919145318-efcd4e0f9787 // indirect
)
