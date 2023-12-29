package probe

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// Ready 检查进程是否完成所有的初始化
func Ready(ctx context.Context, in map[string]interface{}) error {
	log.Context(ctx).Debug("receive Readiness probe")

	// TODO 进行必要的检查

	return nil
}
