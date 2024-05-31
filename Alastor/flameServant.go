package alastor

import (
	"context"
	"errors"
	"github.com/spf13/afero"
	"net.dalva.GvawSkinSync/checksum"
	"net.dalva.GvawSkinSync/logger"
	"net.dalva.GvawSkinSync/ssErrors"
	"time"
)

type ChunkTracker struct {
	dq   DataQuery
	dest *afero.File
}

type FlameServant struct {
	addressPort string
	c           AlastorClient
}

func (f *FlameServant) Run(supervisor *FlameWeaver) {
	supervisor.notifyServantActive()
	go f.Runtime(supervisor)
}

func (f *FlameServant) Runtime(supervisor *FlameWeaver) {

	con, err := newConnection(f.addressPort)
	if err != nil {
		supervisor.notifyServantDead()
		logger.Log.Error().Err(err).Msg("Servant creation error")
		return
	}
	defer con.Close()
	f.c = NewAlastorClient(con)

	for {
		next, err := supervisor.getNextDQ()
		if err != nil {
			break
		}

		var chunk []byte

		//download till success
		for {
			var err error
			chunk, err = f.downloadChunk(&next.dq)
			if err != nil {
				logger.Log.Error().Err(err).Int64("chunk", next.dq.ChunkOffset).Msg("Error downloading chunk")
			} else {
				//done downloading
				break
			}
		}

		offset := next.dq.ChunkOffset * int64(next.dq.ChunkSize)

		//write till success
		for {
			var err error
			err = f.saveChunk(chunk, next.dest, offset)
			if err != nil {
				logger.Log.Error().Err(err).Int64("chunk", next.dq.ChunkOffset).Msg("Error saving chunk")
			} else {
				//done downloading
				break
			}
		}
	}

	supervisor.notifyServantDead()
}

func (f *FlameServant) downloadChunk(dq *DataQuery) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	data, err := f.c.GetFileData(ctx, dq)
	if err != nil {
		return nil, err
	}
	if data.Error.Code != 0 {
		return nil, errors.New(data.Error.Msg)
	}
	if checksum.CRC(data.ChunkData) != data.ChunkCrc32 {
		return nil, ssErrors.NewDataCrcError("chunk crc error")
	}

	return data.ChunkData, nil
}

func (f *FlameServant) saveChunk(bytes []byte, file *afero.File, offset int64) error {

	q := *file
	n, err := q.WriteAt(bytes, offset)
	if err != nil {
		return err
	}
	if n != len(bytes) {
		return errors.New("out of space maybe")
	}

	return nil
}
