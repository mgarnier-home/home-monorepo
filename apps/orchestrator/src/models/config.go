package models

const (
	ModeFullLocal string = "local"
	ModeFullApi   string = "remote"
	ModeHybrid    string = "hybrid"
)

type OrchestratorConfig struct {
	// Commun CLI/API
	S3Endpoint  string `key:"ORCHESTRATOR_S3_ENDPOINT" required:"true"`
	S3AccessKey string `key:"ORCHESTRATOR_S3_ACCESS_KEY" required:"true"`
	S3SecretKey string `key:"ORCHESTRATOR_S3_SECRET_KEY" required:"true"`
	S3Bucket    string `key:"ORCHESTRATOR_S3_BUCKET" required:"true"`

	GitToken string `key:"ORCHESTRATOR_GIT_TOKEN" default-value:""`
	GitRepo  string `key:"ORCHESTRATOR_GIT_REPO" default-value:""`

	SSHPrivateKey  string `key:"ORCHESTRATOR_SSH_PRIVATE_KEY" default-value:""`
	ComposeDirPath string `key:"ORCHESTRATOR_COMPOSE_DIRECTORY_PATH" default-value:""`

	// CLI
	Mode   string `key:"ORCHESTRATOR_MODE" default-value:"local"`
	ApiUrl string `key:"ORCHESTRATOR_API_URL" default-value:"http://localhost:3000"`

	// API
	ServerPort   int    `key:"ORCHESTRATOR_SERVER_PORT" default-value:"3000"`
	BinariesPath string `key:"ORCHESTRATOR_BINARIES_PATH" default-value:"/dist"`
}
