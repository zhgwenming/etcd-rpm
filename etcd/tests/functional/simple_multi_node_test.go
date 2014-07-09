package test

import (
	"os"
	"testing"
	"time"

	"github.com/coreos/etcd/third_party/github.com/coreos/go-etcd/etcd"
)

func TestSimpleMultiNode(t *testing.T) {
	templateTestSimpleMultiNode(t, false)
}

func TestSimpleMultiNodeTls(t *testing.T) {
	templateTestSimpleMultiNode(t, true)
}

// Create a three nodes and try to set value
func templateTestSimpleMultiNode(t *testing.T, tls bool) {
	procAttr := new(os.ProcAttr)
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}

	clusterSize := 3

	_, etcds, err := CreateCluster(clusterSize, procAttr, tls)

	if err != nil {
		t.Fatalf("cannot create cluster: %v", err)
	}

	defer DestroyCluster(etcds)

	time.Sleep(time.Second)

	c := etcd.NewClient(nil)

	if c.SyncCluster() == false {
		t.Fatal("Cannot sync cluster!")
	}

	// Test Set
	result, err := c.Set("foo", "bar", 100)
	if err != nil {
		t.Fatal(err)
	}

	node := result.Node
	if node.Key != "/foo" || node.Value != "bar" || node.TTL < 95 {
		t.Fatalf("Set 1 failed with %s %s %v", node.Key, node.Value, node.TTL)
	}

	time.Sleep(time.Second)

	result, err = c.Set("foo", "bar", 100)
	if err != nil {
		t.Fatal(err)
	}

	node = result.Node
	if node.Key != "/foo" || node.Value != "bar" || node.TTL < 95 {
		t.Fatalf("Set 2 failed with %s %s %v", node.Key, node.Value, node.TTL)
	}

}
