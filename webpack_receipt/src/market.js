import {default as Web3} from 'web3'


//账户地址和用户id
var account;
var userId;

var orderBuy = {
	UserId:		0,
	ListId:		0,
	ReceiptId: 	0,
	QtyBuy:		0,
	NonceBuy:   0,
	SigBuy:		"",
	AddrBuy:	""
}


function start(){
	web3.eth.getAccounts(function(err, rst){
		account = rst[0];
		$("#userAddr")[0].innerHTML = account;
		getUserId();
		getReceiptList();
		delist()
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
					{title: "挂牌编号",	name: "ListId"},
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

//摘牌
function delist(){
	$("#bt_delist").on("click", function(){
		orderBuy.UserId 	=	userId
		orderBuy.ListId 	=	parseInt($("#listId").val())
		orderBuy.ReceiptId	=	parseInt($("#receiptId").val())
		orderBuy.QtyBuy 	=	parseInt($("#qtyBuy").val())
		$.getJSON("http://222.22.64.80:8081/user/getUserNonce", {userAddr: account}, async function(data){
			var nonceLast = data.Nonce;
			console.log("nonceLast: %d", nonceLast)
			orderBuy.NonceBuy 	= 	nonceLast + 1;
			var orderBuyHash	= 	await web3.utils.soliditySha3(orderBuy.ReceiptId, orderBuy.QtyBuy, orderBuy.NonceBuy);
			console.log("orderBuyHash", orderBuyHash);
			orderBuy.SigBuy 	= 	await web3.eth.sign(orderBuyHash, account)
			console.log("SigBuy", orderBuy.SigBuy);
			orderBuy.AddrBuy= account;
			//序列化并发送
			$.ajax({
				type:			'post',
				url:			"http://222.22.64.80:8081/market/delist",
				data:			JSON.stringify(orderBuy),
				dataType:		"json",
				contentType:	"application/json",
				success:		function(data, state, xhr){
					alert(data.Reply)
					console.log("orderBuy发送成功");
				},
				error:			function(xhr, state, error){
					console.log("orderBuy 发送失败");
					console.log(error.toString());
				}
			});
		});
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
