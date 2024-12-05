module mgarnier11/go-proxy

go 1.23.0

require (
	github.com/charmbracelet/lipgloss v1.0.0
	github.com/docker/docker v27.2.0+incompatible
	github.com/go-ping/ping v1.1.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.30.0
	gopkg.in/yaml.v3 v3.0.1
	mgarnier11/go/colors v0.0.0-00010101000000-000000000000
	mgarnier11/go/dockerssh v0.0.0-00010101000000-000000000000
	mgarnier11/go/logger v0.0.0-00010101000000-000000000000
	mgarnier11/go/utils v0.0.0-00010101000000-000000000000

)

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/charmbracelet/log v0.4.0 // indirect
	github.com/charmbracelet/x/ansi v0.5.2 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/sdk v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/exp v0.0.0-20241108190413-2d47ceb2692f // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	gotest.tools/v3 v3.5.1 // indirect
)

replace mgarnier11/go/utils => ../../../libs/go/utils

replace mgarnier11/go/logger => ../../../libs/go/logger

replace mgarnier11/go/colors => ../../../libs/go/colors

replace mgarnier11/go/dockerssh => ../../../libs/go/dockerssh
