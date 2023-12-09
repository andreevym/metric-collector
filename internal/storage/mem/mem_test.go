package mem

import (
	"testing"

	"github.com/andreevym/metric-collector/internal/storage"
)

func TestStorage_Create(t *testing.T) {
	i := int64(1)
	tests := []struct {
		name    string
		storage map[string]*storage.Metric
		m       *storage.Metric
		wantErr bool
	}{
		{
			name:    "create",
			storage: map[string]*storage.Metric{},
			m: &storage.Metric{
				ID:    "k",
				MType: storage.MTypeCounter,
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
			if err := s.Create(nil, tt.m); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	i := int64(1)
	tests := []struct {
		name    string
		storage map[string]*storage.Metric
		m       *storage.Metric
		wantErr bool
	}{
		{
			name: "update exists",
			storage: map[string]*storage.Metric{
				"k" + storage.MTypeCounter: {
					ID:    "k",
					MType: storage.MTypeCounter,
					Delta: &i,
				},
			},
			m: &storage.Metric{
				ID:    "k",
				MType: storage.MTypeCounter,
				Delta: &i,
			},
			wantErr: false,
		},
		{
			name:    "update not exists",
			storage: map[string]*storage.Metric{},
			m: &storage.Metric{
				ID:    "k",
				MType: storage.MTypeCounter,
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
			if err := s.Update(nil, tt.m); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
