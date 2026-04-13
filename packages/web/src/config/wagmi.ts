import { getDefaultConfig } from "@rainbow-me/rainbowkit";
import { http } from "wagmi";
import { zgTestnet } from "./chains";

export const wagmiConfig = getDefaultConfig({
  appName: "LinkWorld",
  projectId: "linkworld-dev",
  chains: [zgTestnet],
  transports: {
    [zgTestnet.id]: http(),
  },
});
