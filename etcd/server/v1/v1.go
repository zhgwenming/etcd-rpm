package v1

import (
	"net/http"

	"github.com/coreos/etcd/store"
	"github.com/coreos/etcd/third_party/github.com/goraft/raft"
)

// The Server interface provides all the methods required for the v1 API.
type Server interface {
	CommitIndex() uint64
	Term() uint64
	Store() store.Store
	Dispatch(raft.Command, http.ResponseWriter, *http.Request) error
}
