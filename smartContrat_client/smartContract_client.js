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
		InsertUserReceipt: 	insertUserReceipt,
		InsertUserFunds:	insertUserFunds
	});
	server.bind('0.0.0.0:50051', grpc.ServerCredentials.createInsecure());	//绑定地址
	server.start();			//启动服务器
	console.log("服务器已启动，正在监听.....")
}

//服务函数
async function insertUserReceipt(call, callback) {

	var userAddr 	= 	call.request.UserAddr;
	var ReceiptId	=	call.request.ReceiptId;
	var TotalQty	=	call.request.TotalQty;
	await instance.insertUserReceipt(userAddr, ReceiptId, TotalQty, {from: account, gas: 9000000});	//为用户添加仓单

	callback(null, {RstCode: 0, RstDetails : "智能合约，为用户添加仓单成功"});

	var length = await instance.getReceiptArrayLength(userAddr); //获取用户仓单的种类数量
	console.log(length.c[0])

	var myArray = new Array();
	for(var i = 0; i < length.c[0]; i++){
		rst = await instance.getReceipt(userAddr, i);
		myArray[i] = new Receipt(rst[0].c[0], rst[1].c[0], rst[2].c[0], rst[3].c[0]); 
	}

	console.log(userAddr)
	console.log(myArray);
}

//为用户添加资金
async function insertUserFunds(call, callback){
	var userAddr 	=	call.request.UserAddr;
	var totalFunds	=	call.request.TotalFunds;

	//调用智能合约，为用户添加资金
	await instance.insertUserFunds(userAddr, totalFunds, {from: account, gas: 9000000});

	callback(null, {RstCode: 0, RstDetails: "智能合约，为用户添加资金成功"})

}
main();
