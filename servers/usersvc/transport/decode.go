package transport

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/kum0/go-mircosvc/common"
	userPb "github.com/kum0/go-mircosvc/pb/user"
	"github.com/kum0/go-mircosvc/servers/usersvc/endpoints"
	"net/http"
	"strconv"
)

// GerUser
func decodeGetUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	uid, ok := vars["UID"]
	if !ok {
		return nil, common.ErrRouteArgs
	}
	return uid, nil
}

func decodeGRPCGetUserRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	r, ok := grpcReq.(*userPb.GetUserRequest)
	if !ok {
		return nil, errors.New("decodeGRPCGetUserRequest: interface conversion error")
	}
	return r.Uid, nil
}

func decodeGRPCGetUserResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	rp, ok := grpcResponse.(*userPb.GetUserResponse)
	if !ok {
		return nil, errors.New("decodeGRPCGetUserResponse: interface conversion error")
	}

	// r := new(userPb.GetUserResponse)
	// if err := utils.StructCopy(rp, r); err != nil {
	// 	return nil, err
	// }
	return rp, nil
}

// Login
func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*userPb.LoginRequest)
	if !ok {
		return nil, errors.New("decodeGRPCLoginRequest: interface conversion error")
	}
	return &endpoints.LoginRequest{Username: req.Username, Password: req.Password}, nil
}

func decodeGRPCLoginResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	rp, ok := grpcResponse.(*userPb.LoginResponse)
	if !ok {
		return nil, errors.New("decodeGRPCLoginResponse: interface conversion error")
	}

	return rp, nil
}

// SendCode
func decodeGRPCSendCodeResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	rp, ok := grpcResponse.(*userPb.SendCodeResponse)
	if !ok {
		return nil, errors.New("decodeGRPCSendCodeResponse: interface conversion error")
	}

	// r := new(userPb.SendCodeResponse)
	// if err := utils.StructCopy(rp, r); err != nil {
	// 	return nil, err
	// }
	return rp, nil
}

// Register
func decodeGRPCRegisterRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*userPb.RegisterRequest)
	if !ok {
		return nil, errors.New("decodeGRPCRegisterRequest: interface conversion error")
	}
	return &endpoints.RegisterRequest{Username: req.Username, Password: req.Password, CodeID: req.CodeID}, nil
}

// UserList
func DecodeUserListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	q := r.URL.Query()
	page, err := strconv.ParseInt(q.Get("page"), 10, 0)
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(q.Get("size"), 10, 0)
	if err != nil {
		return nil, err
	}
	return &endpoints.UserListRequest{Page: int32(page), Size: int32(size)}, nil
}

func decodeGRPCUserListResponse(_ context.Context, grpcResponse interface{}) (interface{}, error) {
	rp, ok := grpcResponse.(*userPb.UserListResponse)
	if !ok {
		return nil, errors.New("decodeGRPCUserListResponse: interface conversion error")
	}

	return rp, nil
}

func decodeGRPCUserListRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req, ok := grpcReq.(*userPb.UserListRequest)
	if !ok {
		return nil, errors.New("decodeGRPCUserListRequest: interface conversion error")
	}
	return &endpoints.UserListRequest{Page: req.Page, Size: req.Size}, nil
}

// Logout
func decodeLogoutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	sid, ok := ctx.Value(common.SessionKey).(string)
	if !ok {
		return nil, common.NewError(http.StatusUnauthorized, "cookie 不存在")
	}

	return &endpoints.LogoutRequest{sid}, nil
}

func decodeGRPCLogoutRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	rp, ok := grpcReq.(*userPb.LogoutRequest)
	if !ok {
		return nil, errors.New("decodeGRPCUserListResponse: interface conversion error")
	}

	return &endpoints.LogoutRequest{rp.Sid}, nil
}
