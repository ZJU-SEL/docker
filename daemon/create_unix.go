// +build !windows

package daemon

import (
	"fmt"
	"os"

	"github.com/docker/docker/image"
	"github.com/docker/docker/runconfig"
	"github.com/docker/docker/volume"
	"github.com/opencontainers/runc/libcontainer/label"
)

// createContainerPlatformSpecificSettings performs platform specific container create functionality
func createContainerPlatformSpecificSettings(container *Container, config *runconfig.Config, img *image.Image) error {
	for spec := range config.Volumes {
		name, destination, err := parseVolumes(spec)
		if err != nil {
			return err
		}

		// Skip volumes for which we already have something mounted on that
		// destination because of a --volume-from.
		if container.isDestinationMounted(destination) {
			continue
		}
		path, err := container.GetResourcePath(destination)
		if err != nil {
			return err
		}

		stat, err := os.Stat(path)
		if err == nil && !stat.IsDir() {
			return fmt.Errorf("cannot mount volume over existing file, file exists %s", path)
		}

		volumeDriver := config.VolumeDriver
		if destination != "" && img != nil {
			if _, ok := img.ContainerConfig.Volumes[destination]; ok {
				// check for whether bind is not specified and then set to local
				if _, ok := container.MountPoints[destination]; !ok {
					volumeDriver = volume.DefaultDriverName
				}
			}
		}

		v, err := createVolume(name, volumeDriver)
		if err != nil {
			return err
		}
		if err := label.Relabel(v.Path(), container.MountLabel, "z"); err != nil {
			return err
		}

		// never attempt to copy existing content in a container FS to a shared volume
		if volumeDriver == volume.DefaultDriverName || volumeDriver == "" {
			if err := container.copyImagePathContent(v, destination); err != nil {
				return err
			}
		}

		container.addMountPointWithVolume(destination, v, true)
	}
	return nil
}
