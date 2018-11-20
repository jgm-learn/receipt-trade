import {default as Web3} from 'web3'
import {default as Contract} from 'truffle-contract'
import abi 	from "/home/jgm/goApp/src/receipt-trade/truffle_receipt/build/contracts/ReceiptTrade.json"

var ReceiptTrade = Contract(abi);
ReceiptTrade.setProvider(web3.currentProvider);
var instance 	//合约实例
//获取公钥
var account 
var publicKey = { pucKey: ""}

var userProperty = {
	totalFunds:		0,
	remainingFunds: 0,
	frozenFunds:	0
}

function Receipt(receiptId, totalQty, remainQty, frozenQty){
	this.receiptId 	= 	receiptId;
	this.totalQty 	=	totalQty;
	this.remainQty 	=	remainQty;
	this.frozenQty 	=	frozenQty;
}

async function start(){
	instance = await ReceiptTrade.deployed();	//创建合约实例
	//获取用户的以太坊地址
	web3.eth.getAccounts(function(err, rst){
		account = rst[0]
		$("#publicKey")[0].innerHTML = account
		publicKey.pucKey = account
		getFunds()
		getReceipt()
	});
}

function getFunds(){
	$.getJSON("http://222.22.64.80:8081/user/getFunds", publicKey, function(data){
		var colums = [{title: "总资金"}, {title: "可用资金"}, {title: "冻结资金"}];

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
		table.lastChild.appendChild(tr);
		var td = document.createElement("td");
		tr.appendChild(td);
		td.innerHTML = data.TotalFunds;
		var td = document.createElement("td");
		tr.appendChild(td);
		td.innerHTML = data.RemainingFunds;
		var td = document.createElement("td");
		tr.appendChild(td);
		td.innerHTML = data.FrozenFunds;
	})
}

async function getReceipt(){
	var colums = [{
			title: "仓单编号",
			name: "receiptId"
		},{
			title: "仓单总量",
			name: "totalQty"
		},{
			title: "剩余量",
			name: "remainQty"
		},{
			title: "冻结量",
			name: "frozenQty"
		}];	

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
	console.log(length.c[0])

	var myArray = new Array();
	for(var i = 0; i < length.c[0]; i++){
		var rst = await instance.getReceipt(account, i);
		myArray[i] = new Receipt(rst[0].c[0], rst[1].c[0], rst[2].c[0], rst[3].c[0]); 
	}
	console.log(account);
	console.log(myArray);

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




window.addEventListener('load', function(){
	if (typeof web3 != 'undefined'){
		console.log("Using MetaMask 3.");
		window.web3 = new Web3(web3.currentProvider);
	}else {
		alert("请安装MetaMask插件。")
	}

	start()
});
