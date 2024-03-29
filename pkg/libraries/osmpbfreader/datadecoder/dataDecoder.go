package datadecoder

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/filter"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/getdata"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmproto"
	"google.golang.org/protobuf/proto"
	"time"
)

const typicalPrimitiveBlockFeatureCount = 8000
const degMultiplier = 1e-9

type DataDecoder interface {
	Decode(blob *osmproto.Blob, filter filter.Filter) ([]interface{}, error)
}

type impl struct {
	data   []interface{}
	filter filter.Filter
}

func NewDataDecoder() DataDecoder {
	return &impl{}
}

func (i *impl) Decode(blob *osmproto.Blob, filter filter.Filter) ([]interface{}, error) {
	i.data = make([]interface{}, 0, typicalPrimitiveBlockFeatureCount)
	i.filter = filter

	data, err := getdata.GetData(blob)
	if err != nil {
		return nil, fmt.Errorf("error getting data: %s", err.Error())
	}

	primitiveBlock := &osmproto.PrimitiveBlock{}
	if err := proto.Unmarshal(data, primitiveBlock); err != nil {
		return nil, fmt.Errorf("error unmarshalling data: %s", err.Error())
	}

	i.parsePrimitiveBlock(primitiveBlock)

	return i.data, nil
}

func (i *impl) parsePrimitiveBlock(block *osmproto.PrimitiveBlock) {
	for _, group := range block.GetPrimitivegroup() {
		i.parsePrimitiveGroup(block, group)
	}
}

func (i *impl) parsePrimitiveGroup(block *osmproto.PrimitiveBlock, group *osmproto.PrimitiveGroup) {
	i.parseNodes(block, group.GetNodes())
	i.parseDenseNodes(block, group.GetDense())
	i.parseWays(block, group.GetWays())
	i.parseRelations(block, group.GetRelations())
}

func (i *impl) parseNodes(block *osmproto.PrimitiveBlock, nodes []*osmproto.Node) {
	if i.filter != nil && i.filter.FilterNodes() {
		return
	}

	stringTable := block.GetStringtable().GetS()
	granularity := int64(block.GetGranularity())
	dateGranularity := int64(block.GetDateGranularity())

	latOffset := block.GetLatOffset()
	lonOffset := block.GetLonOffset()

	for _, node := range nodes {
		id := node.GetId()
		lat := node.GetLat()
		lon := node.GetLon()

		latitude := degMultiplier * float64(latOffset+(granularity*lat))
		longitude := degMultiplier * float64(lonOffset+(granularity*lon))

		tags := extractTags(stringTable, node.GetKeys(), node.GetVals())
		info := extractInfo(stringTable, node.GetInfo(), dateGranularity)

		i.data = append(i.data, osmpbfreaderdata.Node{
			ID:   id,
			Lat:  latitude,
			Lon:  longitude,
			Tags: tags,
			Info: info,
		})
	}

}

func (i *impl) parseDenseNodes(block *osmproto.PrimitiveBlock, denseNodes *osmproto.DenseNodes) {
	if i.filter != nil && i.filter.FilterNodes() {
		return
	}

	stringTable := block.GetStringtable().GetS()
	granularity := int64(block.GetGranularity())
	latOffset := block.GetLatOffset()
	lonOffset := block.GetLonOffset()
	dateGranularity := int64(block.GetDateGranularity())
	ids := denseNodes.GetId()
	lats := denseNodes.GetLat()
	lons := denseNodes.GetLon()
	denseinfo := denseNodes.GetDenseinfo()

	tu := tagUnpacker{
		stringTable: stringTable,
		keysVals:    denseNodes.GetKeysVals(),
		index:       0,
	}
	var id, lat, lon int64
	var state denseInfoState
	for index := range ids {
		id = ids[index] + id
		lat = lats[index] + lat
		lon = lons[index] + lon
		latitude := degMultiplier * float64(latOffset+(granularity*lat))
		longitude := degMultiplier * float64(lonOffset+(granularity*lon))
		tags := tu.next()
		info := extractDenseInfo(stringTable, &state, denseinfo, index, dateGranularity)

		i.data = append(i.data, osmpbfreaderdata.Node{
			ID:   id,
			Lat:  latitude,
			Lon:  longitude,
			Tags: tags,
			Info: info,
		})
	}
}

func (i *impl) parseWays(block *osmproto.PrimitiveBlock, ways []*osmproto.Way) {
	if i.filter != nil && i.filter.FilterWays() {
		return
	}

	stringTable := block.GetStringtable().GetS()
	dateGranularity := int64(block.GetDateGranularity())

	for _, way := range ways {
		id := way.GetId()

		tags := extractTags(stringTable, way.GetKeys(), way.GetVals())

		refs := way.GetRefs()
		var nodeID int64
		nodeIDs := make([]int64, len(refs))
		for index := range refs {
			nodeID = refs[index] + nodeID // delta encoding
			nodeIDs[index] = nodeID
		}

		info := extractInfo(stringTable, way.GetInfo(), dateGranularity)

		i.data = append(i.data, osmpbfreaderdata.Way{
			ID:      id,
			Tags:    tags,
			NodeIDs: nodeIDs,
			Info:    info,
		})
	}
}

func (i *impl) parseRelations(block *osmproto.PrimitiveBlock, relations []*osmproto.Relation) {
	if i.filter != nil && i.filter.FilterRelations() {
		return
	}

	stringTable := block.GetStringtable().GetS()
	dateGranularity := int64(block.GetDateGranularity())

	for _, rel := range relations {
		id := rel.GetId()
		tags := extractTags(stringTable, rel.GetKeys(), rel.GetVals())
		members := extractMembers(stringTable, rel)
		info := extractInfo(stringTable, rel.GetInfo(), dateGranularity)

		i.data = append(i.data, osmpbfreaderdata.Relation{
			ID:      id,
			Tags:    tags,
			Members: members,
			Info:    info,
		})
	}
}

// Make relation members from stringtable and three parallel arrays of IDs.
func extractMembers(stringTable [][]byte, relation *osmproto.Relation) []osmpbfreaderdata.Member {
	memIDs := relation.GetMemids()
	types := relation.GetTypes()
	roleIDs := relation.GetRolesSid()

	var memID int64
	members := make([]osmpbfreaderdata.Member, len(memIDs))
	for index := range memIDs {
		memID = memIDs[index] + memID // delta encoding

		var memType osmpbfreaderdata.MemberType
		switch types[index] {
		case osmproto.Relation_NODE:
			memType = osmpbfreaderdata.NodeType
		case osmproto.Relation_WAY:
			memType = osmpbfreaderdata.WayType
		case osmproto.Relation_RELATION:
			memType = osmpbfreaderdata.RelationType
		}

		role := stringTable[roleIDs[index]]

		members[index] = osmpbfreaderdata.Member{
			ID:   memID,
			Type: memType,
			Role: string(role),
		}
	}

	return members
}

func extractInfo(stringTable [][]byte, protoInfo *osmproto.Info, dateGranularity int64) osmpbfreaderdata.Info {
	info := osmpbfreaderdata.Info{Visible: true}

	if protoInfo != nil {
		info.Version = protoInfo.GetVersion()

		millisec := time.Duration(protoInfo.GetTimestamp()*dateGranularity) * time.Millisecond
		info.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()

		info.Changeset = protoInfo.GetChangeset()

		info.Uid = protoInfo.GetUid()

		info.User = string(stringTable[protoInfo.GetUserSid()])

		if protoInfo.Visible != nil {
			info.Visible = protoInfo.GetVisible()
		}
	}

	return info
}

type denseInfoState struct {
	timestamp int64
	changeset int64
	uid       int32
	userSid   int32
}

func extractDenseInfo(stringTable [][]byte, state *denseInfoState, denseInfo *osmproto.DenseInfo, index int, dateGranularity int64) osmpbfreaderdata.Info {
	info := osmpbfreaderdata.Info{Visible: true}

	versions := denseInfo.GetVersion()
	if len(versions) > 0 {
		info.Version = versions[index]
	}

	timestamps := denseInfo.GetTimestamp()
	if len(timestamps) > 0 {
		state.timestamp = timestamps[index] + state.timestamp
		millisec := time.Duration(state.timestamp*dateGranularity) * time.Millisecond
		info.Timestamp = time.Unix(0, millisec.Nanoseconds()).UTC()
	}

	changesets := denseInfo.GetChangeset()
	if len(changesets) > 0 {
		state.changeset = changesets[index] + state.changeset
		info.Changeset = state.changeset
	}

	uids := denseInfo.GetUid()
	if len(uids) > 0 {
		state.uid = uids[index] + state.uid
		info.Uid = state.uid
	}

	usersids := denseInfo.GetUserSid()
	if len(usersids) > 0 {
		state.userSid = usersids[index] + state.userSid
		info.User = string(stringTable[state.userSid])
	}

	visibleArray := denseInfo.GetVisible()
	if len(visibleArray) > 0 {
		info.Visible = visibleArray[index]
	}

	return info
}
