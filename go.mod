module github.com/SENERGY-Platform/device-manager

go 1.20

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414 // indirect
	github.com/segmentio/kafka-go v0.4.42
	github.com/wvanbergen/kazoo-go v0.0.0-20180202103751-f72d8611297a
)

require (
	github.com/SENERGY-Platform/models/go v0.0.0-20230406081245-2b17534509d4
	github.com/SENERGY-Platform/permission-search v0.0.0-20230802063812-b193f9f6cff3
	github.com/testcontainers/testcontainers-go v0.19.0
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/PaesslerAG/gval v1.2.2 // indirect
	github.com/PaesslerAG/jsonpath v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/containerd/containerd v1.6.19 // indirect
	github.com/cpuguy83/dockercfg v0.3.1 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v23.0.1+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mdaverde/jsonpath v0.2.1 // indirect
	github.com/moby/patternmatcher v0.5.0 // indirect
	github.com/moby/sys/sequential v0.5.0 // indirect
	github.com/moby/term v0.0.0-20221128092401-c43b287e0e0f // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/opensearch-project/opensearch-go v1.1.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	google.golang.org/genproto v0.0.0-20220617124728-180714bec0ad // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

//replace github.com/SENERGY-Platform/models/go => ../models/go
