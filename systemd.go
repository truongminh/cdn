package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

type control struct {
	unit     string
	username string
	group    string
	runDir   string
	logDir   string
	progFile string
}

func (c *control) systemdDescription() string {
	str := `[Unit]
	Description=` + c.unit + `
	ConditionPathExists=` + c.runDir + `
	After=network.target
	
	[Service]
	Type=simple
	User=` + c.username + `
	Group=` + c.group + `
	LimitNOFILE=4096
	
	Restart=on-failure
	RestartSec=10
	StartLimitIntervalSec=60
	
	WorkingDirectory=` + c.runDir + `
	ExecStart=` + c.progFile + `
	
	# make sure log directory exists and owned by syslog
	PermissionsStartOnly=true
	ExecStartPre=/bin/mkdir -p ` + c.logDir + `
	ExecStartPre=/bin/chmod 755 ` + c.logDir + `
	StandardOutput=syslog
	StandardError=syslog
	SyslogIdentifier=` + c.unit + `
	
	[Install]
	WantedBy=multi-user.target
	`
	return str
}

func (c *control) Install(ctx context.Context) (err error) {
	u, err := user.Lookup(c.username)
	if err != nil {
		return
	}
	c.group = u.Gid
	if c.runDir == "" {
		c.runDir = filepath.Join(u.HomeDir, c.unit)
	}
	err = os.MkdirAll(c.runDir, 0700)
	if err != nil {
		return
	}

	c.logDir = filepath.Join(c.runDir, "log")
	err = os.MkdirAll(c.logDir, 0700)
	if err != nil {
		return
	}
	// copy file
	// argv0 := os.Args[0]
	c.progFile = c.runDir + "/" + c.unit
	// err = copyFile(argv0, c.progFile)
	// if err != nil {
	// return
	// }
	// err = os.Chmod(c.progFile, 0700)
	// if err != nil {
	// return
	// }

	service := []byte(c.systemdDescription())
	file := fmt.Sprintf("/lib/systemd/system/%s.service", c.unit)
	err = ioutil.WriteFile(file, service, 0700)
	if err != nil {
		return
	}
	log.Printf("installed service %s", c.unit)
	return
}

func copyFile(src string, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return
}

func systemctlCmd(ctx context.Context, args ...string) (cmd *exec.Cmd, err error) {
	cmd = exec.CommandContext(ctx, "systemctl", args...)
	err = cmd.Run()
	return
}
