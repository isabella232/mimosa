package ui

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	// "html"
	"context"
	"io"
	"log"
	"text/template"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

//Host is a host from firestore
type Host struct {
	ID        string
	PublicDNS string
	PublicIP  string
	State     string
	Source    string
}

// HandleHTTPRequest handles a request and serves a UI by pulling host data from firestore
// gcloud functions deploy HandleHTTPRequest --runtime go111  --trigger-http
func HandleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	hosts, err := getHosts()
	if err != nil {
		log.Fatal(err)
	}
	err = merge(w, hosts)
	if err != nil {
		log.Fatal(err)
	}
}

func getHosts() ([]Host, error) {
	ctx := context.Background()
	fc, err := firestore.NewClient(ctx, firestore.DetectProjectID)
	if err != nil {
		return nil, err
	}

	hosts := []Host{}
	iter := fc.Collection("hosts").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		d := doc.Data()
		fmt.Println(d)

		h := hostFromMap(d)
		h.ID = doc.Ref.ID
		hosts = append(hosts, h)
	}

	// Sort for demo purposes
	sort.Slice(hosts, func(i, j int) bool {
		result := strings.Compare(hosts[i].Source, hosts[j].Source)
		if result != 0 {
			return result < 0
		}
		return hosts[i].PublicDNS < hosts[j].PublicDNS
	})

	return hosts, nil
}

func merge(w io.Writer, hosts []Host) error {
	// would probably be better to use a separate file, but uploading along with HTTP func may not be possible
	// tmpl := template.Must(template.ParseFiles("template.html"))
	tmpl := template.Must(template.New("anything").Parse(getHTML()))
	return tmpl.Execute(w, hosts)
}

func hostFromMap(m map[string]interface{}) Host {
	h := Host{}
	if v, ok := m["public_dns"]; ok {
		h.PublicDNS = v.(string)
	}
	if v, ok := m["public_ip"]; ok {
		h.PublicIP = v.(string)
	}
	if v, ok := m["state"]; ok {
		h.State = v.(string)
	}
	if v, ok := m["source"]; ok {
		h.Source = v.(string)
	}
	return h
}
