package tomasql

import (
	"container/heap"
	"fmt"
	"slices"
	"strings"
)

// DBGraphData represents a directed graph of database tables and their relationships.
// Note that this graph only contains 'forward' relationships, i.e. from the table that has a foreign key
// to the table that is referenced by the foreign key.
type DBGraphData struct {
	// indexed by source table, target table, and source column.
	relationships map[Table]map[Table]map[Column]Column
}

func newDBGraphData(relationships map[Table]map[Table]map[Column]Column) *DBGraphData {
	dup := deepCopyRelationships(relationships)
	return &DBGraphData{
		relationships: dup,
	}
}

func NewDBGraphData(relationships map[Table]map[Table]map[Column]Column) *DBGraphData {
	return newDBGraphData(relationships)
}

func (g *DBGraphData) getNeighbors(table Table) []Table {
	if table == nil {
		return nil
	}

	neighbors := make([]Table, 0, len(g.relationships[table]))
	for target := range g.relationships[table] {
		neighbors = append(neighbors, target)
	}
	return neighbors
}

func (g *DBGraphData) AddLink(sourceCol, targetCol Column) *DBGraphData {
	g.addLink(sourceCol, targetCol)
	g.addLink(targetCol, sourceCol) // Add the reverse link as well
	return g
}

func (g *DBGraphData) addLink(sourceCol, targetCol Column) *DBGraphData {
	sourceTable := sourceCol.Table()
	targetTable := targetCol.Table()
	if _, ok := g.relationships[sourceTable]; !ok {
		g.relationships[sourceTable] = make(map[Table]map[Column]Column)
	}
	if _, ok := g.relationships[sourceTable][targetTable]; !ok {
		g.relationships[sourceTable][targetTable] = make(map[Column]Column)
	}
	g.relationships[sourceTable][targetTable][sourceCol] = targetCol
	return g
}

func (g *DBGraphData) RemoveLink(sourceCol, targetCol Column) *DBGraphData {
	sourceTable := sourceCol.Table()
	targetTable := targetCol.Table()
	if _, ok := g.relationships[sourceTable]; !ok {
		return g
	}
	if _, ok := g.relationships[sourceTable][targetTable]; !ok {
		return g
	}
	delete(g.relationships[sourceTable][targetTable], sourceCol)
	return g
}

func (g *DBGraphData) RemoveTable(table Table) *DBGraphData {
	delete(g.relationships, table)
	// Remove all relationships that have this table as a target
	for _, targets := range g.relationships {
		delete(targets, table)
	}

	return g
}

type JoinItem struct {
	Target      Table
	OnCondition Condition
}

// MinimalJoins returns the minimal join paths from the `from` table to each of the `targets` tables,
// reusing tables where possible.
// Returns a map where the keys are target tables and the values are slices of JoinItems containing all the tables that
// need to be joined to reach the target table and the conditions for the joins.
func (g *DBGraphData) MinimalJoins(from Table, targets []Table) (joinItems []*JoinItem, err error) {
	if from == nil || len(targets) == 0 {
		return nil, fmt.Errorf("from and target tables must not be nil")
	}
	prev := g.dijkstra(from)

	included := map[Table]bool{from: true}

	for _, target := range targets {
		current := target
		targetJoinItems := make([]*JoinItem, 0)
		for current != from {
			if prev[current] == nil {
				// if we reached a table that has no previous table, it means we can't reach the target from the source
				tableNames := make([]string, 0, len(targetJoinItems))
				for _, item := range targetJoinItems {
					tableNames = append(tableNames, item.Target.TableName())
				}
				return nil, fmt.Errorf("target table %s is not reachable from source table %s. Last table in path: %s. Path: %s",
					target.TableName(), from.TableName(), current.TableName(), strings.Join(tableNames, "->"))
			}
			if _, ok := included[current]; ok {
				// if we already included this table in another join path we can skip it
				break
			}
			included[current] = true

			links := g.relationships[prev[current]][current]

			// this whole thing only works if there is exactly one link column between the two tables, for the future:
			// maybe we can take an additional parameter to specify which link column to use.
			sourceCol, targetCol, err := getSingleLinkColumn(links)
			if err != nil {
				return nil, fmt.Errorf("error getting single link column from %s to %s: %w", prev[current].TableName(), current.TableName(), err)
			}

			item := &JoinItem{
				Target:      current,
				OnCondition: sourceCol.Eq(targetCol),
			}
			targetJoinItems = append(targetJoinItems, item)

			current = prev[current]
		}
		if len(targetJoinItems) > 0 {
			// reverse the paths to have them in the correct order
			slices.Reverse(targetJoinItems)
			joinItems = append(joinItems, targetJoinItems...)
		}
	}

	for _, target := range targets {
		if _, ok := included[target]; !ok {
			// if we didn't include the target table in the join path, it means we couldn't reach it from the source table
			return nil, fmt.Errorf("target table %s is not reachable from source table %s", target.TableName(), from.TableName())
		}
	}

	return joinItems, nil
}

// getSingleLinkColumn returns the single link column and its condition from the links map or an error if there are
// multiple columns linking to the same table.
func getSingleLinkColumn(links map[Column]Column) (source, target Column, err error) {
	if len(links) != 1 {
		keys := make([]string, 0, len(links))
		for k := range links {
			keys = append(keys, k.Table().TableName()+"."+k.Name())
		}
		return nil, nil, fmt.Errorf(
			"expected exactly one link column, got %d: %s",
			len(links),
			strings.Join(keys, ", "),
		)
	}
	for sourceCol, targetCol := range links {
		return sourceCol, targetCol, nil
	}
	return nil, nil, fmt.Errorf("no link column found")
}

func (g *DBGraphData) dijkstra(source Table) (prev map[Table]Table) {
	const inf = int(^uint(0) >> 1) // Max int value

	dist := make(map[Table]int)
	prev = make(map[Table]Table)
	inHeap := make(map[Table]*minHeapItem)

	h := &tableMinHeap{}

	addToHeap := func(v Table) {
		if inHeap[v] != nil {
			return // Already in heap, skip
		}

		d := inf
		if v == source {
			d = 0
		}
		item := &minHeapItem{table: v, dist: d}
		dist[v] = d
		prev[v] = nil
		heap.Push(h, item)
		inHeap[v] = item
	}

	// Initialize distances and heap
	for v, rel := range g.relationships {
		for target := range rel {
			addToHeap(target)
		}
		addToHeap(v)
	}

	for h.Len() > 0 {
		uItem := heap.Pop(h).(*minHeapItem)
		u := uItem.table
		delete(inHeap, u)

		// For each neighbor v of u
		for _, v := range g.getNeighbors(u) {
			if _, ok := inHeap[v]; !ok {
				continue // v is not in the heap, skip
			}
			alt := dist[u] + 1 // Assuming all edges have weight 1
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
				if vItem, ok := inHeap[v]; ok {
					vItem.dist = alt
					heap.Fix(h, vItem.index)
				}
			}
		}
	}

	return prev
}

func deepCopyRelationships(relationships map[Table]map[Table]map[Column]Column) map[Table]map[Table]map[Column]Column {
	newRelationships := make(map[Table]map[Table]map[Column]Column)
	for source, targets := range relationships {
		newTargets := make(map[Table]map[Column]Column)
		for target, links := range targets {
			newLinks := make(map[Column]Column)
			for sourceCol, targetCol := range links {
				newLinks[sourceCol] = targetCol
			}
			newTargets[target] = newLinks
		}
		newRelationships[source] = newTargets
	}
	return newRelationships
}

// minHeapItem is a helper struct for the heap.
type minHeapItem struct {
	table Table
	dist  int
	index int // Needed for heap.Fix, not used here but good practice
}

// tableMinHeap implements heap.Interface for minHeapItem.
type tableMinHeap []*minHeapItem

func (h *tableMinHeap) Len() int           { return len(*h) }
func (h *tableMinHeap) Less(i, j int) bool { return (*h)[i].dist < (*h)[j].dist }
func (h *tableMinHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index = i
	(*h)[j].index = j
}

func (h *tableMinHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(*minHeapItem)
	item.index = n
	*h = append(*h, item)
}

func (h *tableMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*h = old[0 : n-1]
	return item
}
