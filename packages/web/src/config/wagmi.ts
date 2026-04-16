import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { http } from "wagmi";
import { hardhatLocal, zgTestnet } from "./chains";

const isLocalChain = import.meta.env.VITE_CHAIN_ID === "31337";
const chain = isLocalChain ? hardhatLocal : zgTestnet;

export const wagmiConfig = getDefaultConfig({
  appName: "LinkWorld",
  projectId: "21fef48091f12692cad574a6f7753643",
  chains: [chain],
  transports: {
    [chain.id]: http(),
  },
});
