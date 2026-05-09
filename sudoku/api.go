package sudoku

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/fbundle/http_transport/http_transport"
)

// N is the sudoku block size (N*N x N*N board).
const N = 3

// Store is a key-value store for game sessions.
type Store interface {
	Set(key string, value any)
	Get(key string) any
	NumActiveKey() int
}

type keyView struct {
	Key string `json:"key"`
}
type boardView struct {
	Board string `json:"board"`
}
type placeView struct {
	Key string `json:"key"`
	Val int    `json:"val"`
}
type pointView struct {
	Key string `json:"key"`
	Row int    `json:"row"`
	Col int    `json:"col"`
}

func decode[T any](body []byte) (T, bool) {
	var v T
	return v, json.Unmarshal(body, &v) == nil
}

func encode(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func getGame(store Store, key string) Game {
	v := store.Get(key)
	if v == nil {
		return nil
	}
	return v.(Game)
}

// RegisterRoutes registers all API routes on the given Router.
func RegisterRoutes(r http_transport.Router, store Store, rnd *rand.Rand) {
	r.POST("api/new", func(body []byte) (int, []byte) {
		bv, _ := decode[boardView](body)
		board, ok := FromString(N, bv.Board)
		intKey := rnd.Int()
		if !ok {
			board = Generate(N, intKey)
		}
		key := keyView{Key: strconv.Itoa(intKey)}
		game, ok := NewGame(N, board)
		if !ok {
			return http.StatusBadRequest, nil
		}
		store.Set(key.Key, game)
		return http.StatusOK, encode(key)
	})

	r.POST("api/login", func(body []byte) (int, []byte) {
		kv, ok := decode[keyView](body)
		if ok && store.Get(kv.Key) != nil {
			return http.StatusOK, encode(kv)
		}
		return http.StatusBadRequest, nil
	})

	r.POST("api/view", func(body []byte) (int, []byte) {
		kv, ok := decode[keyView](body)
		if !ok {
			return http.StatusBadRequest, nil
		}
		game := getGame(store, kv.Key)
		if game == nil {
			return http.StatusNotFound, nil
		}
		return http.StatusOK, encode(game.View())
	})

	r.POST("api/point", func(body []byte) (int, []byte) {
		pv, ok := decode[pointView](body)
		if !ok {
			return http.StatusBadRequest, nil
		}
		game := getGame(store, pv.Key)
		if game == nil {
			return http.StatusNotFound, nil
		}
		game.Point(CellView{Row: pv.Row, Col: pv.Col})
		return http.StatusOK, nil
	})

	r.POST("api/place", func(body []byte) (int, []byte) {
		pv, ok := decode[placeView](body)
		if !ok {
			return http.StatusBadRequest, nil
		}
		game := getGame(store, pv.Key)
		if game == nil {
			return http.StatusNotFound, nil
		}
		game.Place(pv.Val)
		return http.StatusOK, nil
	})

	r.POST("api/undo", func(body []byte) (int, []byte) {
		kv, ok := decode[keyView](body)
		if !ok {
			return http.StatusBadRequest, nil
		}
		game := getGame(store, kv.Key)
		if game == nil {
			return http.StatusNotFound, nil
		}
		game.Undo()
		return http.StatusOK, nil
	})

	r.POST("api/implication", func(body []byte) (int, []byte) {
		kv, ok := decode[keyView](body)
		if !ok {
			return http.StatusBadRequest, nil
		}
		game := getGame(store, kv.Key)
		if game == nil {
			return http.StatusNotFound, nil
		}
		game.Implication()
		return http.StatusOK, nil
	})

	r.POST("api/access", func(body []byte) (int, []byte) {
		kv, _ := decode[keyView](body)
		store.Get(kv.Key)
		return http.StatusOK, nil
	})

	r.POST("api/global_stats", func(body []byte) (int, []byte) {
		return http.StatusOK, encode(map[string]any{
			"number of active users": store.NumActiveKey(),
		})
	})
}
