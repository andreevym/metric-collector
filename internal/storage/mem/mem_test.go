package mem

import (
	"context"
	"testing"

	"github.com/andreevym/metric-collector/internal/storage/store"
)

func TestStorage_Create(t *testing.T) {
	i := int64(1)
	tests := []struct {
		name    string
		storage map[string]*store.Metric
		m       *store.Metric
		wantErr bool
	}{
		{
			name:    "create",
			storage: map[string]*store.Metric{},
			m: &store.Metric{
				ID:    "k",
				MType: store.MTypeCounter,
				Delta: &i,
				Value: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				data: tt.storage,
			}
			if err := s.Create(context.TODO(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	i := int64(1)
	tests := []struct {
		name    string
		storage map[string]*store.Metric
		m       *store.Metric
		wantErr bool
	}{
		{
			name: "update exists",
			storage: map[string]*store.Metric{
				"k" + store.MTypeCounter: {
					ID:    "k",
					MType: store.MTypeCounter,
					Delta: &i,
				},
			},
			m: &store.Metric{
				ID:    "k",
				MType: store.MTypeCounter,
				Delta: &i,
			},
			wantErr: false,
		},
		{
			name:    "update not exists",
			storage: map[string]*store.Metric{},
			m: &store.Metric{
				ID:    "k",
				MType: store.MTypeCounter,
				Delta: &i,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				data: tt.storage,
			}
			if err := s.Update(context.TODO(), tt.m); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
