package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server is an MCP server.
type Server struct {
	server *mcp.Server
}

// NewServer creates a new MCP server.
func NewServer(name, version string) *Server {
	return &Server{mcp.NewServer(&mcp.Implementation{Name: name, Version: version}, nil)}
}

// AddTool adds a new tool to the MCP server.
func AddTool[I, O any](server *Server, name, description string, f func(context.Context, *I) (*O, error)) {
	mcp.AddTool(server.server, &mcp.Tool{Name: name, Description: description},
		func(ctx context.Context, _ *mcp.CallToolRequest, in *I) (*mcp.CallToolResult, *O, error) {
			out, err := f(ctx, in)
			return nil, out, err
		})
}

// Run runs the MCP server.
func (server *Server) Run(ctx context.Context) error {
	return server.server.Run(ctx, new(mcp.StdioTransport))
}
