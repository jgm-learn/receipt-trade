import {default as Web3} from 'web3'
import {default as contract } from 'truffle-contract'



var info = {
	UserName:		"",
	PassWord:		"",
	PublicKey:		""
};
var accounts;
async function start(){
	accounts = await web3.eth.getAccounts();
	submitInfo();
}

function submitInfo()
{
	$("#b1").on('click', function(){
		info.UserName = $("#account_in").val();
		info.PassWord	= $("#password_in").val();
		info.PublicKey = accounts[0];
		console.log("PublicKey = %s", info.PublicKey);
		console.log(info);
		console.log(JSON.stringify(info));
		$("#json_p").html(JSON.stringify(info));
		postInfo();
	});
}

function postInfo(){
	$.ajax({
		type:	'post',
		url:	"http://222.22.64.80:8081/register",
		data:	JSON.stringify(info),
		contentType:	"application/json",
		success:		function(data, state, xhr){
							console.log(state);
		},
		error:	function(xhr,state,error){
			console.log(state);
			console.log(error.toString());
		}
	});
	console.log("postInfo done");
}


window.addEventListener('load', function(){
	if (typeof web3 != 'undefined'){
		console.log("Using MetaMask web3.");
		window.web3 = new Web3(web3.currentProvider);
	}else {
		console.warn("No web3 detected.");
		window.web3 = new Web3(new Web3.providers.HttpProvider("http://192.168.22.247:8546"));
	}
	
	start();
});

