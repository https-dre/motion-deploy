package core

type CoreInstance interface {
	ListApplications() []CoreApplication
	BuildAndRunService(projectPath, imageName, containerName string) error
}