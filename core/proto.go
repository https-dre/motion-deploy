package core

import "motion/pkgs/models"

type CoreInstance interface {
	ListApplications() []models.CoreApplication
	BuildAndRunService(projectPath, imageName, containerName string) error
}