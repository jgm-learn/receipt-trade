var Web3 		= 	require('web3');
var Contract 	= 	require('truffle-contract')
var abi 		= 	require("/home/jgm/goApp/src/receipt-trade/truffle_receipt/build/contracts/ReceiptTrade.json")
var grpc		=	require("grpc")	
var protoLoader = 	require("@grpc/proto-loader")

var web3 = new Web3();
web3.setProvider(new Web3.providers.HttpProvider("http://127.0.0.1:8545"));



function Receipt(receiptId, totalQty, remainQty, frozenQty){
	this.receiptId 	= 	receiptId;
	this.totalQty 	=	totalQty;
	this.remainQty 	=	remainQty;
	this.frozenQty 	=	frozenQty;
}


function ReceiptQty(totalQty, remainQty, frozenQty) {
	this.totalQty = totalQty;
	this.remainQty = remainQty;
	this.frozenQty = frozenQty;
}


var ReceiptTrade = Contract(abi);
ReceiptTrade.setProvider(web3.currentProvider);
var instance 	//合约实例
var account = web3.eth.accounts[0];	//获取管理员以太坊地址

async  function start(){

	//console.log(account);

	await instance.insertUserReceipt(userAddr, 1, 100, {from: account, gas: 9000000});	//为用户添加仓单
	await instance.insertUserReceipt(userAddr, 2, 10, {from: account, gas: 9000000});	//为用户添加仓单

	var rst =	await instance.getReceipt(userAddr, 0);
	console.log(rst[0].c[0]);

	var length = await instance.getReceiptArrayLength(userAddr); //获取用户仓单的种类数量
	console.log(length.c[0])

	var myArray = new Array();
	for(var i = 0; i < length.c[0]; i++){
		rst = await instance.getReceipt(userAddr, i);
		myArray[i] = new Receipt(rst[0].c[0], rst[1].c[0], rst[2].c[0], rst[3].c[0]); 
	}

	console.log(myArray);


}

//start();

async function main() {
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
		InsertReceipt: InsertReceipt
	});
	server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());	//绑定地址
	server.start();			//启动服务器
}

//服务函数
async function InsertReceipt(call, callback) {

	var userAddr 	= 	call.request.UserAddr;
	var ReceiptId	=	call.request.ReceiptId;
	var TotalQty	=	call.request.TotalQty;
	await instance.insertUserReceipt(userAddr, ReceiptId, TotalQty, {from: account, gas: 9000000});	//为用户添加仓单

	callback(null, {Rst: 1});
}

main();
