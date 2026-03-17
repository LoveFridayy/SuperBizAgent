package plan_execute_replan

import (
	"SuperBizAgent/internal/ai/models"
	"SuperBizAgent/internal/ai/tools"
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/adk/prebuilt/planexecute"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewExecutor(ctx context.Context) (adk.Agent, error) {
	// log
	mcpTool, err := tools.GetLogMcpTool()
	if err != nil {
		return nil, err
	}
	toolList := make([]tool.BaseTool, 0)
	// 只有当 MCP 工具可用时才添加
	if len(mcpTool) > 0 {
		toolList = append(toolList, mcpTool...)
	}
	// alerts
	toolList = append(toolList, tools.NewPrometheusAlertsQueryTool())
	// file
	toolList = append(toolList, tools.NewQueryInternalDocsTool())
	// time
	toolList = append(toolList, tools.NewGetCurrentTimeTool())
	execModel, err := models.OpenAIForDeepSeekV3Quick(ctx)
	if err != nil {
		return nil, err
	}
	return planexecute.NewExecutor(ctx, &planexecute.ExecutorConfig{
		Model: execModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: toolList,
			},
		},
		MaxIterations: 999999,
	})
}
