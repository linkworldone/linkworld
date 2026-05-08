import { ethers } from "hardhat";
import axios from "axios";

/**
 * 从 Binance API 获取 USDT 活期利率并更新到 Oracle 合约
 * 需要先安装 axios: npm install axios
 */
async function main() {
    // 获取 Oracle 合约实例
    const Oracle = await ethers.getContractFactory("Oracle");
    const oracle = await Oracle.attach(
        // 替换为实际部署的 Oracle 合约地址
        process.env.ORACLE_ADDRESS || "0xYourOracleContractAddress"
    );

    console.log("正在从 Binance API 获取利率...");

    try {
        // 从 Binance API 获取活期理财产品列表
        const response = await axios.get(
            "https://api.binance.com/sapi/v1/lending/daily/product/list",
            {
                params: {
                    status: "ALL",
                },
            }
        );

        // 查找 USDT 产品
        const usdtProduct = response.data.data.find(
            (p: any) => p.asset === "USDT" && p.status === "PURCHASING"
        );

        if (!usdtProduct) {
            throw new Error("未找到可用的 USDT 理财产品");
        }

        // 计算年化利率（日利率 * 365）
        const dailyRate = parseFloat(usdtProduct.dailyInterestRate);
        const annualRate = dailyRate * 365;

        // 转换为基点（如 2.5% -> 250 basis points）
        const rateBasisPoints = Math.round(annualRate * 100);

        console.log(`当前 Binance USDT 年化利率: ${(annualRate * 100).toFixed(2)}%`);
        console.log(`转换为基点: ${rateBasisPoints}`);

        // 更新到 Oracle 合约
        console.log("正在更新利率到合约...");
        const tx = await oracle.updateInterestRateFromAPI(rateBasisPoints);
        await tx.wait();

        console.log("利率更新成功！交易哈希:", tx.hash);
    } catch (error) {
        console.error("获取或更新利率失败:", error);
        process.exit(1);
    }
}

main().catch((error) => {
    console.error(error);
    process.exitCode = 1;
});