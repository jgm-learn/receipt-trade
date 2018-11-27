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

	trade();
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

//处理摘牌
async function trade() {
	//var tradeValues 	= 	[1,10,6,1];
	var tradeValues 	= 	26;
	var tradeAddress 	=	[account];

	var orderInfo 		= 	"010";
	var hash 			= 	web3.sha3('1','16');
	var hash2			= 	web3.sha3('116');
	//var hash 			= 	web3.utils.soliditySha3(26);
	var sig 			= 	web3.eth.sign(account, hash);
	var r				=	sig.slice(0, 66);
	var s				=	'0x' + sig.slice(66, 130);
	var v 				=	web3.toDecimal(sig.slice(130,132)) + 27;
	var vArray			=	[v];
	var rsArray			=	[{r,s}]

	console.log(hash);
	console.log(hash2);

	 //await instance.trade(tradeValues, tradeAddress, vArray, rsArray, {from: account, gas: 9000000});
	 //await instance.trade(account, v, r, s,  {from: account, gas: 9000000});
	 await instance.trade(tradeValues, account, v, r, s,  {from: account, gas: 9000000});

	var event = await instance.getErrCode();

	//console.log(event);

	event.watch(function(err, rst){
		if (!err){
			console.log(rst.args.errCode.c[0]);
		}
	})

	event = await instance.getHash();
	event.watch(function(err, rst){
		if (!err){
			console.log(rst.args.hash);
		}
	})

}
main();
