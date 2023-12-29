// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             v3.21.12
// source: example/sample.proto

package example

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationSampleCheckSignupStatus = "/example.sample/CheckSignupStatus"

type SampleHTTPServer interface {
	// CheckSignupStatus Sends a greeting
	CheckSignupStatus(context.Context, *CheckSignupStatusRequest) (*CheckSignupStatusResponse, error)
}

func RegisterSampleHTTPServer(s *http.Server, srv SampleHTTPServer) {
	r := s.Route("/")
	r.GET("/sample/signup/status", _Sample_CheckSignupStatus0_HTTP_Handler(srv))
}

func _Sample_CheckSignupStatus0_HTTP_Handler(srv SampleHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CheckSignupStatusRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationSampleCheckSignupStatus)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CheckSignupStatus(ctx, req.(*CheckSignupStatusRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CheckSignupStatusResponse)
		return ctx.Result(200, reply)
	}
}

type SampleHTTPClient interface {
	CheckSignupStatus(ctx context.Context, req *CheckSignupStatusRequest, opts ...http.CallOption) (rsp *CheckSignupStatusResponse, err error)
}

type SampleHTTPClientImpl struct {
	cc *http.Client
}

func NewSampleHTTPClient(client *http.Client) SampleHTTPClient {
	return &SampleHTTPClientImpl{client}
}

func (c *SampleHTTPClientImpl) CheckSignupStatus(ctx context.Context, in *CheckSignupStatusRequest, opts ...http.CallOption) (*CheckSignupStatusResponse, error) {
	var out CheckSignupStatusResponse
	pattern := "/sample/signup/status"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationSampleCheckSignupStatus))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}