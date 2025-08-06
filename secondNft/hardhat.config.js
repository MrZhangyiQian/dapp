require("@nomicfoundation/hardhat-toolbox");
require("@openzeppelin/hardhat-upgrades"); // 添加这一行
require("hardhat-deploy");
require("dotenv").config();

const SEPOLIA_RPC_URL = process.env.SPEOLIA_URL;
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const ETHERSCAN_API_KEY = process.env.ETHERSCAN_API_KEY;

module.exports = {
  defaultNetwork: "hardhat",
  networks: {
    hardhat: {
      chainId: 31337,
      allowUnlimitedContractSize: true,
      blockConfirmations: 1,
    },
    sepolia: {
      url: SEPOLIA_RPC_URL,
      accounts: [PRIVATE_KEY],
      chainId: 11155111,
      blockConfirmations: 6,
    },
  },
  solidity: "0.8.28",
  namedAccounts: {
    deployer: {
      default: 0,
    },
    seller: {
      default: 1,
    },
    bidder1: {
      default: 2,
    },
    bidder2: {
      default: 3,
    },
  },
  etherscan: {
    apiKey: ETHERSCAN_API_KEY,
  },
  mocha: {
    reporter: 'mochawesome',
    reporterOptions: {
      reportDir: 'test-reports', // 报告生成目录
      overwrite: false,          // 是否覆盖之前的报告
      html: true,                // 生成HTML报告
      json: true,                // 生成JSON报告
      timestamp: 'yyyy-mm-dd_HH-MM-ss' // 时间戳，用于避免覆盖
    }
  },
  sourcify: {
    enabled: false
  },
  settings: {
    debug: {
      revertStrings: "verbose", // 或者 "strip", "debug", "verboseDebug"
    },
    outputSelection: {
      "*": {
        "*": ["storageLayout", "metadata", "abi", "bin", "bin-runtime"]
      }
    }
  },
};