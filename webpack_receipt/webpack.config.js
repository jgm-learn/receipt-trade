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
}]

