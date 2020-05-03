package trigger

import (
    "RuntimeAutoDeploy/common"
    "fmt"
    "net/http"
    "net"
    "encoding/json"
    "io/ioutil"
    log "github.com/Sirupsen/logrus"
)

type Tuple struct {
    Port int 
    L net.Listener
}

// Create a TCP Server listening on a port and return the port
func CreateTCPSocket() Tuple {
    // Keep trying ports until either you exhaust all ports 
    // or one is available
    // If no ports are available, return 0
    var port int
    var l net.Listener
    var err error
    for port = 8000; port <= 65535; port++ {
        l, err = net.Listen("tcp4", fmt.Sprintf(":%d", port))
        if err == nil {
            fmt.Println(fmt.Sprintf("Port %d selected", port) )
            return Tuple{port, l}
        }
    }


    return Tuple{0, nil}
}

func HandleConnection(c net.Conn) {
    // STUB handler function

    fmt.Printf("Serving %s\n", c.RemoteAddr().String())

    c.Close()
}

func HandleConfigFile(config *common.RADConfig) {
    // Parse the config file
    // Download the git repository into local Trigger/build folder
    // Check for Dockerfile. If does not exist, quit
    // else build docker image and store within Trigger/images
    fmt.Println(config.GitRepoLink)
}

func AcceptTCPConnection(l net.Listener) {
    c, err := l.Accept()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer c.Close()

    HandleConnection(c)
}

// AppPreferencesHandler handles the requests from the apps
func AppPreferencesHandler(w http.ResponseWriter, r *http.Request) {
    var port int
    if r.Method == "POST" {
        // ctx, _ := context.WithCancel(r.Context())
        ret := CreateTCPSocket()
        port = ret.Port
        go AcceptTCPConnection(ret.L)
        w.Write([]byte(fmt.Sprintf("%d", port)))
        
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

        HandleConfigFile(&data)
    } else {
        log.Error("error. Received incorrect HTTP method. Expecting POST")
    }
    return
}

