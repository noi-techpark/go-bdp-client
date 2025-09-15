// SPDX-FileCopyrightText: 2024 NOI Techpark <digital@noi.bz.it>
//
// SPDX-License-Identifier: MPL-2.0

package bdpmock

import (
	"github.com/noi-techpark/go-bdp-client/bdplib"
)

type BdpMockCalls struct {
	SyncedDataTypes [][]bdplib.DataType             `json:"syncedDataTypes"`
	SyncedData      map[string][]bdplib.DataMap     `json:"syncedData"`
	SyncedStations  map[string][]BdpMockStationCall `json:"syncedStations"`
}

type BdpMockStationCall struct {
	Stations     []bdplib.Station
	SyncState    bool
	OnlyActivate bool
}

type BdpMock struct {
	ProvenanceUuid string
	Prv            string
	Prn            string
	Origin         string

	// stationType - []DataType
	SyncedDataTypes [][]bdplib.DataType
	// stationType - []DataMap
	SyncedData map[string][]bdplib.DataMap
	// stationType - []BdpMockStationCall
	SyncedStations map[string][]BdpMockStationCall
}

func MockFromEnv(e bdplib.BdpEnv) bdplib.Bdp {
	b := BdpMock{}
	b.Prv = e.BDP_PROVENANCE_VERSION
	b.Prn = e.BDP_PROVENANCE_NAME
	b.Origin = e.BDP_ORIGIN
	b.SyncedData = make(map[string][]bdplib.DataMap)
	b.SyncedStations = make(map[string][]BdpMockStationCall)
	b.SyncedDataTypes = [][]bdplib.DataType{}
	return &b
}

func (b *BdpMock) CreateDataMap() bdplib.DataMap {
	var dataMap = bdplib.DataMap{
		Name:       "(default)",
		Provenance: b.ProvenanceUuid,
		Branch:     make(map[string]bdplib.DataMap),
	}
	return dataMap
}

func (b *BdpMock) SyncDataTypes(dataTypes []bdplib.DataType) error {
	b.SyncedDataTypes = append(b.SyncedDataTypes, dataTypes)
	return nil
}

func (b *BdpMock) SyncStations(stationType string, stations []bdplib.Station, syncState bool, onlyActivate bool) error {
	call := BdpMockStationCall{
		Stations:     stations,
		SyncState:    syncState,
		OnlyActivate: onlyActivate,
	}
	if _, ok := b.SyncedStations[stationType]; ok {
		b.SyncedStations[stationType] = append(b.SyncedStations[stationType], call)
	} else {
		b.SyncedStations[stationType] = []BdpMockStationCall{call}
	}
	return nil
}

func (b *BdpMock) PushData(stationType string, dataMap bdplib.DataMap) error {
	if _, ok := b.SyncedData[stationType]; ok {
		b.SyncedData[stationType] = append(b.SyncedData[stationType], dataMap)
	} else {
		b.SyncedData[stationType] = []bdplib.DataMap{dataMap}
	}
	return nil
}

func (b *BdpMock) GetOrigin() string {
	return b.Origin
}

func (b *BdpMock) Requests() BdpMockCalls {
	return BdpMockCalls{
		SyncedDataTypes: b.SyncedDataTypes,
		SyncedData:      b.SyncedData,
		SyncedStations:  b.SyncedStations,
	}
}
