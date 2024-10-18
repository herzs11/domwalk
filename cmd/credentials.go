package cmd

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func getCredentials() (*oauth2.Token, error) {
	ctx := context.Background()

	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate default credentials, make sure your GOOGLE_APPLICATION_CREDENTIALS environment variable is set: %w",
			err,
		)
	}

	ts, err := idtoken.NewTokenSource(ctx, ENRICH_DOMAIN_CF_URL, option.WithCredentials(credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create NewTokenSource: %w", err)
	}

	tok, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to receive token: %w", err)
	}
	return tok, nil
}
