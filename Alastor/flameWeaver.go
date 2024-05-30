package alastor

import (
	"context"
	"errors"
	"github.com/spf13/afero"
	"log"
	"sync"
	"time"
)

type FlameWeaver struct {
	AuthCode    string
	AddressPort string
	FileName    string
	ChunkSize   int
	Destination *afero.File
	dqs         []ChunkTracker
	nextDq      int
	servants    []FlameServant
	awaiter     sync.WaitGroup
	mu          sync.Mutex
}

func (f *FlameWeaver) notifyServantActive() {
	f.awaiter.Add(1)
}

func (f *FlameWeaver) getNextDQ() (*ChunkTracker, error) {
	f.mu.Lock()
	if f.nextDq >= len(f.dqs) {
		f.mu.Unlock()
		return nil, errors.New("no DataQuery found")
	} else {
		next := f.dqs[f.nextDq]
		f.nextDq++
		f.mu.Unlock()
		return &next, nil
	}
}

func (f *FlameWeaver) notifyServantDead() {
	f.awaiter.Done()
}

func (f *FlameWeaver) Weave() error {

	q := &FileQuery{
		ApiKey:            f.AuthCode,
		RequestedFilename: f.FileName,
	}

	con, err := newConnection(f.AddressPort)
	if err != nil {
		return err
	}
	defer con.Close()
	c := NewAlastorClient(con)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	info, err := c.GetFileInfo(ctx, q)
	if err != nil {
		return err
	}
	if info.Error.Code != 0 {
		return errors.New(info.Error.Msg)
	}
	if info.IsDirectory || len(info.Info) != 1 {
		return errors.New("is directory")
	}

	chunks := info.Info[0].FileSize / int64(f.ChunkSize)
	remaining := info.Info[0].FileSize % int64(f.ChunkSize)
	if remaining > 0 {
		chunks++
	}

	for i := range chunks {
		f.dqs = append(f.dqs, ChunkTracker{
			dq: DataQuery{
				ApiKey:            f.AuthCode,
				RequestedFilename: f.FileName,
				ChunkSize:         int32(f.ChunkSize),
				ChunkOffset:       i,
			},
			dest: f.Destination,
		})
	}

	for range 20 {
		newServant := &FlameServant{c}
		go newServant.Run(f)
		f.servants = append(f.servants, *newServant)
	}

	f.awaiter.Wait()
	log.Println("done")
	return nil

}
