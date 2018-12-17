var ReceiptTrade = artifacts.require("./ReceiptTrade.sol");
var ReceiptMap 	=	artifacts.require("./ReceiptMap.sol")

module.exports = function(deployer) {
	deployer.deploy(ReceiptMap);
	deployer.link(ReceiptMap, ReceiptTrade);
  	deployer.deploy(ReceiptTrade);

};
