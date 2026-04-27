package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// UnlockCircuit 定义了解锁电路的约束系统
// 它验证：H == MiMC(Pre, Sn)
type UnlockCircuit struct {
	// 公开输入 (public) - 暴露在链上
	H  frontend.Variable `gnark:",public"` // H 也是 public
	Sn frontend.Variable `gnark:",public"` // [!] 关键修正：Sn 必须是 public

	// 私有输入 (witness) - 保持秘密
	Pre frontend.Variable `gnark:",secret"` // 这是唯一需要保密的

}

// Define 定义电路的约束
func (circuit *UnlockCircuit) Define(api frontend.API) error {
	// 使用 MiMC 哈希函数
	mimcHash, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	// 计算 hash = MiMC(Pre, Sn)
	// MiMC 接收私有输入 Pre 和公开输入 Sn
	mimcHash.Write(circuit.Pre)
	mimcHash.Write(circuit.Sn)
	computedHash := mimcHash.Sum()

	// 约束：计算出的哈希必须等于公开的 H
	api.AssertIsEqual(computedHash, circuit.H)

	return nil
}
