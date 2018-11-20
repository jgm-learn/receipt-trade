pragma solidity ^0.4.23;

contract ReceiptTrade {
	struct Receipt {
		uint 	receiptId;
		uint	totalQty;
		uint 	remainQty;
		uint 	frozenQty;
	}

	mapping(address => Receipt[]) userReceipt;

	int _a;

	function insertUserReceipt(address userAddr, uint receiptId, uint totalQty) public {
		Receipt memory r = Receipt(receiptId, totalQty, 0, 0);
		userReceipt[userAddr].push(r);
	}

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
}
