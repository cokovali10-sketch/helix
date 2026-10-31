// Package api provides gRPC service definitions for the Helix platform
package api

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventServiceServer is the gRPC server for event operations
type EventServiceServer struct {
	UnimplementedEventServiceServer
}

// Ingest receives events from agents
func (s *EventServiceServer) Ingest(ctx context.Context, req *EventRequest) (*EventResponse, error) {
	return &EventResponse{
		Success: true,
		Message: "Event ingested successfully",
		Id:      fmt.Sprintf("evt-%d", time.Now().UnixNano()),
	}, nil
}

// IngestBatch receives a batch of events
func (s *EventServiceServer) IngestBatch(ctx context.Context, req *EventBatchRequest) (*EventBatchResponse, error) {
	ids := make([]string, len(req.Events))
	for i := range req.Events {
		ids[i] = fmt.Sprintf("evt-%d", time.Now().UnixNano()+int64(i))
	}
	return &EventBatchResponse{
		Success: true,
		Message: fmt.Sprintf("Ingested %d events", len(req.Events)),
		Ids:     ids,
	}, nil
}

// Query executes a query against stored events
func (s *EventServiceServer) Query(ctx context.Context, req *QueryRequest) (*QueryResponse, error) {
	start := time.Now()
	
	results := []*QueryResult{
		{
			Metric:    map[string]string{"host": "server1"},
			Value:     50.0,
			Timestamp: time.Now().Unix(),
		},
	}
	
	return &QueryResponse{
		Results:       results,
		Error:         "",
		ExecutionTime: time.Since(start).Seconds(),
	}, nil
}

// StreamEvents streams events in real-time
func (s *EventServiceServer) StreamEvents(req *StreamEventsRequest, stream EventService_StreamEventsServer) error {
	for {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			event := &EventResponse{
				Success: true,
				Message: "streaming event",
				Id:      fmt.Sprintf("evt-%d", time.Now().UnixNano()),
			}
			if err := stream.Send(event); err != nil {
				return err
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// TenantServiceServer is the gRPC server for tenant operations
type TenantServiceServer struct {
	UnimplementedTenantServiceServer
}

// CreateTenant creates a new tenant
func (s *TenantServiceServer) CreateTenant(ctx context.Context, req *CreateTenantRequest) (*TenantResponse, error) {
	return &TenantResponse{
		Id:        fmt.Sprintf("tenant-%d", time.Now().UnixNano()),
		Name:      req.Name,
		RateLimit: req.RateLimit,
		Quota:     req.Quota,
	}, nil
}

// GetTenant gets a tenant by ID
func (s *TenantServiceServer) GetTenant(ctx context.Context, req *GetTenantRequest) (*TenantResponse, error) {
	return &TenantResponse{
		Id:        req.Id,
		Name:      "Test Tenant",
		RateLimit: 1000,
		Quota:     10000,
	}, nil
}

// ListTenants lists all tenants
func (s *TenantServiceServer) ListTenants(ctx context.Context, req *ListTenantsRequest) (*ListTenantsResponse, error) {
	tenants := []*TenantResponse{
		{Id: "tenant-1", Name: "SOC Team Alpha", RateLimit: 5000, Quota: 100000},
		{Id: "tenant-2", Name: "SOC Team Beta", RateLimit: 3000, Quota: 50000},
	}
	return &ListTenantsResponse{Tenants: tenants}, nil
}

// UpdateTenant updates a tenant
func (s *TenantServiceServer) UpdateTenant(ctx context.Context, req *UpdateTenantRequest) (*TenantResponse, error) {
	return &TenantResponse{
		Id:        req.Id,
		Name:      req.Name,
		RateLimit: req.RateLimit,
		Quota:     req.Quota,
	}, nil
}

// DeleteTenant deletes a tenant
func (s *TenantServiceServer) DeleteTenant(ctx context.Context, req *DeleteTenantRequest) (*DeleteTenantResponse, error) {
	return &DeleteTenantResponse{Success: true}, nil
}

// DashboardServiceServer is the gRPC server for dashboard operations
type DashboardServiceServer struct {
	UnimplementedDashboardServiceServer
}

// CreateDashboard creates a new dashboard
func (s *DashboardServiceServer) CreateDashboard(ctx context.Context, req *CreateDashboardRequest) (*DashboardResponse, error) {
	return &DashboardResponse{
		Id:          fmt.Sprintf("dashboard-%d", time.Now().UnixNano()),
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

// GetDashboard gets a dashboard by ID
func (s *DashboardServiceServer) GetDashboard(ctx context.Context, req *GetDashboardRequest) (*DashboardResponse, error) {
	return &DashboardResponse{
		Id:          req.Id,
		Title:       "Test Dashboard",
		Description: "A test dashboard",
	}, nil
}

// ListDashboards lists all dashboards
func (s *DashboardServiceServer) ListDashboards(ctx context.Context, req *ListDashboardsRequest) (*ListDashboardsResponse, error) {
	dashboards := []*DashboardResponse{
		{Id: "dashboard-1", Title: "SOC Overview", Description: "Overall SOC metrics"},
		{Id: "dashboard-2", Title: "Network Traffic", Description: "Network analysis"},
	}
	return &ListDashboardsResponse{Dashboards: dashboards}, nil
}

// UpdateDashboard updates a dashboard
func (s *DashboardServiceServer) UpdateDashboard(ctx context.Context, req *UpdateDashboardRequest) (*DashboardResponse, error) {
	return &DashboardResponse{
		Id:          req.Id,
		Title:       req.Title,
		Description: req.Description,
	}, nil
}

// DeleteDashboard deletes a dashboard
func (s *DashboardServiceServer) DeleteDashboard(ctx context.Context, req *DeleteDashboardRequest) (*DeleteDashboardResponse, error) {
	return &DeleteDashboardResponse{Success: true}, nil
}

// AlertServiceServer is the gRPC server for alert operations
type AlertServiceServer struct {
	UnimplementedAlertServiceServer
}

// CreateAlertRule creates a new alert rule
func (s *AlertServiceServer) CreateAlertRule(ctx context.Context, req *CreateAlertRuleRequest) (*AlertRuleResponse, error) {
	return &AlertRuleResponse{
		Id:         fmt.Sprintf("rule-%d", time.Now().UnixNano()),
		Name:       req.Name,
		Expression: req.Expression,
		Severity:   req.Severity,
	}, nil
}

// ListAlertRules lists all alert rules
func (s *AlertServiceServer) ListAlertRules(ctx context.Context, req *ListAlertRulesRequest) (*ListAlertRulesResponse, error) {
	rules := []*AlertRuleResponse{
		{Id: "rule-1", Name: "HighCPU", Expression: "cpu > 90", Severity: "warning"},
		{Id: "rule-2", Name: "HighMemory", Expression: "memory > 90", Severity: "warning"},
	}
	return &ListAlertRulesResponse{Rules: rules}, nil
}

// GetAlerts gets all alerts
func (s *AlertServiceServer) GetAlerts(ctx context.Context, req *GetAlertsRequest) (*GetAlertsResponse, error) {
	alerts := []*AlertResponse{
		{Id: "alert-1", Name: "HighCPU", Status: "firing", Severity: "warning"},
		{Id: "alert-2", Name: "HighMemory", Status: "resolved", Severity: "warning"},
	}
	return &GetAlertsResponse{Alerts: alerts}, nil
}

// RegisterServices registers all gRPC services
func RegisterServices(grpcServer *grpc.Server) {
	// In production, would register all service implementations
	// grpcServer.RegisterService(&EventService_ServiceDesc, &EventServiceServer{})
	// grpcServer.RegisterService(&TenantService_ServiceDesc, &TenantServiceServer{})
	// grpcServer.RegisterService(&DashboardService_ServiceDesc, &DashboardServiceServer{})
	// grpcServer.RegisterService(&AlertService_ServiceDesc, &AlertServiceServer{})
}

// gRPC message types (simplified - in production would use protobuf)

type EventRequest struct {
	Timestamp int64
	Source    string
	Metrics   map[string]float64
	Tags      map[string]string
	Tenant    string
}

type EventResponse struct {
	Success bool
	Message string
	Id      string
}

type EventBatchRequest struct {
	Events []*EventRequest
}

type EventBatchResponse struct {
	Success bool
	Message string
	Ids     []string
}

type QueryRequest struct {
	Expression string
	TimeRange  []int64
	Tenant     string
	Step       int64
}

type QueryResult struct {
	Metric     map[string]string
	Value      float64
	Timestamp  int64
}

type QueryResponse struct {
	Results       []*QueryResult
	Error         string
	ExecutionTime float64
}

type StreamEventsRequest struct {
	Tenant string
	Filter string
}

type CreateTenantRequest struct {
	Name      string
	RateLimit int
	Quota     int64
}

type GetTenantRequest struct {
	Id string
}

type ListTenantsRequest struct {
	Limit  int
	Offset int
}

type UpdateTenantRequest struct {
	Id        string
	Name      string
	RateLimit int
	Quota     int64
}

type DeleteTenantRequest struct {
	Id string
}

type TenantResponse struct {
	Id        string
	Name      string
	RateLimit int
	Quota     int64
}

type ListTenantsResponse struct {
	Tenants []*TenantResponse
}

type DeleteTenantResponse struct {
	Success bool
}

type CreateDashboardRequest struct {
	Title       string
	Description string
	Panels      []*PanelRequest
}

type PanelRequest struct {
	Title  string
	Type   string
	Query  string
	Width  int
	Height int
}

type GetDashboardRequest struct {
	Id string
}

type ListDashboardsRequest struct {
	Limit  int
	Offset int
	Owner  string
}

type UpdateDashboardRequest struct {
	Id          string
	Title       string
	Description string
}

type DeleteDashboardRequest struct {
	Id string
}

type DashboardResponse struct {
	Id          string
	Title       string
	Description string
}

type ListDashboardsResponse struct {
	Dashboards []*DashboardResponse
}

type DeleteDashboardResponse struct {
	Success bool
}

type CreateAlertRuleRequest struct {
	Name        string
	Expression  string
	Duration    int64
	Severity    string
	Labels      map[string]string
	Annotations map[string]string
}

type ListAlertRulesRequest struct {
	Limit  int
	Offset int
}

type GetAlertsRequest struct {
	Status   string
	Severity string
	Limit    int
	Offset   int
}

type AlertRuleResponse struct {
	Id         string
	Name       string
	Expression string
	Severity   string
}

type ListAlertRulesResponse struct {
	Rules []*AlertRuleResponse
}

type AlertResponse struct {
	Id       string
	Name     string
	Status   string
	Severity string
}

type GetAlertsResponse struct {
	Alerts []*AlertResponse
}

// EventService_StreamEventsServer is the streaming interface
type EventService_StreamEventsServer interface {
	Send(*EventResponse) error
	Recv() (*StreamEventsRequest, error)
	Context() context.Context
}

// UnimplementedEventServiceServer is the base implementation
type UnimplementedEventServiceServer struct{}

// UnimplementedTenantServiceServer is the base implementation
type UnimplementedTenantServiceServer struct{}

// UnimplementedDashboardServiceServer is the base implementation
type UnimplementedDashboardServiceServer struct{}

// UnimplementedAlertServiceServer is the base implementation
type UnimplementedAlertServiceServer struct{}

// EventService_ServiceDesc is the service descriptor
var EventService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "helix.EventService",
	HandlerType:  (*EventServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
}

// TenantService_ServiceDesc is the service descriptor
var TenantService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "helix.TenantService",
	HandlerType: (*TenantServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
}

// DashboardService_ServiceDesc is the service descriptor
var DashboardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "helix.DashboardService",
	HandlerType: (*DashboardServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
}

// AlertService_ServiceDesc is the service descriptor
var AlertService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "helix.AlertService",
	HandlerType: (*AlertServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams:     []grpc.StreamDesc{},
}

// Helper for timestamp conversion
func toTimestamp(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func fromTimestamp(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}