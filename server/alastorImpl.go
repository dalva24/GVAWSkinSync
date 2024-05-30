package main

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=.. Alastor/alastor.proto

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
	alastor "net.dalva.GvawSkinSync/Alastor"
	"net.dalva.GvawSkinSync/checksum"
	"net.dalva.GvawSkinSync/conf"
	"net.dalva.GvawSkinSync/data"
	"net.dalva.GvawSkinSync/logger"
	"strings"
)

func serveAlastor() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Cfg.Port))
	if err != nil {
		logger.Log.Err(err).Msg("Error Listening")
	}
	s := grpc.NewServer()
	alastor.RegisterAlastorServer(s, &server{})
	logger.Log.Info().Msg(fmt.Sprintf("server listening at %v", lis.Addr()))
	if err := s.Serve(lis); err != nil {
		logger.Log.Err(err).Msg("Error Serving")
	}
}

type server struct {
	alastor.UnimplementedAlastorServer
}

var noError = &alastor.ErrorMsg{
	Code: 0,
	Msg:  "",
}

var generalIoError = alastor.ErrorMsg{
	Code: 10,
	Msg:  "",
}

func fileInfoError(code int, err error) *alastor.FileInfo {
	return &alastor.FileInfo{
		Error: &alastor.ErrorMsg{
			Code: int32(code),
			Msg:  err.Error(),
		},
		IsDirectory: false,
		Info:        nil,
	}
}

func fileDataError(code int, err error) *alastor.FileData {
	return &alastor.FileData{
		Error: &alastor.ErrorMsg{
			Code: int32(code),
			Msg:  err.Error(),
		},
		ChunkData:  nil,
		ChunkCrc32: 0,
	}
}

func (s *server) GetFileInfo(ctx context.Context, fq *alastor.FileQuery) (*alastor.FileInfo, error) {
	if fq.ApiKey == "" || !strings.EqualFold(fq.ApiKey, authCode) {
		return nil, status.Error(codes.Unauthenticated, "Invalid Auth Key")
	}

	f, err := fs.Stat(fq.RequestedFilename)
	if err != nil {
		return fileInfoError(10, err), nil
	}

	switch mode := f.Mode(); {

	//wanted file is a directory
	case mode.IsDir():

		dir, err := afero.ReadDir(fs, fq.RequestedFilename)
		if err != nil {
			return fileInfoError(10, err), nil
		}

		var files []*alastor.File

		for _, info := range dir {
			files = append(files, &alastor.File{
				FileName:      info.Name(),
				FileSize:      info.Size(),
				FileTimestamp: info.ModTime().UnixMilli(),
			})
		}

		return &alastor.FileInfo{
			Error:       noError,
			IsDirectory: true,
			Info:        files,
		}, nil

	//wanted file is a file
	case mode.IsRegular():
		file := make([]*alastor.File, 1)
		file[0] = &alastor.File{
			FileName:      f.Name(),
			FileSize:      f.Size(),
			FileTimestamp: f.ModTime().UnixMilli(),
		}
		return &alastor.FileInfo{
			Error:       noError,
			IsDirectory: false,
			Info:        file,
		}, nil

	//wanted file is... whatever it is...
	default:
		return &alastor.FileInfo{
			Error: &alastor.ErrorMsg{
				Code: 10,
				Msg:  "unknown or not exist",
			},
			IsDirectory: false,
			Info:        nil,
		}, nil

	}
}

func (s *server) GetFileData(ctx context.Context, dq *alastor.DataQuery) (*alastor.FileData, error) {
	if dq.ApiKey == "" || !strings.EqualFold(dq.ApiKey, authCode) {
		return nil, status.Error(codes.Unauthenticated, "Invalid Auth Key")
	}

	f, err := fs.Open(dq.RequestedFilename)
	if err != nil {
		return fileDataError(10, err), nil
	}

	chunk, err := data.ReadOffsetChunk(f, dq.ChunkOffset, dq.ChunkSize)
	if err != nil {
		return fileDataError(12, err), nil
	}

	return &alastor.FileData{
		Error:      noError,
		ChunkData:  chunk,
		ChunkCrc32: checksum.CRC(chunk),
	}, nil
}

func (s *server) Command(ctx context.Context, cq *alastor.CommandQuery) (*alastor.CommandReply, error) {
	p, _ := peer.FromContext(ctx)
	logger.Log.Info().
		Str("apiKey", cq.ApiKey).
		Str("IP", p.Addr.String()).
		Msg("New Connection")
	if cq.ApiKey == "" || !strings.EqualFold(cq.ApiKey, authCode) {
		return nil, status.Error(codes.Unauthenticated, "Invalid Auth Key")
	}

	if cq.Command == "aircraft" {
		aircraft := data.LoadLines("aircraft.txt", fs)
		return &alastor.CommandReply{
			Error:        noError,
			ReturnString: aircraft,
			ReturnInt32:  nil,
		}, nil
	}

	return &alastor.CommandReply{
		Error:        noError,
		ReturnString: nil,
		ReturnInt32:  nil,
	}, nil
}
