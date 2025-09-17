package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/xuanviet96/seta-training/internal/config"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
)

type ESClient struct {
	Client *elastic.Client
}

func New(cfg config.Config, log *zap.Logger) (*ESClient, error) {
	if cfg.ESAddr == "" {
		return nil, errors.New("ES_ADDR empty")
	}
	es, err := elastic.NewClient(elastic.Config{
		Addresses: []string{cfg.ESAddr},
	})
	if err != nil {
		return nil, err
	}
	// ping
	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	log.Info("connected to elasticsearch", zap.String("addr", cfg.ESAddr))
	return &ESClient{Client: es}, nil
}

func EnsureIndex(ctx context.Context, es *ESClient, index string, log *zap.Logger) error {
	// check exists
	res, err := es.Client.Indices.Exists([]string{index})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return nil
	}
	// create with mapping
	body := map[string]any{
		"mappings": map[string]any{
			"properties": map[string]any{
				"id":      map[string]any{"type": "integer"},
				"title":   map[string]any{"type": "text", "fields": map[string]any{"keyword": map[string]any{"type": "keyword"}}},
				"content": map[string]any{"type": "text"},
				"tags":    map[string]any{"type": "keyword"},
			},
		},
	}
	buf, _ := json.Marshal(body)
	cr, err := es.Client.Indices.Create(index, es.Client.Indices.Create.WithBody(bytes.NewReader(buf)))
	if err != nil {
		return err
	}
	defer cr.Body.Close()
	if cr.IsError() {
		b, _ := io.ReadAll(cr.Body)
		return fmt.Errorf("create index error: %s", string(b))
	}
	log.Info("created es index", zap.String("index", index))
	return nil
}

type PostDoc struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags,omitempty"`
}

func IndexPost(ctx context.Context, es *ESClient, index string, doc PostDoc) error {
	data, _ := json.Marshal(doc)
	res, err := es.Client.Index(index, bytes.NewReader(data), es.Client.Index.WithDocumentID(fmt.Sprint(doc.ID)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("index error: %s", string(b))
	}
	return nil
}

func SearchPosts(ctx context.Context, es *ESClient, index, query string) ([]PostDoc, int, error) {
	body := map[string]any{
		"query": map[string]any{
			"multi_match": map[string]any{
				"query":  query,
				"fields": []string{"title", "content"},
			},
		},
	}
	b, _ := json.Marshal(body)
	res, err := es.Client.Search(
		es.Client.Search.WithIndex(index),
		es.Client.Search.WithBody(bytes.NewReader(b)),
		es.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()
	if res.IsError() {
		raw, _ := io.ReadAll(res.Body)
		return nil, 0, fmt.Errorf("es search error: %s", string(raw))
	}
	var out struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source PostDoc `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, 0, err
	}
	items := make([]PostDoc, 0, len(out.Hits.Hits))
	for _, h := range out.Hits.Hits {
		items = append(items, h.Source)
	}
	return items, out.Hits.Total.Value, nil
}
