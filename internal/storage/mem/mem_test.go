package mem

import "testing"

func TestStorage_Create(t *testing.T) {
	type fields struct {
		m map[string][]string
	}
	type args struct {
		key string
		val string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create",
			fields: fields{
				map[string][]string{},
			},
			args: args{
				key: "k",
				val: "v",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				data: tt.fields.m,
			}
			if err := s.Create(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_Update(t *testing.T) {
	type fields struct {
		m map[string][]string
	}
	type args struct {
		key string
		val []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update exists",
			fields: fields{
				map[string][]string{
					"k": {"a", "b"},
				},
			},
			args: args{
				key: "k",
				val: []string{"c"},
			},
			wantErr: false,
		},
		{
			name: "update not exists",
			fields: fields{
				map[string][]string{},
			},
			args: args{
				key: "k",
				val: []string{"c"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				data: tt.fields.m,
			}
			if err := s.Update(tt.args.key, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
