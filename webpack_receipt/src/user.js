import {default as Web3} from 'web3'
import {default as contract} from 'truffle-contract'

var userProperty = {
	totalFunds:		0,
	remainingFunds: 0,
	frozenFunds:	0
}

//获取公钥
var account 
var publicKey = { pucKey: ""}
function start(){
	web3.eth.getAccounts(function(err, rst){
		account = rst[0]
		$("#publicKey")[0].innerHTML = account
		publicKey.pucKey = account
		getFunds()
		getReceipt()
	})
}

function getFunds(){
	$.getJSON("http://222.22.64.80:8081/user/getFunds", publicKey, function(data){
		$("#totalFunds")[0].innerHTML 		= 	data.TotalFunds
		$("#remainingFunds")[0].innerHTML 	= 	data.RemainingFunds
		$("#frozenFunds")[0].innerHTML 		= 	data.FrozenFunds
	})
}

function getReceipt(){
	$.getJSON("http://222.22.64.80:8081/user/getReceipt", publicKey, function(data){
		var colums = [{
				title: "仓单编号",
				name: "ReceiptId"
			},{
				title: "仓单总量",
				name: "TotalQuantity"
			},{
				title: "剩余量",
				name: "RemainingQuantity"
			},{
				title: "冻结量",
				name: "FrozenQuantity"
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

		for (var i = 0; i < data.length; i++){
			var tr = document.createElement("tr")
			table.lastChild.appendChild(tr)
			for(var j = 0; j < colums.length; j++){
				var td = document.createElement("td")
				tr.appendChild(td)
				td.innerHTML = data[i][colums[j].name]
			}
		}

	})


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
