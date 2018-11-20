import {default as Web3} from 'web3'
import {default as contract} from 'truffle-contract'

//仓单属性
var receipt = {
	Class: 				"",	//品种
	ProductionDate:	"", //产期
	Level:				"",	//等级
	Warehouse:			"", //仓库
	Provenance:			""  //产地
}
//用户持有的仓单
var userReceipt = {
	UserId:			1,
	ReceiptId:		1,
	TotalQuantity:	1
}
//用户资金数据
var userFunds = {
	UserId:			1,
	TotalFunds:		1
}

//发送仓单信息
function addReceiptInfo(){
	//获取页面数据
	$("#bt_addReceiptInfo").on('click', function(){
		receipt.Class 			= 	$("#class").val();
		receipt.ProductionDate 	= 	$("#productionDate").val();
		receipt.Level			= 	$("#level").val();
		receipt.Warehouse		= 	$("#warehouse").val();
		receipt.Provenance		= 	$("#provenance").val();
		//序列化并发送
		$.ajax({
			type:			'post',
			url:			"http://222.22.64.80:8081/admin/addReceipt",
			data:			JSON.stringify(receipt),
			contentType:	"application/json",
			success:		function(data, state, xhr){
				console.log("发送成功");
				console.log(state);
			},
			error:			function(xhr, state, error){
				console.log("发送失败");
				console.log(error.toString());
			}
		});
	});
}

//为用户增加仓单
function addReceiptQty(){
	//获取页面数据
	$("#bt_addReceiptQty").on('click', function(){
		userReceipt.UserId = parseInt($("#userId1").val());
		userReceipt.ReceiptId = parseInt($("#receiptId").val());
		userReceipt.TotalQuantity 	= parseInt($("#totalQty").val());
		//序列化并发送
		$.ajax({
			type:			'post',
			url:			"http://222.22.64.80:8081/admin/addUserReceipt",
			data:			JSON.stringify(userReceipt),
			contentType:	"application/json",
			success:		function(data, state, xhr){
				console.log("userReceipt 发送成功");
				console.log(state);
			},
			error:			function(xhr, state, error){
				console.log("userReceipt 发送失败");
				console.log(error.toString());
			}
		});
	});
}

//为用户增加资金
function addUserFunds(){
	console.log("addUserFunds()")
	//获取页面数据
	$("#bt_addFundsQty").on('click', function(){
		userFunds.UserId 		= 	parseInt($("#userId2").val())
		userFunds.TotalFunds	= 	parseInt($("#fundsTotalQty").val())
		//序列化并发送
		$.ajax({
			type:			'post',
			url:			"http://222.22.64.80:8081/admin/addUserFunds",
			data:			JSON.stringify(userFunds),
			contentType:	"application/json",
			success:		function(data, state, xhr){
				console.log("userFunds 发送成功");
				console.log(state);
			},
			error:			function(xhr, state, error){
				console.log("userFunds 发送失败");
				console.log(error.toString());
			}
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
	addReceiptInfo();
	addReceiptQty();
	addUserFunds();
});
