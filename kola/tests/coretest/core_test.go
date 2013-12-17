package coretest

import (
	"os"
	"os/exec"
	"path"
	"testing"
	"time"
)

const (
	CaPath      = "/usr/share/coreos-ca-certificates/"
	CmdTimeout  = time.Second * 3
	DbusTimeout = time.Second * 20
	HttpTimeout = time.Second * 3
	PortTimeout = time.Second * 3
	UpdateUrl   = "https://api.core-os.net/v1/update/"
)

func TestPortSsh(t *testing.T) {
	t.Parallel()
	err := CheckPort("tcp", "127.0.0.1:22", PortTimeout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateEngine(t *testing.T) {
	t.Parallel()

	errc := make(chan error, 1)
	go func() {
		c := exec.Command("update_engine_client", "-omaha_url", UpdateUrl)
		err := c.Run()
		errc <- err
	}()

	select {
	case <-time.After(CmdTimeout):
		t.Fatalf("update_engine_client timed out after %s.", CmdTimeout)
	case err := <-errc:
		if err != nil {
			t.Error(err)
		}
	}

	err := CheckDbusInterface("org.chromium.UpdateEngineInterface", DbusTimeout)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDockerEcho(t *testing.T) {
	t.Parallel()
	errc := make(chan error, 1)
	go func() {
		c := exec.Command("docker", "run", "busybox", "echo")
		err := c.Run()
		errc <- err
	}()
	select {
	case <-time.After(CmdTimeout):
		t.Fatalf("DockerEcho timed out after %s.", CmdTimeout)
	case err := <-errc:
		if err != nil {
			t.Error(err)
		}
	}
}

func TestUpdateServiceHttp(t *testing.T) {
	t.Parallel()
	err := CheckHttpStatus("http://api.core-os.net/v1/c10n/group", HttpTimeout)
	if err != nil {
		t.Error(err)
	}
}

func TestSymlinkResolvConf(t *testing.T) {
	t.Parallel()
	f, err := os.Lstat("/etc/resolv.conf")
	if err != nil {
		t.Fatal(err)
	}
	if !IsLink(f) {
		t.Fatal("/etc/resolv.conf is not a symlink.")

	}
}

func TestInstalledCACerts(t *testing.T) {
	t.Parallel()
	caCerts := []string{
		"CoreOS_Internet_Authority.pem",
		"CoreOS_Network_Authority.pem",
	}
	for _, fileName := range caCerts {
		_, err := os.Stat(path.Join(CaPath, fileName))
		if err != nil {
			t.Error(err)
		}
	}
}

func TestInstalledUpdateEngineRsaKeys(t *testing.T) {

}

func TestReadOnlyFs(t *testing.T) {
	mountModes := make(map[string]bool)
	mounts, err := GetMountTable()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range mounts {
		mountModes[m.MountPoint] = m.Options[0] == "ro"
	}
	if mp, ok := mountModes["/"]; ok {
		if mp {
			return
		} else {
			t.Fatalf("/ is not mounted ro.")
		}
	}
	t.Fatal("could not find rootfs.")
}
