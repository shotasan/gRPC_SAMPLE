package handler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"grpc_sample/api/gen/api"
)

func init() {
	// パンケーキの仕上がりに影響するseedを初期化する。
	rand.Seed(time.Now().UnixNano())
}

type BakerHandler struct {
	report *report
}

type report struct {
	// 並列処理の関係（排他制御）
	// 複数人が同時に接続しても大丈夫にしておく
	sync.Mutex
	data map[api.Pancake_Menu]int
}

// BakerHandlerを初期化する関数
func NewBakerHandler() *BakerHandler {
	return &BakerHandler{
		report: &report{
			data: make(map[api.Pancake_Menu]int),
		},
	}
}

// Bakeメソッドの定義
// PancakeBakerServiceインターフェースで定義されたメソッドを実装する
// 指定されたメニューのパンケーキを焼いて、焼けたパンをResponseとして返す
func (h *BakerHandler) Bake(
	ctx context.Context,
	req *api.BakeRequest,
) (*api.BakeResponse, error) {

	// バリデーション
	// api.Pancake_UNKNOWN →　定数の呼び出し
	// req.Menu > api.Pancake_SPICY_CURRY　→　リクエストの値が既定値以外の場合
	if req.Menu == api.Pancake_UNKOWN || req.Menu > api.Pancake_SPICY_CURRY {
		return nil, status.Errorf(codes.InvalidArgument, "パンケーキを選んでください！")
	}

	now := time.Now()
	// func (*Mutex) Lock
	h.report.Lock()
	// パンケーキを焼いて、数を記録する
	h.report.data[req.Menu] = h.report.data[req.Menu] + 1
	// func (*Mutex) Unlock
	h.report.Unlock()

	// レスポンスの作成
	return &api.BakeResponse{
		// gen/apiで定義されるPancake型を使用する
		Pancake: &api.Pancake{
			Menu:           req.Menu,
			ChefName:       "gami",
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

// Reportメソッド
// 焼けたパンケーキの数を報告する
func (h *BakerHandler) Report(
	ctx context.Context,
	req *api.ReportRequest,
) (*api.ReportResponse, error) {
	// Report_BakeCount型のポインタスライス
	counts := make([]*api.Report_BakeCount, 0)

	h.report.Lock()

	for k, v := range h.report.data {
		counts = append(counts, &api.Report_BakeCount{
			Menu:  k,
			Count: int32(v),
		})
	}

	h.report.Unlock()

	return &api.ReportResponse{
		Report: &api.Report{
			BakeCounts: counts,
		},
	}, nil
}
