pragma solidity ^0.4.23;

library ReceiptMap{
	struct itmap{
		mapping(uint256 => QtyTotal) data;
		ReceiptId[] keys;
		uint size;
	}

	struct ReceiptId {
		uint256 receiptId;
		bool	deleted;
	}
	struct QtyTotal{
		uint keyIndex;
		uint256 qtyTotal;
	}

	function insert(itmap storage self, uint256 receiptId, uint256 qtyTotal) returns(bool replaced){
		uint keyIndex = self.data[receiptId].keyIndex;
		self.data[receiptId].qtyTotal += qtyTotal;

		if (keyIndex > 0) {
			return true;
		} else {
			keyIndex = self.keys.length++;
			self.data[receiptId].keyIndex = keyIndex + 1;
			self.keys[keyIndex].receiptId = receiptId;
			self.size++;
			return false;
		}
	}

	function minus(itmap storage self, uint256 receiptId, uint256 qtyBuy) returns (bool){
		if (self.data[receiptId].qtyTotal < qtyBuy)
			return false;
		else {
			self.data[receiptId].qtyTotal -= qtyBuy;
			return true;
		}
	}
}
