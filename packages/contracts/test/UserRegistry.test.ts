import { expect } from "chai";
import { ethers } from "hardhat";

describe("UserRegistry", function () {
  it("should register a user and mint NFT", async function () {
    const [owner, user1] = await ethers.getSigners();

    const UserRegistry = await ethers.getContractFactory("UserRegistry");
    const registry = await UserRegistry.deploy();

    await registry.connect(user1).register("user1@example.com");

    const info = await registry.getUserInfo(user1.address);
    expect(info.email).to.equal("user1@example.com");
    expect(info.isActive).to.be.true;
    expect(await registry.ownerOf(0)).to.equal(user1.address);
  });

  it("should reject duplicate registration", async function () {
    const [, user1] = await ethers.getSigners();

    const UserRegistry = await ethers.getContractFactory("UserRegistry");
    const registry = await UserRegistry.deploy();

    await registry.connect(user1).register("user1@example.com");
    await expect(
      registry.connect(user1).register("user1@example.com")
    ).to.be.revertedWith("Already registered");
  });
});
