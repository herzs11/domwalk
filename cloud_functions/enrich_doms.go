package cloud_functions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"dev.azure.com/Unum/Mkt_Analytics/_git/cloud_functions/types"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/herzs11/domwalk/domains"
	"github.com/herzs11/domwalk/stores/bq"
)

func init() {
	bqs, err := bq.NewBQStore("unum-marketing-data-assets", "domwalk", "domains")
	if err != nil {
		log.Fatal(err)
	}
	functions.HTTP("enrich", handleDomainEnrichment(bqs))
}

func handleDomainEnrichment(bqs *bq.BQStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rParams types.RequestParams
		err := json.NewDecoder(r.Body).Decode(&rParams)
		if err != nil {
			writeJSON(
				w, http.StatusBadRequest,
				map[string]string{"error": fmt.Sprintf("Unable to parse json: %s", err.Error())},
			)
			return
		}
		log.Println(rParams)
		doms, err := bqs.GetDomainsByNames(context.Background(), rParams.DomainNames)
		if err != nil {
			writeJSON(
				w, http.StatusInternalServerError, map[string]string{"error": "Unable to get domains from BQ"},
			)
			return
		}
		enrichDomains(doms, rParams.ProcessConfig)
		go log.Println(bqs.PutDomains(context.Background(), doms))
		if rParams.NoResponse {
			writeJSON(w, http.StatusOK, map[string]string{"message": "Enriched domains"})
			return
		}
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
