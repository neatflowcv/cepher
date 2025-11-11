package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/neatflowcv/cepher/api"
	"github.com/neatflowcv/cepher/internal/app/flow"
)

var _ api.StrictServerInterface = (*Handler)(nil)

type Handler struct {
	service       *flow.Service
	scheduler     gocron.Scheduler
	jobSliders    map[string]*Slider
	levelDuration map[int]time.Duration
}

func NewHandler(service *flow.Service) (*Handler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	const (
		stableDuration = 6 * time.Minute
		warnDuration   = 3 * time.Minute
		errDuration    = 1 * time.Minute
	)

	handler := &Handler{
		service:    service,
		scheduler:  scheduler,
		jobSliders: make(map[string]*Slider),
		levelDuration: map[int]time.Duration{
			0: stableDuration,
			1: warnDuration,
			2: errDuration,
		},
	}

	clusters, err := service.ListClusters(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}

	for _, cluster := range clusters {
		handler.addJob(cluster.ID)
	}

	scheduler.Start()

	return handler, nil
}

func (h *Handler) Get() http.Handler {
	return api.HandlerFromMux(api.NewStrictHandler(h, nil), chi.NewMux())
}

func (h *Handler) Close() {
	err := h.scheduler.Shutdown()
	if err != nil {
		log.Printf("failed to shutdown scheduler: %v", err)
	}
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

	h.addJob(cluster.ID) //nolint:contextcheck

	return api.RegisterCluster201JSONResponse{
		Id:       cluster.ID,
		Name:     cluster.Name,
		Status:   api.ClusterStatus(cluster.Status),
		IsStable: cluster.IsStable,
		Detail:   &cluster.Detail,
	}, nil
}

func (h *Handler) ListClusters(
	ctx context.Context,
	request api.ListClustersRequestObject,
) (api.ListClustersResponseObject, error) {
	log.Println("ListClusters")

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
			Detail:   &cluster.Detail,
		})
	}

	return api.ListClusters200JSONResponse(apiClusters), nil
}

func (h *Handler) refreshCluster(clusterID string) {
	now := time.Now()
	log.Printf("refreshCluster: %s at %v", clusterID, now)

	ok, err := h.service.RefreshCluster(context.Background(), clusterID, now)
	if err != nil {
		log.Printf("failed to refresh cluster %s: %v", clusterID, err)

		return
	}

	slider := h.jobSliders[clusterID]
	if ok {
		slider = slider.Down()
	} else {
		slider = slider.Up()
	}

	h.jobSliders[clusterID] = slider
}

func (h *Handler) afterJobRuns(jobID uuid.UUID, clusterID string) {
	slider := h.jobSliders[clusterID]

	_, err := h.scheduler.NewJob(
		gocron.OneTimeJob(gocron.OneTimeJobStartDateTime(time.Now().Add(h.levelDuration[slider.value]))),
		gocron.NewTask(h.refreshCluster, clusterID),
		gocron.WithName(clusterID),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(h.afterJobRuns),
		),
	)
	if err != nil {
		log.Printf("failed to create job: %v", err)
	}
}

func (h *Handler) addJob(clusterID string) {
	const maxValue = 2

	h.jobSliders[clusterID] = NewSlider(0, maxValue, maxValue)
	h.afterJobRuns(uuid.Nil, clusterID)

	_, err := h.scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func() {
			log.Printf("UpdateMonitor at %v", time.Now())

			err := h.service.UpdateMonitor(context.Background(), clusterID)
			if err != nil {
				log.Printf("failed to update monitor %s: %v", clusterID, err)
			}
		}),
		gocron.JobOption(gocron.WithStartImmediately()),
	)
	if err != nil {
		log.Printf("failed to create job: %v", err)
	}
}
