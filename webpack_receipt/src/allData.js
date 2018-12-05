function start(){
	getUser()
	getReceipt()
	getFunds()
	getOrderSell()
}

function createTable(tableId, row, data){
	//创建标题
		var table 	= 	document.getElementById(tableId)
		var thead 	=	document.createElement("thead")
		table.appendChild(thead)
			var tr 		= document.createElement("tr")
			thead.appendChild(tr)
		var tbody	= document.createElement("tbody")
		table.appendChild(tbody)
		for( var i=0; i < row.length; i++){
			var th = document.createElement("th")
			tr.appendChild(th)
			th.innerHTML = row[i].title
		}

		//将挂单数据填入表格
		for (var i = 0; i < data.length; i++){
			var tr = document.createElement("tr")
			table.lastChild.appendChild(tr)
			for(var j = 0; j < row.length; j++){
				var td = document.createElement("td")
				tr.appendChild(td)
				td.innerHTML = data[i][row[j].name]
			}
	}
}

function getUser(){
	$.getJSON("http://222.22.64.80:8081/allData/getUser", function(data){
		var row = [
			{title:"用户ID", 	name: "UserId"},
			{title:"用户名", 	name: "UserName"},
			{title:"密码", 		name: "PassWord"},
			{title:"账户地址", 	name: "PublicKey"},
			{title:"nonce", 	name: "Nonce"}];

		createTable("userTable", row, data);
	});

}

function getFunds(){
	$.getJSON("http://222.22.64.80:8081/allData/getFunds", function(data){
		var row = [
			{title:"用户ID", 	name: "UserId"},
			{title:"总资金", 	name: "TotalFunds"},
			{title:"可用资金", 	name: "AvailableFunds"},
			{title:"冻结资金", 	name: "FrozenFunds"}];

		createTable("fundsTable", row, data);
	});
}

function getReceipt(){
	$.getJSON("http://222.22.64.80:8081/allData/getReceipt", function(data){
		var row = [	
					{title: "用户",		name: "UserId"},
					{title: "仓单编号",	name: "ReceiptId"},
					{title: "总量",		name: "QtyTotal"},
					{title: "可用量",	name: "QtyAvailable"},
					{title: "冻结量",	name: "QtyFrozen"}];	
		
		createTable("receiptTable", row, data)
	});
		
}

function getOrderSell(){
	$.getJSON("http://222.22.64.80:8081/allData/getOrderSell", function(data){
		var row = [	
					{title: "挂单编号",		name: "Id"},
					{title: "用户ID",		name: "UserId"},
					{title: "仓单编号",		name: "ReceiptId"},
					{title: "挂牌数量",		name: "QtySell"},
					{title: "NonceSell",	name: "NonceSell"}];
		
		createTable("orderSellTable", row, data)
	});
}

start()
