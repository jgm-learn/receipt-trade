const path = require('path');

module.exports = [{
	entry: './src/register.js',
	output: {
		filename: 'register.js',
		path: '/home/jgm/goApp/src/receipt-trade/web/static/js'
	}
},{
	entry: './src/admin.js',
	output: {
		filename: 'admin.js',
		path: '/home/jgm/goApp/src/receipt-trade/web/static/js'
	}
},{
	entry: './src/user.js',
	output: {
		filename: 'user.js',
		path: '/home/jgm/goApp/src/receipt-trade/web/static/js'
	}
},{
	entry: './src/market.js',
	output: {
		filename: 'market.js',
		path: '/home/jgm/goApp/src/receipt-trade/web/static/js'
	} 
},{
	entry: './src/allData.js',
	output: {
		filename: 'allData.js',
		path: '/home/jgm/goApp/src/receipt-trade/web/static/js'
	} 
}
]

