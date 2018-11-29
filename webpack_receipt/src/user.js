import {default as Web3} from 'web3'
import {default as Contract} from 'truffle-contract'
import abi 	from "/home/jgm/goApp/src/receipt-trade/truffle_receipt/build/contracts/ReceiptTrade.json"

var ReceiptTrade = Contract(abi);
ReceiptTrade.setProvider(web3.currentProvider);
var instance 	//合约实例
//获取公钥
var account;
var userId;

var userProperty = {
	totalFunds:		0,
	remainingFunds: 0,
	frozenFunds:	0
}

//Receipt对象
function Receipt(receiptId, totalQty, remainQty, frozenQty){
	this.receiptId 	= 	receiptId;
	this.totalQty 	=	totalQty;
	this.remainQty 	=	remainQty;
	this.frozenQty 	=	frozenQty;
}
//卖方订单
var orderSell = {
	UserId:		0,
	ReceiptId:	0,
	Price:		0,
	QtySell:	0,
	NonceSell:  0,
	SigSell:  	"",
	AddrSell:  	""
}

async function start(){
	instance = await ReceiptTrade.deployed();	//创建合约实例
	//获取用户的以太坊地址
	web3.eth.getAccounts(function(err, rst){
		account = rst[0]
		$("#publicKey")[0].innerHTML = account
		getUserId();
		getFunds()
		getReceipt()
		setNonce();
		listTrade();
	});
}

//获取userId
function getUserId(){
	$.getJSON("http://222.22.64.80:8081/user/getUserId", {userAddr: account}, function(data){
		userId = data.UserId;
		$("#userId")[0].innerHTML = userId;
	})
}
//从智能合约获取资金
async function getFunds(){
	//$.getJSON("http://222.22.64.80:8081/user/getFunds", publicKey, function(data){
	var rst = await instance.getFunds(account);
		var colums = [{title: "总资金"}, {title: "可用资金"}, {title: "冻结资金"}];

		//构建表格
		var table 	= 	document.getElementById("fundsTable")
		var thead 	=	document.createElement("thead")
		table.appendChild(thead)
		var tr 		= document.createElement("tr")
		thead.appendChild(tr)
		var tbody	= document.createElement("tbody")
		table.appendChild(tbody)
		for( var i=0; i < colums.length; i++){
			var th = document.createElement("th")
			tr.appendChild(th)
			th.innerHTML = colums[i].title
		}

	var tr = document.createElement("tr");
	table.lastChild.appendChild(tr)
	for(var j = 0; j < colums.length; j++){
			var td = document.createElement("td")
			tr.appendChild(td)
			td.innerHTML = rst[j].c[0]
	}
}

//从智能合约获取仓单
async function getReceipt(){
	var colums = [	{title: "仓单编号",	name: "receiptId"},
					{title: "仓单总量",	name: "totalQty"},
					{title: "剩余量",	name: "remainQty"},
					{title: "冻结量",	name: "frozenQty"}];	

	var table 	= 	document.getElementById("simpleTable")
	var thead 	=	document.createElement("thead")
	table.appendChild(thead)
		var tr 		= document.createElement("tr")
		thead.appendChild(tr)
	var tbody	= document.createElement("tbody")
	table.appendChild(tbody)
	for( var i=0; i < colums.length; i++){
		var th = document.createElement("th")
		tr.appendChild(th)
		th.innerHTML = colums[i].title
	}

	//从智能合约中获取用户的所有仓单
	var length = await instance.getReceiptArrayLength(account); //获取用户仓单的种类数量

	var myArray = new Array();
	for(var i = 0; i < length.c[0]; i++){
		var rst = await instance.getReceipt(account, i);
		myArray[i] = new Receipt(rst[0].c[0], rst[1].c[0], rst[2].c[0], rst[3].c[0]); 
	}
	console.log(account);

	//将数据填入表格
	for (var i = 0; i < myArray.length; i++){
		var tr = document.createElement("tr")
		table.lastChild.appendChild(tr)
		for(var j = 0; j < colums.length; j++){
			var td = document.createElement("td")
			tr.appendChild(td)
			td.innerHTML = myArray[i][colums[j].name]
		}
	}
}

async function setNonce() {
	$("#bt_setNonce").on('click', async function(){
		await instance.setNonce(6, {from: account, gasLimit: 9000000});
	})
}

function listTrade(){
	$("#bt_listTrade").on('click', function(){
		orderSell.UserId =	userId;
		orderSell.ReceiptId = parseInt($("#receiptId").val());	
		orderSell.Price 	= parseInt($("#price").val());
		orderSell.QtySell 	= parseInt($("#QtySell").val());
		$.getJSON("http://222.22.64.80:8081/user/getUserNonce", {userAddr: account}, async function(data){
			var nonceLast = data.Nonce;
			console.log("nonceLast: %d", nonceLast)
			orderSell.NonceSell = nonceLast + 1;
			var orderHash = await web3.utils.soliditySha3(orderSell.ReceiptId, orderSell.Price, orderSell.QtySell, orderSell.NonceSell);
			console.log("orderHash");
			console.log(orderHash);
			orderSell.SigSell = await web3.eth.sign(orderHash, account)
			console.log("SigSell");
			console.log(orderSell.SigSell);
			orderSell.AddrSell = account;
			//序列化并发送
			$.ajax({
				type:			'post',
				url:			"http://222.22.64.80:8081/user/listTrade",
				data:			JSON.stringify(orderSell),
				contentType:	"application/json",
				success:		function(data, state, xhr){
					console.log("orderSell 发送成功");
					console.log(state);
				},
				error:			function(xhr, state, error){
					console.log("orderSell 发送失败");
					console.log(error.toString());
				}
			});
		});
	});
}

window.addEventListener('load', function(){
	if (typeof web3 != 'undefined'){
		console.log("Using MetaMask 3.");
		window.web3 = new Web3(web3.currentProvider);
	}else {
		alert("请安装MetaMask插件。")
	}

	start()
});
