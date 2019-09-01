package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// export GRAPH_DB_NAME="arango_graphs" # name of database in arango
// export GRAPH_DB_COLLECTION_NAME="wikipedia" # collection name within arango db name
// export GRAPH_DB_ARANGO_ENDPOINTS="http://localhost:9520" #list of arango db endpoints
// export GRAPH_DB_NAME="wikipedia-graph" # name of graph within collection
func TestConnectToDB(t *testing.T) {
	// mock out log.Fatalf
	origLogFatalf := logFatalf
	defer func() { logFatalf = origLogFatalf }()
	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}
	t.Run("connects to db that doesnt already exist and connects to graph that does exist", func(t *testing.T) {
		errors = []string{}
		os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
		os.Setenv("GRAPH_DB_NAME", "wikipedia-graph-1")
		g, nodes, edges := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		// connect to graph we just created
		g, nodes, edges = ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		require.Nil(t, g.Remove(nil))
	})
	t.Run("connects to same DB with new graph name", func(t *testing.T) {
		os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
		os.Setenv("GRAPH_DB_NAME", "wikipedia-graph-2")
		errors = []string{}
		dbName2 := "graph-testing-2"
		os.Setenv("GRAPH_DB_NAME", dbName2)
		g, nodes, edges := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		require.Nil(t, g.Remove(nil))
	})
	t.Run("bad url endpoints", func(t *testing.T) {
		errors = []string{}
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8000")
		g, nodes, edges := ConnectToDB()
		assert.Nil(t, nodes)
		assert.Nil(t, edges)
		assert.Nil(t, g)
		assert.Equal(t, []string{"Could not establish connection to DB [Could not check if databse exists create database at [http://localhost:8000]: Get http://localhost:8000/_db/graph-testing-2/_api/database/current: dial tcp 127.0.0.1:8000: connect: connection refused]"}, errors)
		errors = []string{}
	})
	t.Run("bad db name", func(t *testing.T) {
		errors = []string{}
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
		os.Setenv("GRAPH_DB_NAME", "sldjf093ur2n093r2039d[2e9ufsdf - -CC]")
		g, nodes, edges := ConnectToDB()
		assert.Nil(t, nodes)
		assert.Nil(t, edges)
		assert.Nil(t, g)
		assert.Equal(t, []string{"Could not establish connection to DB [Failed to initialize database: database name invalid]"}, errors)
		errors = []string{}
	})
}

func TestAddEdgesDB(t *testing.T) {
	// mock out log.Fatalf
	origLogFatalf := logFatalf
	defer func() { logFatalf = origLogFatalf }()
	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}
	os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
	os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
	os.Setenv("GRAPH_DB_NAME", "wikipedia-graph-1")
	g, nodes, edges := ConnectToDB()
	assert.NotNil(t, g)
	assert.NotNil(t, nodes)
	assert.NotNil(t, edges)
	require.Equal(t, []string{}, errors)
	defer require.Nil(t, g.Remove(nil))

	type Test struct {
		Before             func()
		Name               string
		Node               string
		Neighbors          []string
		ExpectedError      error
		ExpectedNodesAdded []string
	}

	testTable := []Test{
		Test{
			Before:             func() {},
			Name:               "addes all new edges",
			Node:               "new-node-1",
			Neighbors:          []string{"new-node-2", "new-node-3"},
			ExpectedError:      nil,
			ExpectedNodesAdded: []string{"new-node-2", "new-node-3"},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Before()
			e, nAdded := AddEdges(test.Node, test.Neighbors)
			assert.Equal(t, test.ExpectedError, e)
			assert.Equal(t, test.ExpectedNodesAdded, nAdded)
		})
	}

}
