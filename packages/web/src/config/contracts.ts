import hardhatDeployment from "./deployments/hardhat.json";
import arbitrumSepoliaDeployment from "./deployments/arbitrum_sepolia.json";

// 单一出口：合约地址 / USDT / 精度统一从 deployments/<net>.json 读，杜绝手抄漂移。
// 31337 取 deployments/hardhat.json；421614 取 deployments/arbitrum_sepolia.json
// Arbitrum Sepolia 部署地址来自 packages/contracts/deployments/arbitrum_sepolia.json。

type ContractName =
  | "UserRegistry"
  | "FeeManager"
  | "Deposit"
  | "ServiceManager"
  | "Payment"
  | "Oracle"
  | "TrafficCardNFT";

export type ContractAddresses = Record<ContractName, `0x${string}`> & {
  usdt: `0x${string}`;
  usdtDecimals: number;
};

type Deployment = {
  proxies: Record<string, string>;
  usdt: string;
  usdtDecimals: number;
};

function fromDeployment(d: Deployment): ContractAddresses {
  return {
    UserRegistry: d.proxies.UserRegistry as `0x${string}`,
    FeeManager: d.proxies.FeeManager as `0x${string}`,
    Deposit: d.proxies.Deposit as `0x${string}`,
    ServiceManager: d.proxies.ServiceManager as `0x${string}`,
    Payment: d.proxies.Payment as `0x${string}`,
    Oracle: d.proxies.Oracle as `0x${string}`,
    TrafficCardNFT: d.proxies.TrafficCardNFT as `0x${string}`,
    usdt: d.usdt as `0x${string}`,
    usdtDecimals: d.usdtDecimals,
  };
}

export const CONTRACTS: Record<number, ContractAddresses> = {
  31337: fromDeployment(hardhatDeployment as Deployment),
  421614: fromDeployment(arbitrumSepoliaDeployment as Deployment),
};

const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

export function getContractAddress(
chainId: number,
name: ContractName
): `0x${string}` {
const addresses = CONTRACTS[chainId];
if (!addresses) {
throw new Error(`No contract addresses for chainId ${chainId}`);
}

const addr = addresses[name];
if (
!addr ||
addr.toLowerCase() === ZERO_ADDRESS.toLowerCase()
) {
throw new Error(`Contract ${name} not deployed on chainId ${chainId}`);
}

return addr;
}

/** 不抛错地检查合约地址是否已配置；未配置时返回 undefined。 */
export function getOptionalContractAddress(
chainId: number,
name: ContractName
): `0x${string}` | undefined {
try {
return getContractAddress(chainId, name);
} catch {
return undefined;
}
}

/** 不抛错地检查 USDT 地址是否已配置；未配置时返回 undefined。 */
export function getOptionalUsdtAddress(
chainId: number
): `0x${string}` | undefined {
try {
return getUsdt(chainId);
} catch {
return undefined;
}
}

/** 获取当前网络的 USDT 地址；未配置或为零地址时抛出错误。 */
export function getUsdt(chainId: number): `0x${string}` {
  const addresses = CONTRACTS[chainId];

  if (!addresses) {
    throw new Error(`No contract addresses for chainId ${chainId}`);
  }

  const addr = addresses.usdt;

  if (!addr || addr.toLowerCase() === ZERO_ADDRESS.toLowerCase()) {
    throw new Error(`USDT not deployed on chainId ${chainId}`);
  }

  return addr;
}

/** 获取当前网络的 USDT 小数位数。 */
export function getUsdtDecimals(chainId: number): number {
  const addresses = CONTRACTS[chainId];

  if (!addresses) {
    throw new Error(`No contract addresses for chainId ${chainId}`);
  }

  return addresses.usdtDecimals;
}
