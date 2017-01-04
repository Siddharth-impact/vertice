package carton

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/cmd"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/pairs"
	"gopkg.in/yaml.v2"
	"io"
	"encoding/json"
	"strings"
	"time"
)

const (
	SNAPSHOTBUCKET = "snapshots"
	DISKSBUCKET    = "disks"
	ACCOUNTID      = "account_id"
	ASSEMBLYID     = "asm_id"
)


type ApiSnaps struct {
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	Results    []Snaps  `json:"results" cql:"results"`
}


//The grand elephant for megam cloud platform.
type Snaps struct {
	Id         string `json:"id" cql:"id"`
	ImageId    string `json:"image_id" cql:"image_id"`
	OrgId      string `json:"org_id" cql:"org_id"`
	AccountId  string `json:"account_id" cql:"account_id"`
	Name       string `json:"name" cql:"name"`
	AssemblyId string `json:"asm_id" cql:"asm_id"`
	JsonClaz   string `json:"json_claz" cql:"json_claz"`
	CreatedAt  string `json:"created_at" cql:"created_at"`
	Status     string `json:"status" cql:"status"`
	Tosca      string `json:"tosca_type" cql:"tosca_type"`
	Inputs     pairs.JsonPairs `json:"inputs" cql:"inputs"`
	Outputs    pairs.JsonPairs `json:"inputs" cql:"inputs"`
}

func (a *Snaps) String() string {
	if d, err := yaml.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(d)
	}
}

// ChangeState runs a state increment of a machine or a container.
func SaveImage(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].SaveImage(opts.B, writer)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}
	slog := outBuffer.String()
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))
	return nil
}

// ChangeState runs a state increment of a machine or a container.
func DeleteImage(opts *DiskOpts) error {
	var outBuffer bytes.Buffer
	start := time.Now()
	logWriter := LogWriter{Box: opts.B}
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, &logWriter)
	err := ProvisionerMap[opts.B.Provider].DeleteImage(opts.B, writer)
	elapsed := time.Since(start)

	if err != nil {
		return err
	}
	slog := outBuffer.String()
	log.Debugf("%s in (%s)\n%s",
		cmd.Colorfy(opts.B.GetFullName(), "cyan", "", "bold"),
		cmd.Colorfy(elapsed.String(), "green", "", "bold"),
		cmd.Colorfy(slog, "yellow", "", ""))
	return nil
}

/** A public function which pulls the snapshot for disk save as image.
and any others we do. **/
func GetSnap(id , email string) (*Snaps, error) {
	cl := api.NewClient(newArgs(email, ""), "/snapshots/" + id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	res := &ApiSnaps{}
	err = json.Unmarshal(response, res)
	if err != nil {
		return nil, err
	}
	a := &res.Results[0]
	log.Debugf("Snaps %v", a)
	return a, nil
}

func (s *Snaps) UpdateSnap() error {
	cl := api.NewClient(newArgs(s.AccountId, s.OrgId),"/snapshots/update" )
	if _, err := cl.Post(s); err != nil {
		return err
	}
	return nil

}


func (a *Snaps) RemoveSnap() error {
	cl := api.NewClient(newArgs(a.AccountId, a.OrgId), "/snapshots/" + a.Id)
	if	_, err := cl.Delete(); err != nil {
		return err
	}
	return nil
}

//make cartons from snaps.
func (a *Snaps) MkCartons() (Cartons, error) {
	newCs := make(Cartons, 0, 1)
	if len(strings.TrimSpace(a.AssemblyId)) > 1 {
		if ca, err := mkCarton(a.Id, a.AssemblyId, a.AccountId); err != nil {
			return nil, err
		} else {
			ca.toBox()                //on success, make a carton2box if BoxLevel is BoxZero
			newCs = append(newCs, ca) //on success append carton
		}
	}
	log.Debugf("Cartons %v", newCs)
	return newCs, nil
}
