package entrypoints

type CLI struct {
	Job JobCmd `cmd:"" help:"pull jobs from certain sources (upwork, weworkremotely)"`
}
