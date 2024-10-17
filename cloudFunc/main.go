package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"domwalk/domains"
	"domwalk/stores/bq"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

type requestParams struct {
	DomainNames []string `json:"domain_names"`
	processConfig
	OnlyMatchedDomains bool `json:"only_matched_domains,omitempty"`
}

func init() {
	bqs, err := bq.NewBQStore("unum-marketing-data-assets", "domwalk_dev", "domains_test")
	if err != nil {
		log.Fatal(err)
	}
	functions.HTTP("handleDomainEnrichment", handleDomainEnrichment(bqs))
}

func main() {
	if err := funcframework.Start("8080"); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}

func handleDomainEnrichment(bqs *bq.BQStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rParams requestParams
		err := json.NewDecoder(r.Body).Decode(&rParams)
		if err != nil {
			writeJSON(
				w, http.StatusBadRequest,
				map[string]string{"error": fmt.Sprintf("Unable to parse json: %s", err.Error())},
			)
			return
		}
		doms, err := bqs.GetDomainsByNames(context.Background(), rParams.DomainNames)
		if err != nil {
			writeJSON(
				w, http.StatusInternalServerError, map[string]string{"error": "Unable to get domains from BQ"},
			)
			return
		}
		enrichDomains(doms, rParams.processConfig)
		go bqs.PutDomains(context.Background(), doms)
		if rParams.OnlyMatchedDomains {
			matchedDoms := make(map[string]domains.MatchedDomainsByStrategy)
			for _, dom := range doms {
				matchedDoms[dom.DomainName] = dom.GetAllMatchedDomains()
			}
			writeJSON(w, http.StatusOK, matchedDoms)
		} else {
			writeJSON(w, http.StatusOK, doms)
		}
		return
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	return json.NewEncoder(rw).Encode(v)
}
