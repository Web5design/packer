package common

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// A driver is able to talk to VirtualBox and perform certain
// operations with it. Some of the operations on here may seem overly
// specific, but they were built specifically in mind to handle features
// of the VirtualBox builder for Packer, and to abstract differences in
// versions out of the builder steps, so sometimes the methods are
// extremely specific.
type Driver interface {
	// Create a SATA controller.
	CreateSATAController(vm string, controller string) error

	// Delete a VM by name
	Delete(string) error

	// Import a VM
	Import(string, string) error

	// Checks if the VM with the given name is running.
	IsRunning(string) (bool, error)

	// Stop stops a running machine, forcefully.
	Stop(string) error

	// SuppressMessages should do what needs to be done in order to
	// suppress any annoying popups from VirtualBox.
	SuppressMessages() error

	// VBoxManage executes the given VBoxManage command
	VBoxManage(...string) error

	// Verify checks to make sure that this driver should function
	// properly. If there is any indication the driver can't function,
	// this will return an error.
	Verify() error

	// Version reads the version of VirtualBox that is installed.
	Version() (string, error)
}

func NewDriver() (Driver, error) {
	var vboxmanagePath string

	// On Windows, we check VBOX_INSTALL_PATH env var for the path
	if runtime.GOOS == "windows" {
		if installPath := os.Getenv("VBOX_INSTALL_PATH"); installPath != "" {
			log.Printf("[DEBUG] builder/virtualbox: VBOX_INSTALL_PATH: %s",
				installPath)
			for _, path := range strings.Split(installPath, ";") {
				path = filepath.Join(path, "VBoxManage.exe")
				if _, err := os.Stat(path); err == nil {
					vboxmanagePath = path
					break
				}
			}
		}
	}

	if vboxmanagePath == "" {
		var err error
		vboxmanagePath, err = exec.LookPath("VBoxManage")
		if err != nil {
			return nil, err
		}
	}

	log.Printf("VBoxManage path: %s", vboxmanagePath)
	driver := &VBox42Driver{vboxmanagePath}
	if err := driver.Verify(); err != nil {
		return nil, err
	}

	return driver, nil
}
