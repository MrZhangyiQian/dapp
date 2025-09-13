const { ethers } = require("hardhat");

async function main() {
  const [deployer] = await ethers.getSigners();
  
  console.log("Deploying contracts with account:", deployer.address);
  
  // 部署ERC1820注册表
  const ERC1820Registry = await ethers.getContractFactory("ERC1820Registry");
  const erc1820 = await ERC1820Registry.deploy();
  await erc1820.deployed();
  console.log("ERC1820Registry deployed to:", erc1820.address);
  
  // 部署ERC777代币
  const initialSupply = ethers.utils.parseEther("1000000");
  const defaultOperators = []; // 初始为空，稍后添加操作员
  const ERC777Token = await ethers.getContractFactory("ERC777Token");
  const token = await ERC777Token.deploy(initialSupply, defaultOperators);
  await token.deployed();
  console.log("ERC777Token deployed to:", token.address);
  
  // 部署操作员合约
  const TokenOperator = await ethers.getContractFactory("TokenOperator");
  const operator = await TokenOperator.deploy(token.address);
  await operator.deployed();
  console.log("TokenOperator deployed to:", operator.address);
  
  // 部署接收合约
  const TokenReceiver = await ethers.getContractFactory("TokenReceiver");
  const receiver = await TokenReceiver.deploy();
  await receiver.deployed();
  console.log("TokenReceiver deployed to:", receiver.address);
  
  // 将操作员合约添加为默认操作员
  await token.authorizeOperator(operator.address);
  console.log("Operator authorized");
}

main()
  .then(() => process.exit(0))
  .catch(error => {
    console.error(error);
    process.exit(1);
  });