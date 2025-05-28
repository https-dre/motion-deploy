package core

import "motion/pkgs/models"

type CoreInstance interface {
	ListContainers() []models.CoreApplication
	BuildAndRun(projectPath, imageName, containerName string) (models.CoreApplication, error)
}