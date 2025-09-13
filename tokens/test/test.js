const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("ERC777 Token System", function () {
  let token, operator, receiver, erc1820Registry;
  let owner, user1, user2;

  const deployContracts = async () => {
    // 部署ERC1820注册表
    const ERC1820Registry = await ethers.getContractFactory("ERC1820Registry");
    const erc1820RegistryInstance = await ERC1820Registry.deploy();
    await erc1820RegistryInstance.waitForDeployment();
    
    // 确保获取到地址
    const registryAddress = await erc1820RegistryInstance.getAddress();
    console.log("Registry address:", registryAddress);

    // ERC777代币参数
    const initialSupply = ethers.parseEther("1000");
    const defaultOperators = [];

    // 部署ERC777代币
    const ERC777Token = await ethers.getContractFactory("ERC777Token");
    token = await ERC777Token.deploy(
      initialSupply,
      defaultOperators,
      registryAddress
    );
    await token.waitForDeployment();

    // 部署操作员合约
    const TokenOperator = await ethers.getContractFactory("TokenOperator");
    operator = await TokenOperator.deploy(await token.getAddress());
    await operator.waitForDeployment();

    // 部署接收合约
    const TokenReceiver = await ethers.getContractFactory("TokenReceiver");
    receiver = await TokenReceiver.deploy(registryAddress);
    await receiver.waitForDeployment();

    // 授权操作员
    await token.authorizeOperator(await operator.getAddress());
    
    // 返回实例以便在外部使用
    return erc1820RegistryInstance;
  };

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();
    
    // 部署合约
    erc1820Registry = await deployContracts();

    // 分配初始代币
    await token.transfer(user1.address, ethers.parseEther("100"));
  });

  it("应正确初始化代币名称和符号", async function () {
    expect(await token.name()).to.equal("ERC777Token");
    expect(await token.symbol()).to.equal("ERC777");
  });

  describe("基本功能测试", function () {
    it("应该正确部署合约并分配初始供应量", async function () {
      expect(await token.totalSupply()).to.equal(ethers.parseEther("1000"));
      expect(await token.balanceOf(owner.address)).to.equal(ethers.parseEther("900"));
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("100"));
    });

    it("应该允许用户发送代币", async function () {
      const amount = ethers.parseEther("10");
      
      // 用户1发送代币给用户2
      await token.connect(user1).send(user2.address, amount, "0x");
      
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("90"));
      expect(await token.balanceOf(user2.address)).to.equal(amount);
    });

    it("应该允许用户授权和撤销操作员", async function () {
      // 授权操作员
      await token.connect(user1).authorizeOperator(await operator.getAddress());
      expect(await token.isOperatorFor(await operator.getAddress(), user1.address)).to.be.true;
      
      // 撤销操作员
      await token.connect(user1).revokeOperator(await operator.getAddress());
      expect(await token.isOperatorFor(await operator.getAddress(), user1.address)).to.be.false;
    });
  });

  describe("操作员功能测试", function () {
    it("应该允许操作员代表用户发送代币", async function () {
      const amount = ethers.parseEther("10");
      
      // 用户1授权操作员
      await token.connect(user1).authorizeOperator(await operator.getAddress());
      
      // 操作员代表用户1发送代币给接收合约
      await operator.operatorSend(
        user1.address,
        await receiver.getAddress(),
        amount,
        "0x",
        "operator data"
      );
      
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("90"));
      expect(await token.balanceOf(await receiver.getAddress())).to.equal(amount);
    });

    it("应该阻止未授权的操作员代表用户发送代币", async function () {
      const amount = ethers.parseEther("10");
      
      // 尝试使用未授权的操作员发送代币
      await expect(
        operator.operatorSend(
          user1.address,
          await receiver.getAddress(),
          amount,
          "0x",
          "operator data"
        )
      ).to.be.revertedWith("ERC777: caller is not operator");
    });

    it("应该允许操作员代表用户销毁代币", async function () {
      const amount = ethers.parseEther("10");
      
      // 用户1授权操作员
      await token.connect(user1).authorizeOperator(await operator.getAddress());
      
      // 操作员代表用户1销毁代币
      await operator.operatorBurn(
        user1.address,
        amount,
        "0x",
        "operator data"
      );
      
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("90"));
      expect(await token.totalSupply()).to.equal(ethers.parseEther("990"));
    });
  });

  describe("接收钩子功能测试", function () {
    it("应该触发 tokensReceived 钩子函数", async function () {
      const amount = ethers.parseEther("10");
      
      // 发送代币给接收合约
      await token.connect(user1).send(
        await receiver.getAddress(),
        amount,
        "user data"
      );
      
      // 验证接收合约的余额
      // 由于钩子函数会返回10%的代币，所以余额应为90%
      const expectedBalance = amount.mul(9).div(10);
      expect(await token.balanceOf(await receiver.getAddress())).to.equal(expectedBalance);
      
      // 验证用户1的余额变化（应收到10%的返回）
      expect(await token.balanceOf(user1.address)).to.equal(
        ethers.parseEther("100").sub(amount).add(amount.div(10))
      );
    });

    it("应该发出 TokensReceived 事件", async function () {
      const amount = ethers.parseEther("10");
      
      // 发送代币给接收合约
      const tx = await token.connect(user1).send(
        await receiver.getAddress(),
        amount,
        "user data"
      );
      
      const receipt = await tx.wait();
      const event = receipt.events.find(e => e.event === "Sent");
      
      expect(event).to.not.be.undefined;
      expect(event.args.operator).to.equal(user1.address);
      expect(event.args.from).to.equal(user1.address);
      expect(event.args.to).to.equal(await receiver.getAddress());
      expect(event.args.amount).to.equal(amount);
    });

    it("应该正确处理没有实现接收接口的地址", async function () {
      const amount = ethers.parseEther("10");
      
      // 发送代币给普通地址（未实现接收接口）
      await token.connect(user1).send(
        user2.address,
        amount,
        "user data"
      );
      
      // 验证余额变化
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("90"));
      expect(await token.balanceOf(user2.address)).to.equal(amount);
    });
  });

  describe("ERC20 兼容性测试", function () {
    it("应该支持 ERC20 transfer 函数", async function () {
      const amount = ethers.parseEther("10");
      
      // 使用 ERC20 transfer 函数
      await token.connect(user1).transfer(user2.address, amount);
      
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("90"));
      expect(await token.balanceOf(user2.address)).to.equal(amount);
    });

    it("应该支持 ERC20 transferFrom 函数", async function () {
      const amount = ethers.parseEther("10");
      
      // 用户1批准操作员花费代币
      await token.connect(user1).approve(await operator.getAddress(), amount);
      
      // 操作员使用 transferFrom
      await token.connect(operator).transferFrom(
        user1.address,
        user2.address,
        amount
      );
      
      expect(await token.balanceOf(user1.address)).to.equal(ethers.parseEther("90"));
      expect(await token.balanceOf(user2.address)).to.equal(amount);
    });
  });

  describe("边界和安全测试", function () {
    it("应该防止发送到零地址", async function () {
      const amount = ethers.parseEther("10");
      
      await expect(
        token.connect(user1).send(
          ethers.constants.AddressZero,
          amount,
          "0x"
        )
      ).to.be.revertedWith("ERC777: send to zero address");
    });

    it("应该防止从零地址发送", async function () {
      const amount = ethers.parseEther("10");
      
      await expect(
        token.send(
          user2.address,
          amount,
          "0x"
        )
      ).to.be.reverted;
    });

    it("应该防止发送超过余额的代币", async function () {
      const amount = ethers.parseEther("1000");
      
      await expect(
        token.connect(user1).send(
          user2.address,
          amount,
          "0x"
        )
      ).to.be.revertedWith("ERC777: insufficient balance");
    });

    it("应该防止未授权的操作员销毁代币", async function () {
      const amount = ethers.parseEther("10");
      
      await expect(
        operator.operatorBurn(
          user1.address,
          amount,
          "0x",
          "operator data"
        )
      ).to.be.revertedWith("ERC777: caller is not operator");
    });
  });
});