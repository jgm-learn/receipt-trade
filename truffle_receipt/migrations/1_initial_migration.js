var ReceiptTrade = artifacts.require("./ReceiptTrade.sol");

module.exports = function(deployer) {
  deployer.deploy(ReceiptTrade);
};
