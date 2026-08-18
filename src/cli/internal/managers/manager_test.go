package managers

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

// errStub is a service that fails whichever lifecycle phase it is told to fail,
// so the manager's error propagation can be exercised without a real service.
type errStub struct {
	sm         ServiceManager
	failStart  bool
	failClose  bool
	failHealth bool
}

var errStubFailed = errors.New("stub failed")

func (e *errStub) Start(_ context.Context) error {
	if e.failStart {
		return errStubFailed
	}
	return nil
}

func (e *errStub) Close(_ context.Context) error {
	if e.failClose {
		return errStubFailed
	}
	return nil
}

func (e *errStub) Healthy(_ context.Context) error {
	if e.failHealth {
		return errStubFailed
	}
	return nil
}

func (e *errStub) ServiceManager() ServiceManager { return e.sm }

// errLogsService fails as a LogsService.
type errLogsService struct{ errStub }

func (e *errLogsService) WithServiceManager(sm ServiceManager) LogsService {
	e.sm = sm
	return e
}
func (e *errLogsService) Info(_ context.Context, _ string)  {}
func (e *errLogsService) Warn(_ context.Context, _ string)  {}
func (e *errLogsService) Error(_ context.Context, _ string) {}
func (e *errLogsService) Debug(_ context.Context, _ string) {}

// errCommandService fails as a CommandService.
type errCommandService struct{ errStub }

func (e *errCommandService) WithServiceManager(sm ServiceManager) CommandService {
	e.sm = sm
	return e
}
func (e *errCommandService) Execute(_ context.Context, _ []string) error { return nil }

func TestNewManager(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "success"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm := NewManager()
			if sm == nil {
				t.Fatal("NewManager() = nil")
			}
			// Every slot must be usable before anything is registered.
			if sm.LogsService() == nil {
				t.Error("LogsService() = nil, want a no-op")
			}
			if sm.CommandService() == nil {
				t.Error("CommandService() = nil, want a no-op")
			}
		})
	}
}

// TestNewManager_noOpsAreCallable is the property the whole container rests on:
// an unregistered service can be called without a nil check.
func TestNewManager_noOpsAreCallable(t *testing.T) {
	sm := NewManager()
	ctx := context.Background()

	sm.LogsService().Info(ctx, "no-op")
	sm.LogsService().Warn(ctx, "no-op")
	sm.LogsService().Error(ctx, "no-op")
	sm.LogsService().Debug(ctx, "no-op")

	if err := sm.CommandService().Execute(ctx, []string{"anything"}); err != nil {
		t.Errorf("no-op CommandService.Execute() error = %v, want nil", err)
	}
}

func Test_serviceManagerFinal_WithLogsService(t *testing.T) {
	type args struct{ s LogsService }
	sm := NewManager()
	s := NewNoOpsLogsService()
	tests := []struct {
		name string
		m    *serviceManagerFinal
		args args
		want ServiceManager
	}{
		{name: "success", m: sm.(*serviceManagerFinal), args: args{s}, want: sm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.WithLogsService(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithLogsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serviceManagerFinal_LogsService(t *testing.T) {
	sm := NewManager()
	s := NewNoOpsLogsService()
	sm.WithLogsService(s)
	tests := []struct {
		name string
		m    *serviceManagerFinal
		want LogsService
	}{
		{name: "success", m: sm.(*serviceManagerFinal), want: s},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.LogsService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogsService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serviceManagerFinal_WithCommandService(t *testing.T) {
	type args struct{ s CommandService }
	sm := NewManager()
	s := NewNoOpsCommandService()
	tests := []struct {
		name string
		m    *serviceManagerFinal
		args args
		want ServiceManager
	}{
		{name: "success", m: sm.(*serviceManagerFinal), args: args{s}, want: sm},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.WithCommandService(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCommandService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serviceManagerFinal_CommandService(t *testing.T) {
	sm := NewManager()
	s := NewNoOpsCommandService()
	sm.WithCommandService(s)
	tests := []struct {
		name string
		m    *serviceManagerFinal
		want CommandService
	}{
		{name: "success", m: sm.(*serviceManagerFinal), want: s},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.CommandService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommandService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_serviceManagerFinal_registrationInjectsManager verifies the injection
// half of the container: registering a service hands it the manager, which is
// how it later reaches its peers.
func Test_serviceManagerFinal_registrationInjectsManager(t *testing.T) {
	sm := NewManager()
	logs := &errLogsService{}
	cmd := &errCommandService{}

	sm.WithLogsService(logs).WithCommandService(cmd)

	if logs.ServiceManager() != sm {
		t.Error("WithLogsService() did not inject the manager")
	}
	if cmd.ServiceManager() != sm {
		t.Error("WithCommandService() did not inject the manager")
	}
}

func Test_serviceManagerFinal_Start(t *testing.T) {
	tests := []struct {
		name    string
		build   func() ServiceManager
		wantErr bool
	}{
		{
			name:    "success with no-ops",
			build:   NewManager,
			wantErr: false,
		},
		{
			name: "logs service failure aborts",
			build: func() ServiceManager {
				return NewManager().WithLogsService(&errLogsService{errStub{failStart: true}})
			},
			wantErr: true,
		},
		{
			name: "command service failure aborts",
			build: func() ServiceManager {
				return NewManager().WithCommandService(&errCommandService{errStub{failStart: true}})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.build().Start(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test_serviceManagerFinal_StartDoesNotBlock guards the CLI's core requirement:
// Start returns so the command can run, rather than blocking like a server.
func Test_serviceManagerFinal_StartDoesNotBlock(t *testing.T) {
	done := make(chan error, 1)
	go func() {
		done <- NewManager().Start(context.Background())
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Start() error = %v, want nil", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Start() blocked; it must return so the command can execute")
	}
}

func Test_serviceManagerFinal_Close(t *testing.T) {
	tests := []struct {
		name    string
		build   func() ServiceManager
		wantErr bool
	}{
		{
			name:    "success with no-ops",
			build:   NewManager,
			wantErr: false,
		},
		{
			name: "command service failure aborts",
			build: func() ServiceManager {
				return NewManager().WithCommandService(&errCommandService{errStub{failClose: true}})
			},
			wantErr: true,
		},
		{
			name: "logs service failure aborts",
			build: func() ServiceManager {
				return NewManager().WithLogsService(&errLogsService{errStub{failClose: true}})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.build().Close(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_serviceManagerFinal_Healthy(t *testing.T) {
	tests := []struct {
		name    string
		build   func() ServiceManager
		wantErr bool
	}{
		{
			name:    "success with no-ops",
			build:   NewManager,
			wantErr: false,
		},
		{
			name: "logs service unhealthy",
			build: func() ServiceManager {
				return NewManager().WithLogsService(&errLogsService{errStub{failHealth: true}})
			},
			wantErr: true,
		},
		{
			name: "command service unhealthy",
			build: func() ServiceManager {
				return NewManager().WithCommandService(&errCommandService{errStub{failHealth: true}})
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.build().Healthy(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("Healthy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
