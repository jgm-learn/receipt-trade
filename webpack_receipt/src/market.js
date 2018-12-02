import {default as Web3} from 'web3'


//账户地址和用户id
var account;
var userId;

/*
function receiptList(listId, receiptId, price, qtySell, qtyDeal, qtyRemain){
	this.listId 	=	listId;
	this.receiptId 	=	receiptId;
	this.price 		=	price;
	this.qtySell 	=	qtySell;
	this.qtyDeal 	= 	qtyDeal;
	this.qtyRemain  =   qtyRemain;
}
*/

function start(){
	web3.eth.getAccounts(function(err, rst){
		account = rst[0];
		$("#userAddr")[0].innerHTML = account;
		getUserId();
		getReceiptList();
	})
}


//获取userId
function getUserId(){
	$.getJSON("http://222.22.64.80:8081/user/getUserId", {userAddr: account}, function(data){
		userId = data.UserId;
		$("#userId")[0].innerHTML = userId;
	})
}

//展示所有的挂单数据

function getReceiptList(){
	$.getJSON("http://222.22.64.80:8081/market/getReceiptList", function(data){
		var colums = [	
					{title: "挂牌编号",	name: "Id"},
					{title: "仓单编号",	name: "ReceiptId"},
					{title: "价格",		name: "Price"},
					{title: "挂牌量",	name: "QtySell"},
					{title: "成交量",	name: "QtyDeal"},
					{title: "剩余量",	name: "QtyRemain"}];	
		//创建标题
		var table 	= 	document.getElementById("listTable")
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

		var myArray = new Array() 

		//将挂单数据填入表格
		for (var i = 0; i < data.length; i++){
			var tr = document.createElement("tr")
			table.lastChild.appendChild(tr)
			for(var j = 0; j < colums.length; j++){
				var td = document.createElement("td")
				tr.appendChild(td)
				td.innerHTML = data[i][colums[j].name]
			}
		}
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
