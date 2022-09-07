package Era

import (
	"context"
	"fmt"
	"net/http"

	"github.com/terraform-providers/terraform-provider-nutanix/client"
)

type Service interface {
	ListProfiles(ctx context.Context, engine string, profileType string) (*ProfileListResponse, error)
	GetProfiles(ctx context.Context, engine string, profileType string, id string, name string) (*ListProfileResponse, error)
	ListClusters(ctx context.Context) (*ListClusterResponse, error)
	ListSLA(ctx context.Context) (*SLAResponse, error)
}

type ServiceClient struct {
	c *client.Client
}

func (sc ServiceClient) ListProfiles(ctx context.Context, engine string, profile_type string) (*ProfileListResponse, error) {
	var httpReq *http.Request
	var err error

	if engine != "" && profile_type != "" {
		path := fmt.Sprintf("/profiles?engine=%s&type=%s", engine, profile_type)
		httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	} else if engine != "" {
		path := fmt.Sprintf("/profiles?engine=%s", engine)
		httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	} else if profile_type != "" {
		path := fmt.Sprintf("/profiles?type=%s", profile_type)
		httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)
	} else {
		httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, "/profiles", nil)
	}
	// httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}
	res := new(ProfileListResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) GetProfiles(ctx context.Context, engine string, profile_type string, id string, name string) (*ListProfileResponse, error) {
	var httpReq *http.Request
	var err error
	path := makePathProfiles(engine, profile_type, id, name)

	httpReq, err = sc.c.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, err
	}
	res := new(ListProfileResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListClusters(ctx context.Context) (*ListClusterResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/clusters", nil)
	if err != nil {
		return nil, err
	}
	res := new(ListClusterResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func (sc ServiceClient) ListSLA(ctx context.Context) (*SLAResponse, error) {
	httpReq, err := sc.c.NewRequest(ctx, http.MethodGet, "/slas", nil)
	if err != nil {
		return nil, err
	}
	res := new(SLAResponse)

	return res, sc.c.Do(ctx, httpReq, res)
}

func makePathProfiles(engine string, ptype string, id string, name string) string {
	if engine != "" {
		path := "/profiles?engine=" + engine
		if ptype != "" {
			path = path + "&type=" + ptype
		}
		if id != "" {
			path = path + "&id=" + id
		}
		if name != "" {
			path = path + "&name=" + name
		}
		return path
	}
	if ptype != "" {
		path := "/profiles?type=" + ptype
		if id != "" {
			path = path + "&id=" + id
		}
		if name != "" {
			path = path + "&name=" + name
		}
		return path
	}

	if id != "" {
		path := "/profiles?id=" + id
		if name != "" {
			path = path + "&name=" + name
		}
		return path
	}

	if name != "" {
		path := "/profiles?name=" + name
		return path
	}
	return ""
}
