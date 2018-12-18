var Web3 		= 	require('web3');
var Contract 	= 	require('truffle-contract')
var abi 		= 	require("/home/jgm/goApp/src/receipt-trade/truffle_receipt/build/contracts/ReceiptTrade.json")
var grpc		=	require("grpc")	
var protoLoader = 	require("@grpc/proto-loader")

var web3 = new Web3();
web3.setProvider(new Web3.providers.HttpProvider("http://127.0.0.1:8545"));

//Receipt对象构造函数
function Receipt(receiptId, totalQty){
	this.receiptId 	= 	receiptId;
	this.totalQty 	=	totalQty;
}
//ReceiptQty对象构造函数
function ReceiptQty(totalQty, remainQty, frozenQty) {
	this.totalQty = totalQty;
	this.remainQty = remainQty;
	this.frozenQty = frozenQty;
}

var ReceiptTrade = Contract(abi);
ReceiptTrade.setProvider(web3.currentProvider);
var instance 	//合约实例
var accounts = web3.eth.accounts;	//获取管理员以太坊地址

async function main() {
	console.log(accounts[0])
	instance = await ReceiptTrade.deployed();	//创建合约实例

	var PROTO_PATH = "/home/jgm/goApp/src/receipt-trade/web/controllers/grpcPB/grpcPB.proto"
	var packageDefinition = protoLoader.loadSync(
		PROTO_PATH,{
			keepCase: 	true,
			longs:		String,
			enums:		String,
			defaults:	true,
			oneofs:		true
	});
	var rpc_proto	=	grpc.loadPackageDefinition(packageDefinition).grpcPB;

	var server = new grpc.Server();		//创建服务器
	server.addService(rpc_proto.RPCService.service,{
		InsertUserReceipt: 	insertUserReceipt,
		InsertUserFunds:	insertUserFunds,
		Trade: 				trade
	});
	server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());	//绑定地址
	server.start();			//启动服务器
	console.log("服务器已启动，正在监听.....");
	setOwner();
}

async function setOwner(){
	await instance.setOwner(accounts[0], {from:accounts[0], gas: 90000000})

	instance.setOwnerEv(function(err, rst){
		if (!err)
			console.log("previouseOwner: ",rst.args.previousOwner)
			console.log("newOwner: ", rst.args.newOwner)
	})
}

//grpc服务函数
async function insertUserReceipt(call, callback) {

	var userAddr 	= 	call.request.UserAddr;
	var ReceiptId	=	call.request.ReceiptId;
	var TotalQty	=	call.request.TotalQty;
	await instance.insertUserReceipt(userAddr, ReceiptId, TotalQty, {from: accounts[0], gas: 9000000});	//为用户添加仓单

	callback(null, {RstCode: 0, RstDetails : "智能合约，为用户添加仓单成功"});

	var length = await instance.getReceiptArrayLength(userAddr); //获取用户仓单的种类数量
	console.log(length.c[0])

	var myArray = new Array();
	for(var i = 0; i < length.c[0]; i++){
		rst = await instance.getReceipt(userAddr, i);
		myArray[i] = new Receipt(rst[0].c[0], rst[1].c[0]); 
	}

	console.log(userAddr)
	console.log(myArray);
}

//为用户添加资金
async function insertUserFunds(call, callback){
	var userAddr 	=	call.request.UserAddr;
	var totalFunds	=	call.request.TotalFunds;

	//调用智能合约，为用户添加资金
	await instance.insertUserFunds(userAddr, totalFunds, {from: accounts[0], gas: 9000000});

	callback(null, {RstCode: 0, RstDetails: "智能合约，为用户添加资金成功"})
}

//处理摘牌
async function trade(call, callback) {
	var tradeValues	= new Array(); //receiptId price QtySell nonceSell
	var tradeAddresses = new Array();
	var rs = new Array();
	var v = new Array();

	tradeValues[0] = call.request.ReceiptId;
	tradeValues[1] = call.request.Price;
	tradeValues[2] = call.request.QtySell;
	tradeValues[3] = call.request.NonceSell;
	tradeValues[4] = call.request.QtyBuy;
	tradeValues[5] = call.request.NonceBuy;

	tradeAddresses[0] = call.request.AddrSell
	tradeAddresses[1] = call.request.AddrBuy

	var sigSell		= 	call.request.SigSell;
	var sigBuy		= 	call.request.SigBuy;

	rs[0]	=	sigSell.slice(0, 66);
	rs[1]	=	'0x' + sigSell.slice(66, 130);
	v[0]	=  	await web3.toDecimal('0x' + sigSell.slice(130,132));
	rs[2]	=	sigBuy.slice(0, 66);
	rs[3]	=	'0x' + sigBuy.slice(66, 130);
	v[1]	=  	await web3.toDecimal('0x' + sigBuy.slice(130,132));

	//console.log(call.request);
	await instance.trade(tradeValues, tradeAddresses, v, rs, {from: accounts[0], gas: 9000000});


	var event = await instance.getErrCode();

	event.watch(function(err, rst){
		if (!err){
			console.log(rst.args.errCode.toNumber());
		}
	})

	event = await instance.getHash();
	event.watch(function(err, rst){
		if (!err){
			console.log("orderHash: ",rst.args.hash);
		}
	})
	var addrSell;
	event = await instance.getAddrSell();
	event.watch(async function(err, rst){
		if(!err){
			console.log("用户地址%s",rst.args.addr);
			addrSell = rst.args.addr;
		}
	})

	callback(null, {RstCode: 0, RstDetails: "调用智能合约成功"})
}
main();
