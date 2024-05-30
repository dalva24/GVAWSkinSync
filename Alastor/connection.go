package alastor

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net.dalva.GvawSkinSync/ssErrors"
	"time"
)

func newConnection(addressPort string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addressPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func TestConnection(addressPort string, authCode string) *ssErrors.UiError {
	con, err := newConnection(addressPort)
	if err != nil {
		return ssErrors.NewUiError("Connection Error", err.Error())
	}
	defer con.Close()
	c := NewAlastorClient(con)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	t, err := c.Command(ctx, &CommandQuery{
		ApiKey:  authCode,
		Command: "test",
	})
	if err != nil {
		if status.Code(err) == codes.Unauthenticated {
			return ssErrors.NewUiError("Auth Error", "Make sure to use latest Auth Key to sync. Check Discord.")
		} else {
			return ssErrors.NewUiError("Error", err.Error())
		}
	}
	if t.GetError().Code != 0 {
		return ssErrors.NewUiError("Error", t.GetError().Msg)
	}
	return nil
}
