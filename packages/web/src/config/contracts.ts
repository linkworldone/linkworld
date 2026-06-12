import hardhatDeployment from "./deployments/hardhat.json";
import arbitrumSepoliaDeployment from "./deployments/arbitrum_sepolia.json";

// 单一出口：合约地址 / USDT / 精度统一从 deployments/<net>.json 读，杜绝手抄漂移。
// 31337 取 deployments/hardhat.json；421614 取 deployments/arbitrum_sepolia.json
//（真·上链未执行 → 当前为零地址占位，待合约上链回填）。

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
  if (!addresses)
    throw new Error(`No contract addresses for chainId ${chainId}`);
  const addr = addresses[name];
  if (!addr || addr === ZERO_ADDRESS) {
    throw new Error(`Contract ${name} not deployed on chainId ${chainId}`);
  }
  return addr;
}

/** USDT(ERC20) 合约地址；零地址（未部署）抛错保护。 */
export function getUsdt(chainId: number): `0x${string}` {
  const addresses = CONTRACTS[chainId];
  if (!addresses)
    throw new Error(`No contract addresses for chainId ${chainId}`);
  if (!addresses.usdt || addresses.usdt === ZERO_ADDRESS) {
    throw new Error(`USDT not deployed on chainId ${chainId}`);
  }
  return addresses.usdt;
}

/** USDT 精度（本轮 MockUSDT = 6）。从 deployments 读，不硬编码。 */
export function getUsdtDecimals(chainId: number): number {
  const addresses = CONTRACTS[chainId];
  if (!addresses)
    throw new Error(`No contract addresses for chainId ${chainId}`);
  return addresses.usdtDecimals;
}
