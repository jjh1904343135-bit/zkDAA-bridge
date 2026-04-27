const hre = require("hardhat");

async function main() {
    // 直接使用验证器地址
    const verifierAddr = "0x5FbDB2315678afecb367f032d93F642f64180aa3";
    
    console.log("直接测试验证器:", verifierAddr);
    
    const verifier = await hre.ethers.getContractAt("BatchUnlockVerifier", verifierAddr);
    
    // 使用你最新的测试数据
    const proof = [
        "13368861877719866834911251902871094413280731175582118433810691026950071972923",
        "1184973859348026248675423533385659407784818956520722808008997950775670981580",
        "13033890459992367420088521473871045132573803918963449201344554610987096553980",
        "8340599338524072968624523938182243217220151816189572844862798223800854113233",
        "378066865481376785013725786882190914543582922421285639235070387430030450085",
        "9859682356058182196341441473153588372052614326143035811564769646357472268840",
        "10195766642147751489383170932571745008403358088453635420417994568824592086443",
        "20948580297244079446910708042905443333219070905826842562042114736404460729668"
    ];
    
    const publicInputs = [
        "7649593193017438476480190740478439093036199922286055488656277972216650429398",
        "302111519086211411682947601071787884841"
    ];
    
    console.log("\n测试数据:");
    console.log("  Proof 长度:", proof.length);
    console.log("  PublicInputs:", publicInputs);
    
    try {
        console.log("\n调用 verifyProof...");
        await verifier.verifyProof.staticCall(proof, publicInputs);
        console.log("✅ 验证器验证通过！");
    } catch (error) {
        console.log("❌ 验证器验证失败!");
        console.log("错误消息:", error.message);
        
        if (error.data) {
            console.log("错误数据:", error.data);
        }
        
        // 检查是否是 ProofInvalid
        if (error.data === "0x7fcdd1f4") {
            console.log("\n⚠️  错误是 ProofInvalid() - 证明本身无效");
            console.log("这说明:");
            console.log("  1. 证明格式正确（否则会是其他错误）");
            console.log("  2. 但验证计算失败");
            console.log("  3. 可能是 Setup 的验证密钥与证明不匹配");
        }
    }
}

main().catch(console.error);
