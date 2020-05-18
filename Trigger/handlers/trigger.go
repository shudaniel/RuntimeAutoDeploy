package trigger

import (
	"RuntimeAutoDeploy/common"
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	dockertypes "github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/go-git/go-git"
)

func Cleanup(dir string) error {
	// Delete the dir folder and all repos inside
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("error reading directories in the current path")
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err.Error(),
			}).Error("error deleting directories from the current path")
			return err
		}
	}
	return nil
}

func tarDirectory(dir string, buf io.Writer) error {
	// Tar the dir folder
	// Reference: https://gist.github.com/mimoo/25fc9716e0f1353791f5908f94d6e726
	tw := tar.NewWriter(buf)
	filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
		// generate tar header

		// fmt.Println(file)
		// fmt.Println(fi.Name())

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			log.Error("FileInfoHeader")
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file)

		// write header
		if err := tw.WriteHeader(header); err != nil {
			log.Error("WriteHeader")
			return err
		}
		// if not a dir, write file content
		if !fi.IsDir() {

			data, err := os.Open(file)
			if err != nil {
				return err
			}

			if _, err := io.Copy(tw, data); err != nil {
				return err
			}
		}

		return nil
	})

	// produce tar
	if err := tw.Close(); err != nil {
		return err
	}

	return nil
}

func downloadGitRepo(gitrepo string) bool {

	_, err := git.PlainClone(common.GIT_BUILD_FOLDER, false, &git.CloneOptions{
		URL:      gitrepo,
		Progress: os.Stdout,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error cloning repo")
		return false
	}

	if _, err := os.Stat(common.GIT_BUILD_FOLDER + "Dockerfile"); os.IsNotExist(err) {
		// Dockerfile does not exist
		log.Error("git repo missing Dockerfile")
		return false
	}

	return true
}

func buildDockerImage(path string) error {
	// https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
	ctx := context.Background()
	cli, err := dockerclient.NewClientWithOpts(dockerclient.WithVersion("1.40")) // Max supported API version

	if err != nil {
		log.Fatal(err, " :unable to init client")
		return err
	}

	buf := new(bytes.Buffer)

	err = tarDirectory(path, buf)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(" :unable to tar directory")
		return err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		dockertypes.ImageBuildOptions{
			Context:    dockerFileTarReader,
			Dockerfile: common.GIT_BUILD_FOLDER + "Dockerfile",
			Remove:     true})
	if err != nil {

		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(" :unable to build docker image")
		return err
	}
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(err, " :unable to read image build response")
		return err
	}

	return nil
}

// Create a TCP Server listening on a port and return the port
func createTCPSocket() (int, net.Listener) {
	// Keep trying ports until either you exhaust all ports
	// or one is available
	// If no ports are available, return 0
	var port int
	var l net.Listener
	var err error
	for port = 8000; port <= 65535; port++ {
		l, err = net.Listen("tcp4", fmt.Sprintf(":%d", port))
		if err == nil {
			fmt.Println(fmt.Sprintf("Port %d selected", port))
			return port, l
		}
	}
	return 0, nil
}

func acceptTCPConnection(l net.Listener) {
	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	handleConnection(c)
}

func handleConnection(c net.Conn) {
	// STUB handler function

	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	c.Close()
}

func handleConfigFile(config *common.RADConfig) bool {
	// Parse the config file
	// Download the git repository into local Trigger/build folder
	// Check for Dockerfile. If does not exist, quit
	// else build docker image and store within Trigger/images
	fmt.Println(config.GitRepoLink)

	if !downloadGitRepo(config.GitRepoLink) {
		return false
	}
	err := buildDockerImage(common.GIT_BUILD_FOLDER)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error building docker image")
		return false
	}

	return true
}

func RADTriggerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Error("error. Received incorrect HTTP method. Expecting POST")
		return
	}
	port, l := createTCPSocket()
	go acceptTCPConnection(l)
	_, _ = w.Write([]byte(fmt.Sprintf("%d", port)))

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error reading body")
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	data := common.RADConfig{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error unmarshal body")
		return
	}
	_ = handleConfigFile(&data)
	err = Cleanup(common.GIT_BUILD_FOLDER)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error clearing GIT_BUILD_FOLDER")
	}
}
