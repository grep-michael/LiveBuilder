package buildmanager

/*
Master object for doing all the backend building
1. imports all the custom files
2. runs the nessacary live-build commands
3. formats resulting files our specific desired output
*/

type UpdateType string

const (
	START  UpdateType = "start"
	UPDATE UpdateType = "start"
	END    UpdateType = "start"
)

type LogUpdate struct {
	UpdateType
	Message string
	Append  bool // true to append, false to replace
}

type BuildManager struct {
	updateChannel chan LogUpdate
	stepFinished  chan bool
	importer      Importer
	buildPath     string
}

func NewBuilder() *BuildManager {
	builder := &BuildManager{
		updateChannel: make(chan LogUpdate, 100),
		stepFinished:  make(chan bool),
	}
	builder.importer = *NewImporter(builder.updateChannel)
	return builder
}
func (self *BuildManager) Build(buildPath string) {
	self.buildPath = buildPath

}
