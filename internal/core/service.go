package core

import (
	"errors"
	"github.com/LibCyber/cyber/pkg/util"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var (
	ErrNotLinux = errors.New("only support linux currently")
)

func InstallService() error {
	// only support linux currently
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS != "linux" {
		return ErrNotLinux
	}

	isAdmin, err := util.IsAdmin()
	if err != nil {
		return err
	}

	if !isAdmin {
		return ErrNotAdmin
	}

	// judge if service installed (/etc/systemd/system/cyber-core.service)
	servicePath := filepath.Join("/etc", "systemd", "system", "cyber-core.service")
	if _, err = os.Stat(servicePath); os.IsExist(err) {
		return err
	}

	// judge if core installed, ~/.cyber/core/cyber-core
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	corePath := filepath.Join(homeDir, ".cyber", "core", "cyber-core")
	if _, err = os.Stat(corePath); os.IsNotExist(err) {
		return err
	}

	// copy ~/.cyber/core/cyber-core to /usr/bin/cyber-core
	err = util.CopyFile(corePath, filepath.Join("/usr", "bin", "cyber-core"))
	if err != nil {
		return err
	}

	// add executable permission to /usr/bin/cyber-core
	err = os.Chmod(filepath.Join("/usr", "bin", "cyber-core"), 0755)
	if err != nil {
		return err
	}

	// judge if node config file exists
	nodeConfigPath := filepath.Join(homeDir, ".cyber", "node", "config.yaml")
	if _, err = os.Stat(nodeConfigPath); os.IsNotExist(err) {
		return err
	}

	// check if /etc/cyber-core exists, if not, create it
	if _, err = os.Stat(filepath.Join("/etc", "cyber-core")); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join("/etc", "cyber-core"), 0755)
		if err != nil {
			return err
		}
	}

	// copy ~/.cyber/core/Country.mmdb to /etc/cyber-core/Country.mmdb
	err = util.CopyFile(filepath.Join(homeDir, ".cyber", "core", "Country.mmdb"), filepath.Join("/etc", "cyber-core", "Country.mmdb"))
	if err != nil {
		return err
	}

	// copy ~/.cyber/core/cache.db to /etc/cyber-core/cache.db
	err = util.CopyFile(filepath.Join(homeDir, ".cyber", "core", "cache.db"), filepath.Join("/etc", "cyber-core", "cache.db"))
	if err != nil {
		return err
	}

	// copy ~/.cyber/node/config.yaml to /etc/cyber-core/config.yaml
	err = util.CopyFile(nodeConfigPath, filepath.Join("/etc", "cyber-core", "config.yaml"))
	if err != nil {
		return err
	}

	// judge if core running
	pid, err := Status()
	if err != nil {
		return err
	}

	if pid != 0 {
		return errors.New("core is running, please stop it first with `cyber core stop`")
	}

	// check or create log path and files, /var/log/cyber-core, /var/log/cyber-core/coreErr.log, /var/log/cyber-core/coreOut.log
	if _, err = os.Stat(filepath.Join("/var", "log", "cyber-core")); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Join("/var", "log", "cyber-core"), 0755)
		if err != nil {
			return err
		}
	}

	if _, err = os.Stat(filepath.Join("/var", "log", "cyber-core", "coreErr.log")); os.IsNotExist(err) {
		_, err = os.Create(filepath.Join("/var", "log", "cyber-core", "coreErr.log"))
		if err != nil {
			return err
		}
	}

	if _, err = os.Stat(filepath.Join("/var", "log", "cyber-core", "coreOut.log")); os.IsNotExist(err) {
		_, err = os.Create(filepath.Join("/var", "log", "cyber-core", "coreOut.log"))
		if err != nil {
			return err
		}
	}

	// install service
	serviceTemplate := `[Unit]
Description=Cyber Core Service
Documentation=https://docs.libcyber.org/docs/quan-xin-libcyber-ke-hu-duan/cli
After=network.target nss-lookup.target

[Service]
Type=simple
ExecStart=/usr/bin/cyber-core -d /etc/cyber-core -f /etc/cyber-core/config.yaml
Restart=on-failure
LimitNOFILE=1048550
RestartSec=5

[Install]
WantedBy=multi-user.target`

	// create service file, /etc/systemd/system/cyber-core.service, with template above
	err = os.WriteFile(servicePath, []byte(serviceTemplate), 0644)
	if err != nil {
		return err
	}

	// execute systemctl daemon-reload
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func UninstallService() error {
	// only support linux currently
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS != "linux" {
		return ErrNotLinux
	}

	isAdmin, err := util.IsAdmin()
	if err != nil {
		return err
	}

	if !isAdmin {
		return ErrNotAdmin
	}

	// located in /etc/systemd/system/cyber-core.service
	servicePath := filepath.Join("/etc", "systemd", "system", "cyber-core.service")

	// judge if service running
	pid, err := Status()
	if err != nil {
		return err
	}

	if pid != 0 {
		return errors.New("core is running, please stop it first with `cyber core stop`")
	}

	// uninstall service
	err = os.RemoveAll(servicePath)
	if err != nil {
		return err
	}

	// remove /usr/bin/cyber-core
	err = os.RemoveAll(filepath.Join("/usr", "bin", "cyber-core"))
	if err != nil {
		return err
	}

	// remove /etc/cyber-core
	err = os.RemoveAll(filepath.Join("/etc", "cyber-core"))
	if err != nil {
		return err
	}

	// remove /var/log/cyber-core
	err = os.RemoveAll(filepath.Join("/var", "log", "cyber-core"))
	if err != nil {
		return err
	}

	// execute systemctl daemon-reload
	cmd := exec.Command("systemctl", "daemon-reload")
	err = cmd.Run()
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func IsServiceInstalled() (bool, error) {
	// only support linux currently
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS != "linux" {
		return false, nil
	}

	isAdmin, err := util.IsAdmin()
	if err != nil {
		return false, err
	}

	if !isAdmin {
		return false, ErrNotAdmin
	}

	// judge if service installed (/etc/systemd/system/cyber-core.service)
	servicePath := filepath.Join("/etc", "systemd", "system", "cyber-core.service")
	if _, err = os.Stat(servicePath); os.IsNotExist(err) {
		return false, nil
	}

	return true, nil
}

func Purge() error {
	// delete ~/.cyber
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(homeDir, ".cyber"))
	if err != nil {
		return err
	}

	// if linux and service installed, delete /etc/systemd/system/cyber-core.service, /usr/bin/cyber-core, /etc/cyber-core, /var/log/cyber-core
	//goland:noinspection GoBoolExpressions
	if runtime.GOOS == "linux" {
		isServiceInstalled, err := IsServiceInstalled()
		if err != nil {
			return err
		}

		if isServiceInstalled {
			err = UninstallService()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
