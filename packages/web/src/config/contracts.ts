type ContractAddresses = {
  UserRegistry: `0x${string}`;
  FeeManager: `0x${string}`;
  Deposit: `0x${string}`;
  ServiceManager: `0x${string}`;
  Payment: `0x${string}`;
  Oracle: `0x${string}`;
  TrafficCardNFT: `0x${string}`;
};

export const CONTRACTS: Record<number, ContractAddresses> = {
  // Hardhat localhost
  31337: {
    UserRegistry: "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512",
    FeeManager: "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9",
    ServiceManager: "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707",
    Payment: "0xa513E6E4b8f2a923D98304ec87F64353C4D5C853",
    Deposit: "0x8A791620dd6260079BF849Dc5567aDC3F2FdC318",
    Oracle: "0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e",
    TrafficCardNFT: "0x0DCd1Bf9A1b36cE34237eEaFef220932846BCD82",
  },
  // 0G Testnet (TODO: fill after testnet deploy)
  16601: {
    UserRegistry: "0x0000000000000000000000000000000000000000",
    FeeManager: "0x0000000000000000000000000000000000000000",
    Deposit: "0x0000000000000000000000000000000000000000",
    ServiceManager: "0x0000000000000000000000000000000000000000",
    Payment: "0x0000000000000000000000000000000000000000",
    Oracle: "0x0000000000000000000000000000000000000000",
    TrafficCardNFT: "0x0000000000000000000000000000000000000000",
  },
};

export function getContractAddress(
  chainId: number,
  name: keyof ContractAddresses
): `0x${string}` {
  const addresses = CONTRACTS[chainId];
  if (!addresses)
    throw new Error(`No contract addresses for chainId ${chainId}`);
  const addr = addresses[name];
  if (!addr || addr === "0x0000000000000000000000000000000000000000") {
    throw new Error(`Contract ${name} not deployed on chainId ${chainId}`);
  }
  return addr;
}
