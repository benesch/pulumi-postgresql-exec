// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/url"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"github.com/jackc/pgconn"
	"github.com/pulumi/pulumi/pkg/v3/resource/provider"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource/plugin"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/cmdutil"
	rpc "github.com/pulumi/pulumi/sdk/v3/proto/go"
)

// Injected by linker in release builds.
var version string

func main() {
	err := provider.Main("postgresql-exec", func(host *provider.HostClient) (rpc.ResourceProviderServer, error) {
		return &postgresqlExecProvider{
			host: host,
		}, nil
	})
	if err != nil {
		cmdutil.ExitError(err.Error())
	}
}

type postgresqlExecProvider struct {
	host *provider.HostClient
	conn *pgconn.PgConn
}

func (p *postgresqlExecProvider) CheckConfig(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	return &rpc.CheckResponse{Inputs: req.GetNews()}, nil
}

func (p *postgresqlExecProvider) DiffConfig(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	return &rpc.DiffResponse{}, nil
}

func (p *postgresqlExecProvider) Configure(ctx context.Context, req *rpc.ConfigureRequest) (*rpc.ConfigureResponse, error) {
	vars := req.GetVariables()
	host := vars["postgresql-exec:config:host"]
	port := vars["postgresql-exec:config:port"]
	database := vars["postgresql-exec:config:database"]
	user := vars["postgresql-exec:config:user"]
	password := vars["postgresql-exec:config:password"]

	connStr := "postgresql://"
	if len(user) > 0 {
		connStr += url.PathEscape(user)
		if len(password) > 0 {
			connStr += ":"
			connStr += url.PathEscape(password)
		}
		connStr += "@"
	}
	if len(host) > 0 {
		connStr += url.PathEscape(host)
	}
	if len(port) > 0 {
		connStr += ":"
		connStr += url.PathEscape(port)
	}
	if len(database) > 0 {
		connStr += "/"
		connStr += url.PathEscape(database)
	}

	conn, err := pgconn.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	p.conn = conn
	return &rpc.ConfigureResponse{}, nil
}

func (p *postgresqlExecProvider) Invoke(_ context.Context, req *rpc.InvokeRequest) (*rpc.InvokeResponse, error) {
	tok := req.GetTok()
	return nil, fmt.Errorf("Unknown Invoke token '%s'", tok)
}

func (p *postgresqlExecProvider) StreamInvoke(req *rpc.InvokeRequest, server rpc.ResourceProvider_StreamInvokeServer) error {
	tok := req.GetTok()
	return fmt.Errorf("Unknown StreamInvoke token '%s'", tok)
}

func (p *postgresqlExecProvider) Check(ctx context.Context, req *rpc.CheckRequest) (*rpc.CheckResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "postgresql-exec:index:Exec" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	return &rpc.CheckResponse{Inputs: req.News, Failures: nil}, nil
}

func (p *postgresqlExecProvider) Diff(ctx context.Context, req *rpc.DiffRequest) (*rpc.DiffResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "postgresql-exec:index:Exec" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	olds, err := plugin.UnmarshalProperties(req.GetOlds(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	news, err := plugin.UnmarshalProperties(req.GetNews(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	d := olds.Diff(news)
	if d == nil {
		return &rpc.DiffResponse{
			Changes: rpc.DiffResponse_DIFF_NONE,
		}, nil
	}
	replaces := []string{}
	if d.Changed("createSql") {
		replaces = append(replaces, "createSql")
	}
	if d.Changed("destroySql") {
		replaces = append(replaces, "destroySql")
	}
	return &rpc.DiffResponse{
		Changes:             rpc.DiffResponse_DIFF_SOME,
		DeleteBeforeReplace: true,
		Replaces:            replaces,
	}, nil
}

func (p *postgresqlExecProvider) Create(ctx context.Context, req *rpc.CreateRequest) (*rpc.CreateResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "postgresql-exec:index:Exec" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	_, err = p.conn.Exec(ctx, inputs["createSql"].StringValue()).ReadAll()
	if err != nil {
		return nil, err
	}
	outputs := map[string]interface{}{
		"createSql":  inputs["createSql"].StringValue(),
		"destroySql": inputs["destroySql"].StringValue(),
	}
	outputProperties, err := plugin.MarshalProperties(
		resource.NewPropertyMapFromMap(outputs),
		plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true},
	)
	if err != nil {
		return nil, err
	}
	return &rpc.CreateResponse{
		Id:         inputs["createSql"].StringValue(),
		Properties: outputProperties,
	}, nil
}

func (p *postgresqlExecProvider) Read(ctx context.Context, req *rpc.ReadRequest) (*rpc.ReadResponse, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "postgresql-exec:index:Exec" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	return &rpc.ReadResponse{
		Id:         req.Id,
		Properties: req.Properties,
	}, nil
}

func (p *postgresqlExecProvider) Update(ctx context.Context, req *rpc.UpdateRequest) (*rpc.UpdateResponse, error) {
	panic("Update not implemented")
}

func (p *postgresqlExecProvider) Delete(ctx context.Context, req *rpc.DeleteRequest) (*pbempty.Empty, error) {
	urn := resource.URN(req.GetUrn())
	ty := urn.Type()
	if ty != "postgresql-exec:index:Exec" {
		return nil, fmt.Errorf("Unknown resource type '%s'", ty)
	}
	inputs, err := plugin.UnmarshalProperties(req.GetProperties(), plugin.MarshalOptions{KeepUnknowns: true, SkipNulls: true})
	if err != nil {
		return nil, err
	}
	_, err = p.conn.Exec(ctx, inputs["destroySql"].StringValue()).ReadAll()
	if err != nil {
		return nil, err
	}
	return &pbempty.Empty{}, nil
}

func (p *postgresqlExecProvider) Construct(_ context.Context, _ *rpc.ConstructRequest) (*rpc.ConstructResponse, error) {
	panic("Construct not implemented")
}

func (p *postgresqlExecProvider) GetPluginInfo(context.Context, *pbempty.Empty) (*rpc.PluginInfo, error) {
	return &rpc.PluginInfo{
		Version: version,
	}, nil
}

func (p *postgresqlExecProvider) GetSchema(ctx context.Context, req *rpc.GetSchemaRequest) (*rpc.GetSchemaResponse, error) {
	return &rpc.GetSchemaResponse{}, nil
}

func (p *postgresqlExecProvider) Cancel(context.Context, *pbempty.Empty) (*pbempty.Empty, error) {
	return &pbempty.Empty{}, nil
}
