package integrations

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type MysqlDumper struct {
	User     string
	Password string
	Database string
	Outfile  string
}

func (md *MysqlDumper) Dump() error {
	mysqldumpExe, err := exec.LookPath("mysqldump")

	if err != nil {
		return fmt.Errorf("could not lookup the mysqldump location: %w", err)
	}

	cmd := exec.Command(
		"sh", "-c",
		strings.Join([]string{
			mysqldumpExe, "-u", md.User, "-p" + md.Password,
			md.Database, "|", "gzip", ">", md.Outfile,
		}, " "),
	)

	cmd.Env = os.Environ()

	log.Println("Running the command:", cmd.String())

	out, err := cmd.CombinedOutput()

	log.Println("Dump command output:", string(out))

	if err != nil {
		return fmt.Errorf("could not run the docker command: %w", err)
	}

	return nil
}
