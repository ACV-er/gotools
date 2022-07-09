package storage

import (
	"math/rand"
	"testing"
	"time"
)

func TestThreadSafeBitMap(t *testing.T) {
	type fields struct {
		size            uint64
		lockGranularity uint64
	}
	type getArgs struct {
		pos     uint64
		wantRet bool
	}
	type setArgs struct {
		pos   uint64
		value bool
	}
	type args struct {
		set []setArgs
		get []getArgs
	}
	tests := []struct {
		name   string
		race   bool
		fields fields
		args   args
	}{
		{
			name: "简单测试",
			race: false,
			fields: fields{
				size:            1000000,
				lockGranularity: 1000,
			},
			args: args{
				set: []setArgs{
					{15, true},
					{25, true},
					{80, true},
					{999999, true},
					{80, false},
				},
				get: []getArgs{
					{15, true},
					{25, true},
					{80, false},
					{999999, true},
				},
			},
		},
		{
			name: "简单测试2",
			race: true,
			fields: fields{
				size:            1000000,
				lockGranularity: 1000,
			},
			args: args{
				set: []setArgs{
					{80, true},
					{999999, true},
					{999998, true},
					{80, false},
				},
				get: []getArgs{
					{80, false},
					{999999, true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBitMap(tt.fields.size, tt.fields.lockGranularity)
			if tt.race {
				for _, set := range tt.args.set {
					set := set
					go func() {
						b.Set(set.pos, set.value)
					}()
				}
			} else {
				for _, set := range tt.args.set {
					set := set
					b.Set(set.pos, set.value)
				}
			}

			for _, get := range tt.args.get {
				ret := b.Get(get.pos)
				if !tt.race && ret != get.wantRet {
					t.Errorf("pos %v want %v got %v\n", get.pos, get.wantRet, ret)
				}
			}
		})
	}
}

func TestBitMap(t *testing.T) {
	type fields struct {
		size            uint64
		lockGranularity uint64
	}
	type getArgs struct {
		pos     uint64
		wantRet bool
	}
	type setArgs struct {
		pos   uint64
		value bool
	}
	type args struct {
		set []setArgs
		get []getArgs
	}
	tests := []struct {
		name   string
		race   bool
		fields fields
		args   args
	}{
		{
			name: "简单测试",
			race: false,
			fields: fields{
				size:            1000000,
				lockGranularity: 0,
			},
			args: args{
				set: []setArgs{
					{15, true},
					{25, true},
					{80, true},
					{999999, true},
					{80, false},
				},
				get: []getArgs{
					{15, true},
					{25, true},
					{80, false},
					{999999, true},
				},
			},
		},
		{
			name: "简单测试2",
			race: true,
			fields: fields{
				size:            1000000,
				lockGranularity: 0,
			},
			args: args{
				set: []setArgs{
					{80, true},
					{999999, true},
					{999998, true},
					{80, false},
				},
				get: []getArgs{
					{80, false},
					{999999, true},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBitMap(tt.fields.size, tt.fields.lockGranularity)
			if tt.race {
				for _, set := range tt.args.set {
					set := set
					go func() {
						b.Set(set.pos, set.value)
					}()
				}
			} else {
				for _, set := range tt.args.set {
					set := set
					b.Set(set.pos, set.value)
				}
			}

			for _, get := range tt.args.get {
				ret := b.Get(get.pos)
				if !tt.race && ret != get.wantRet {
					t.Errorf("pos %v want %v got %v\n", get.pos, get.wantRet, ret)
				}
			}
		})
	}
}

func generateu64(n int) []uint64 {
	rand.Seed(time.Now().UnixNano())
	nums := make([]uint64, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, uint64(rand.Int63n(10000000)))
	}
	return nums
}
func generate(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, 0)
	for i := 0; i < n; i++ {
		nums = append(nums, rand.Intn(10000000))
	}
	return nums
}

func BenchmarkThreadSafeBitMap(b *testing.B) {
	set_arr := generateu64(10000000)
	get_arr := generateu64(10000000)
	for n := 0; n < b.N; n++ {
		bm := NewBitMap(10000000, 1000)
		for _, set_val := range set_arr {
			bm.Set(set_val, true)
		}
		for _, get_val := range get_arr {
			bm.Get(get_val)
		}
	}
}

func BenchmarkBitMap(b *testing.B) {
	set_arr := generateu64(10000000)
	get_arr := generateu64(10000000)
	for n := 0; n < b.N; n++ {
		bm := NewBitMap(10000000, 0)
		for _, set_val := range set_arr {
			bm.Set(set_val, true)
		}
		for _, get_val := range get_arr {
			bm.Get(get_val)
		}
	}
}
