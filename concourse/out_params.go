package concourse

type OutParams struct {
	Manifest    string                 `json:"manifest"`
	NoRedact    bool                   `json:"no_redact,omitempty"`
	DryRun      bool                   `json:"dry_run,omitempty"`
	MaxInFlight int                    `json:"max_in_flight,omitempty"`
	Recreate    bool                   `json:"recreate,omitempty"`
	SkipDrain   []string               `json:"skip_drain,omitempty"`
	Cleanup     bool                   `json:"cleanup,omitempty"`
	Releases    []string               `json:"releases,omitempty"`
	Stemcells   []string               `json:"stemcells,omitempty"`
	Vars        map[string]interface{} `json:"vars,omitempty"`
	VarsFiles   []string               `json:"vars_files,omitempty"`
	VarFiles    map[string]string      `json:"var_files,omitempty"`
	OpsFiles    []string               `json:"ops_files,omitempty"`
	Delete      DeleteParams           `json:"delete,omitempty"`
}

type DeleteParams struct {
	Enabled bool `json:"enabled,omitempty"`
	Force   bool `json:"force,omitempty"`
}
