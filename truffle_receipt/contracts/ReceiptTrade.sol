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

}
