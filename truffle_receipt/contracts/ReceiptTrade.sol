pragma solidity ^0.4.23;

contract ReceiptTrade {
	struct Receipt {
		uint 	receiptId;
		uint	totalQty;
		uint 	remainQty;
		uint 	frozenQty;
	}

	struct Funds {
		uint totalFunds;
		uint remainFunds;
		uint frozenFunds;
	}

	mapping(address => Receipt[]) 	userReceipt;
	mapping(address => Funds)		userFunds;

	int _a;

	//为用户添加仓单
	function insertUserReceipt(address userAddr, uint receiptId, uint totalQty) public {
		Receipt memory r = Receipt(receiptId, totalQty, 0, 0);
		userReceipt[userAddr].push(r);
	}

	//获取用户的仓单
	function getReceipt(address userAddr, uint index) public view returns(
							uint receiptId, uint totalQty, uint remainQty, uint frozenQty){
		receiptId 	= 	userReceipt[userAddr][index].receiptId;
		totalQty 	=	userReceipt[userAddr][index].totalQty;
		remainQty 	=	userReceipt[userAddr][index].remainQty;
		frozenQty 	=	userReceipt[userAddr][index].frozenQty;

	}

	//获取用户的仓单数组的长度
	function getReceiptArrayLength(address userAddr) public view returns(uint) {
		return userReceipt[userAddr].length;
	}

	//为用户添加资金
	function insertUserFunds(address userAddr, uint totalFunds) public {
		Funds memory f = Funds(totalFunds, 0, 0);
		userFunds[userAddr] = f;
	}

	//获取用户资金
	function getFunds(address userAddr) public view returns(uint totalFunds, uint remainFunds, uint frozenFunds){
		totalFunds 	=	userFunds[userAddr].totalFunds;
		remainFunds =	userFunds[userAddr].remainFunds;
		frozenFunds = 	userFunds[userAddr].frozenFunds;
	}

	event getErrCode(int errCode);
	event getHash(bytes32 hash);
	//function trade(uint[6] tradeValues, address[2] tradeAddress, uint8[2] v, bytes32[4]rs) public  {
	//function trade(address addr, uint8 v, bytes32 r, bytes32 s){
	function trade(uint tradeValues, address addr, uint8 v, bytes32 r, bytes32 s){
		/*
		   tradeValues
		   	[0] receiptId
		   	[1] amountSell
		   	[2] price 
		   	[3] nonceSell 
		   	[4] amountBuy
		   	[5] nonceBuy
		   tradeAddress
		   	[0] addressSell
			[1] addressBuy
		*/

	   //bytes32 orderHash = sha3(tradeValues[0], tradeValues[1], tradeValues[2], tradeValues[3]);
	   //bytes32 orderHash = sha3("11061");
	   bytes32 orderHash = sha3(uint256(1),uint256(16));
	   getHash(orderHash);
	   bytes memory prefix = "\x19Ethereum Signed Message:\n32";
	   bytes32 hash = sha3(prefix, orderHash);
	   if ( ecrecover(hash, v, r, s) != addr)
		   getErrCode(25);
	   else 
		   getErrCode(26);
	}

}
