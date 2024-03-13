package write

import (
	"encoding/json"
	"log"
	h "stage2024/pkg/helper"
	"stage2024/pkg/protogen/common"
	"stage2024/pkg/protogen/occupations"

	"log/slog"

	ravendb "github.com/ravendb/ravendb-go-client"
)

var globalDocumentStore *ravendb.DocumentStore

func createDocumentStore() (*ravendb.DocumentStore, error) {

	if globalDocumentStore != nil {
		return globalDocumentStore, nil
	}

	serverURL := "http://localhost:8080"
	databaseName := "Database"

	urls := []string{serverURL}
	store := ravendb.NewDocumentStore(urls, databaseName)

	//store.Certificate = ...
	//store.TrustStore = ...

	conventions := ravendb.NewDocumentConventions()
	// Modify conventions as needed
	store.SetConventions(conventions)

	err := store.Initialize()
	if err != nil {
		return nil, err
	}

	globalDocumentStore = store

	if globalDocumentStore == nil {
		panic("globalDocumentStore is nil")
	}

	return globalDocumentStore, nil
}

func workWithSession(record string) error {
	session, err := globalDocumentStore.OpenSession("Database")
	h.MaybeDie(err, "Error opening session")

	//   Run your business logic:
	//

	slog.Info("Session opened")
	// record := `{"stationId":"23898","num_bikes_available": 1,"numDocksAvailable":3,"isRenting":1,"isInstalled":1,"isReturning":1,"lastReported":"1708365148","location":{"lon":3.7467302,"lat":51.0623281},"name":"Artevelde Hogeschool Sint Amandsberg"}`

	err = createRecords(record)
	h.MaybeDieErr(err)

	err = session.SaveChanges()
	h.MaybeDie(err, "Error saving changes")

	session.Close()

	return nil
}

// main
func WriteToRaven(record string) error {
	log.Println("Starting RavenDB client")
	_, err := createDocumentStore()
	h.MaybeDieErr(err)
	slog.Info("created document store")

	err = workWithSession(record)
	h.MaybeDieErr(err)

	return nil
}

func createDocument(out *occupations.DonkeyLocation) error {

	slog.Info("Creating document", "documentID", out.StationId)

	session, err := globalDocumentStore.OpenSession("Database")
	h.MaybeDie(err, "Error opening session")

	err = session.Store(out)
	h.MaybeDie(err, "Error storing document")

	err = session.SaveChanges()
	h.MaybeDie(err, "error saving changes")

	return nil
}

func createRecords(record string) error {
	var in DonkeyLocation
	err := json.Unmarshal([]byte(record), &in)
	if err != nil {
		slog.Warn("Failed to unmarshal", "error", err)
	}
	out := occupations.DonkeyLocation{}
	out.StationId = in.StationId
	out.NumBikesAvailable = in.NumBikesAvailable
	out.NumDocksAvailable = in.NumDocksAvailable
	out.IsRenting = in.IsRenting
	out.IsInstalled = in.IsInstalled
	out.IsReturning = in.IsReturning
	out.LastReported = in.LastReported
	out.Location = &common.Location{
		Lon: in.Location.Lon,
		Lat: in.Location.Lat,
	}
	out.Name = in.Name

	slog.Info("record", "StationID", out.StationId, "Num_bikes_available", out.NumBikesAvailable, "Num_docks_available", out.NumDocksAvailable, "Is_renting", out.IsRenting, "Is_installed", out.IsInstalled, "Is_returning", out.IsReturning, "Last_reported", out.LastReported, "Location", out.Location, "Name", out.Name)

	err = createDocument(&out)
	h.MaybeDieErr(err)
	return nil
}

type DonkeyLocation struct {
	StationId         string `json:"stationId"`
	NumBikesAvailable int32  `json:"numBikesAvailable"`
	NumDocksAvailable int32  `json:"numDocksAvailable"`
	IsRenting         int32  `json:"isRenting"`
	IsInstalled       int32  `json:"isInstalled"`
	IsReturning       int32  `json:"isReturning"`
	LastReported      string `json:"lastReported"`
	Location          struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	}
	Name string `json:"name"`
}

func Close() {
	globalDocumentStore.Close()
}
