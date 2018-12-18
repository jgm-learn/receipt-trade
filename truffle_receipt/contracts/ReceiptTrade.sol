pragma solidity ^0.4.23;

import "./ReceiptMap.sol";

contract ReceiptTrade {
	using ReceiptMap for ReceiptMap.itmap;

	mapping(address => ReceiptMap.itmap) 			userReceipt;
	mapping(address => uint256)						userFunds;
	mapping(address => mapping(uint256 => uint256))	orderFill; //地址 => (nonceSell => 成交量)
	address userSell;

	address owner;
	event setOwnerEv(address indexed previousOwner, address indexed newOwner);

	//求和	
	function safeAdd(uint a, uint b) returns (uint) {
		uint c = a + b;
		assert(c >= a && c >= b);
		return c;
	}

	//构造函数
	function ReceiptTrade(){
		owner = msg.sender;
	}

	modifier onlyOwner{
		assert(msg.sender == owner);
		_;
	}
	//设置新的管理员
	function setOwner(address newOwner) onlyOwner {
		setOwnerEv(owner, newOwner);
		owner = newOwner;
	}

	//为用户添加仓单
	function insertUserReceipt(address addr, uint receiptId, uint qtyTotal) public {
		userReceipt[addr].insert(receiptId, qtyTotal);
	}

	//获取用户的仓单
	function getReceipt(address userAddr, uint index) public view returns(
							uint receiptId, uint totalQty){
		receiptId 	= 	userReceipt[userAddr].keys[index].receiptId;
		totalQty 	=	userReceipt[userAddr].data[receiptId].qtyTotal;
	}

	//获取用户的仓单数组的长度
	function getReceiptArrayLength(address userAddr) public view returns(uint) {
		return userReceipt[userAddr].keys.length;
	}

	//为用户添加资金
	function insertUserFunds(address userAddr, uint256 totalFunds) public {
		userFunds[userAddr] = totalFunds;
	}

	//获取用户资金
	function getFunds(address userAddr) public view returns(uint256 totalFunds){
		totalFunds 	=	userFunds[userAddr];
	}

	//取消挂单
	function cancellation(uint256 nonceSell, uint256 qtyRemain){
		orderFill[msg.sender][nonceSell] += qtyRemain;
	}

	event getErrCode(int errCode);
	event getHash(bytes32 hash);
	event getAddrSell(address addr);
	event getNonceSell(uint256 nonceSell);

	function trade(uint256[6] tradeValues, address[2] tradeAddress, uint8[2] v, bytes32[4] rs){
		/*
		   tradeValues
		   	[0] receiptId
		   	[1] price 
		   	[2] qtySell
		   	[3] nonceSell 
		   	[4] qtyBuy
		   	[5] nonceBuy
		   tradeAddress
		   	[0] addressSell
			[1] addressBuy
		*/	
	   bytes32 orderSellHash 	= sha3(tradeValues[0], tradeValues[1], tradeValues[2], tradeValues[3], tradeAddress[0]);
	   bytes32 orderBuyHash 	= sha3(orderSellHash, tradeValues[4], tradeValues[5], tradeAddress[1]);
	   getHash(orderBuyHash);

	   if ( ecrecover(orderSellHash, v[0], rs[0], rs[1]) != tradeAddress[0]) {
		   getErrCode(-1);
		   throw;
	   }
	   if ( ecrecover(orderBuyHash, v[1], rs[2], rs[3]) != tradeAddress[1])  {
		   getErrCode(-2);
		   throw;
	   }
	   //判断挂单剩余量是否足够，卖方可通过管理剩余量来使订单失效
	   if ( safeAdd(orderFill[tradeAddress[0]][tradeValues[3]], tradeValues[4])  > tradeValues[2]) {
		   getErrCode(-3);
		   throw;
	   }
	   //更新成交量
	   orderFill[tradeAddress[0]][tradeValues[3]] = safeAdd( orderFill[tradeAddress[0]][tradeValues[3]], tradeValues[4] );

	   uint256 payment = tradeValues[4] * tradeValues[1]; 	//应付款
	   if ( userFunds[tradeAddress[1]] <  payment ){		//判断资金是否足够
		   getErrCode(-4);
		   throw;
	   }	

	   if ( !userReceipt[tradeAddress[0]].minus(tradeValues[0], tradeValues[4]) ){ //减少卖方仓单
			getErrCode(-5);
			throw;
	   }
	   userReceipt[tradeAddress[1]].insert(tradeValues[0], tradeValues[4]); //增加买方仓单

	   userFunds[tradeAddress[1]] -= payment;
	   userFunds[tradeAddress[0]] += payment;

	   getErrCode(0);
	}
}
