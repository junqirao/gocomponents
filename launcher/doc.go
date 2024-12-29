// Package launcher
//
//	launch procedure:
//	0. prepare
//	1. execute init hooks
//	2. init dependent modules
//	   2.1. init meta
//	   2.2. execute updater if exists
//	   2.3. init registry
//	3. execute before hooks
//	4. execute blocked function
//	[blocked] 6. execute grace exit, wait for exit signal
package launcher
