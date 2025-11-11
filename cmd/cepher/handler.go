package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/neatflowcv/cepher/api"
	"github.com/neatflowcv/cepher/internal/app/flow"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	service *flow.Service
}

func NewHandler(service *flow.Service) http.Handler { //nolint:ireturn
	handler := &Handler{
		service: service,
	}

	return api.HandlerFromMux(api.NewStrictHandler(handler, nil), chi.NewMux())
}

func (h *Handler) RegisterCluster(
	ctx context.Context,
	request api.RegisterClusterRequestObject,
) (api.RegisterClusterResponseObject, error) {
	log.Println("RegisterCluster")

	cluster, err := h.service.RegisterCluster(ctx, &flow.RegisterCluster{
		Name:  request.Body.Name,
		Hosts: request.Body.Hosts,
		Key:   request.Body.Key,
		Now:   time.Now(),
	})
	if err != nil {
		return api.RegisterCluster500JSONResponse{ //nolint:nilerr
			Message: err.Error(),
		}, nil
	}

	return api.RegisterCluster201JSONResponse{
		Id:       cluster.ID,
		Name:     cluster.Name,
		Status:   api.ClusterStatus(cluster.Status),
		IsStable: cluster.IsStable,
	}, nil
}

func (h *Handler) ListClusters(
	ctx context.Context,
	request api.ListClustersRequestObject,
) (api.ListClustersResponseObject, error) {
	clusters, err := h.service.ListClusters(ctx)
	if err != nil {
		return api.ListClusters500JSONResponse{ //nolint:nilerr
			Message: err.Error(),
		}, nil
	}

	if len(clusters) == 0 {
		return api.ListClusters204Response{}, nil
	}

	var apiClusters []api.Cluster
	for _, cluster := range clusters {
		apiClusters = append(apiClusters, api.Cluster{
			Id:       cluster.ID,
			Name:     cluster.Name,
			Status:   api.ClusterStatus(cluster.Status),
			IsStable: cluster.IsStable,
		})
	}

	return api.ListClusters200JSONResponse(apiClusters), nil
}
