package workflow

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/elumbantoruan/feed/cmd/cronjob/mock"
	"github.com/stretchr/testify/assert"
)

func TestWorkflow_Run(t *testing.T) {
	client := mock.MockGRPCClient{}
	crawler := mock.MockCrawler{}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	work := New(&client, logger, &crawler)
	err := work.Run(context.Background())
	assert.NoError(t, err)

	// GetSites only called once
	assert.Equal(t, 1, client.GetSitesCount)
	// UpdateSiteFeedCount called twice
	assert.Equal(t, 2, client.UpdateSiteFeedCount)
	// AddArticleCount called twice
	assert.Equal(t, 2, client.AddArticleCount)

}
