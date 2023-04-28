package helpers

type Context struct {
	DryRun bool `help:"dry run. no persistence will be involved"`
	Debug  bool `help:"increase verbosity"`
}
