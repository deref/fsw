package rest

import (
	"net/http"

	"github.com/deref/fsw/internal/api"
)

type TagCollection struct {
	Service api.Service
}

func (coll *TagCollection) Subresource(key string) interface{} {
	return &TagKey{
		Service: coll.Service,
		Name:    key,
	}
}

type TagKey struct {
	Service api.Service
	Name    string
}

func (tk *TagKey) Subresource(key string) interface{} {
	return &Tag{
		Service:  tk.Service,
		TagName:  tk.Name,
		TagValue: key,
	}
}

type Tag struct {
	Service  api.Service
	TagName  string
	TagValue string
}

func (tag *Tag) Subresource(key string) interface{} {
	switch key {
	case "watchers":
		return &TagWatcherCollection{
			Service:  tag.Service,
			TagName:  tag.TagName,
			TagValue: tag.TagValue,
		}
	default:
		return nil
	}
}

type TagWatcherCollection struct {
	Service  api.Service
	TagName  string
	TagValue string
}

func (coll *TagWatcherCollection) HandleGet(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tags := make(map[string]string)
	tags[coll.TagName] = coll.TagValue
	output, err := coll.Service.DescribeWatchers(ctx, &api.DescribeWatchersInput{
		Tags: tags,
	})
	writeJSONResponse(w, req, output, err)
}

func (coll *TagWatcherCollection) HandleDelete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	tags := make(map[string]string)
	tags[coll.TagName] = coll.TagValue
	err := coll.Service.DeleteWatchers(ctx, &api.DeleteWatchersInput{
		Tags: tags,
	})
	writeErr(w, req, err)
}
