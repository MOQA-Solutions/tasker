package registry 

import "github.com/MOQA-Solutions/tasker/types"

var (
	Workers map[int][]*types.ProtectedChannel
    Pool *types.Pool
    CancelContext *types.CancelContext
   )

func init() {
	Pool = types.NewPool()
	CancelContext = types.NewCancelContext()
}
