package cli

type Validator interface {
	Validate() error
}

type CmdParams struct {
	SrcDir  string
	DestDir string
	IsCopy  bool
	GroupBy string
}
